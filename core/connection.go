package core

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/errs"
	"github.com/tandy9527/js-slot/pkg/utils"
	"github.com/tandy9527/js-util/cache"
	"github.com/tandy9527/js-util/logger"

	"github.com/gorilla/websocket"
)

// Connection 包含 websocket 连接和用户ID
type Connection struct {
	SEQ     int64 // 每条消息的唯一ID，自增
	UID     int64
	Ws      *websocket.Conn
	TraceID string // 方便出问题日志

	sendChan  chan []byte
	mu        sync.Mutex
	closed    bool
	closeOnce sync.Once

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PingInterval time.Duration
}

// // 所有连接
// var (
// 	connections = make(map[int64]*Connection)
// 	connMu      sync.RWMutex
// )

// NewConnection 新建连接
func New(ws *websocket.Conn) *Connection {
	c := &Connection{
		Ws:           ws,
		TraceID:      utils.NewTraceID(),
		SEQ:          1,
		sendChan:     make(chan []byte, 512),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 10 * time.Second,
		PingInterval: 50 * time.Second,
	}
	return c
}
func (c *Connection) NextSeq() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.SEQ++
}

// // Add 注册连接
// func Add(c *Connection) {
// 	connMu.Lock()
// 	defer connMu.Unlock()
// 	if connections[c.UID] == nil {
// 		GetDB("db1").Incr(consts.REDIS_CONN_TOTAL)
// 	}
// 	connections[c.UID] = c
// }

// // Del 注销连接
// func Del(c *Connection) {
// 	connMu.Lock()
// 	defer connMu.Unlock()
// 	if _, ok := connections[c.UID]; ok {
// 		delete(connections, c.UID)
// 		close(c.sendChan)
// 		c.Ws.Close()
// 		GetDB("db1").Decr(consts.REDIS_CONN_TOTAL)
// 	}
// }

// SendJSON 将消息推送到发送队列（非阻塞）
func (c *Connection) SendJSON(v any) error {
	if c.closed {
		logger.Warn("[SendJSON] uid:%d closed", c.UID)
		return errs.ErrConnClosed
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	b, err := json.Marshal(v)
	if err != nil {
		logger.Error("[sendJson] uid:%d error: ", c.UID, err)
		return errs.ErrDataFormatError
	}

	select {
	case c.sendChan <- b:
		return nil
	default:
		return nil
	}
}

// onClose
func (c *Connection) OnClose() {
	c.closeOnce.Do(func() {
		logger.Info("closing connection uid :%d", c.UID)
		c.closed = true
		close(c.sendChan) // 停止写循环
		_ = c.Ws.Close()  // 关闭 websocke
		GlobalConnManager().Remove(c.UID)
		cache.GetDB("db2").ZIncrBy(consts.REDIS_GAME_ONLINE, -1, strconv.Itoa(GConf.GameID))
	})

}

func (c *Connection) AddGlobalConnManager() {
	GlobalConnManager().Add(c)
}

// // CloseAll 关闭所有连接
// func (c *Connection) CloseAll() {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	for uid, conn := range connections {
// 		conn
// 		delete(m.conns, uid)
// 	}
// }

// ReadPump 读取消息
func (c *Connection) ReadPump(onMessage func(c *Connection, msg []byte), onClose func(c *Connection)) {
	defer func() {
		onClose(c)
	}()

	c.Ws.SetReadLimit(512)
	c.Ws.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	c.Ws.SetPongHandler(func(string) error {
		c.Ws.SetReadDeadline(time.Now().Add(c.ReadTimeout))
		return nil
	})

	for {
		_, msg, err := c.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("[ReadPump] User %d read error: %v", c.UID, err)
				c.SendErr("", errs.ErrInternalServerError)
			}
			break
		}
		onMessage(c, msg)
	}
}

// WritePump 发送消息
func (c *Connection) WritePump() {
	ticker := time.NewTicker(c.PingInterval)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()

	for {
		select {
		case msg, ok := <-c.sendChan:
			c.Ws.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
			if !ok {
				c.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.Errorf("[WritePump] User %d write error: %v", c.UID, err)
				return
			}
			c.NextSeq()
		case <-ticker.C:
			c.Ws.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
			if err := c.Ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Errorf("[WritePump] User %d ping error: %v", c.UID, err)
				return
			}
		}
	}
}

// SendByte-异步
func (c *Connection) SendByte(msg []byte) error {
	if c.closed {
		logger.Warn("[SendByte] uid:%d closed", c.UID)
		return errs.ErrConnClosed
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case c.sendChan <- msg:
		return nil
	default:
		logger.Errorf("[SendByte] User %d send queue full, message dropped", c.UID)
		return nil
	}
}

// Broadcast 广播给所有用户
// func Broadcast(msg []byte) {
// 	connMu.RLock()
// 	defer connMu.RUnlock()
// 	for _, c := range connections {
// 		select {
// 		case c.sendChan <- msg:
// 		default:
// 			logger.Errorf("[Broadcast] User %d send queue full, message dropped", c.UID)
// 		}
// 	}
// }

// sendResponse 发送成功消息
func (c *Connection) SendResp(cmd string, data any) error {
	resp := RespMsg{
		Data:  data,
		Cmd:   cmd,
		Code:  errs.ErrSuccess.Code,
		Msg:   errs.ErrSuccess.Msg,
		Seq:   c.SEQ,
		Trace: c.TraceID,
	}
	return c.SendJSON(resp)
}

// sendError 发送错误消息
func (c *Connection) SendErr(cmd string, errs *errs.APIError) error {
	resp := RespMsg{
		Cmd:   cmd,
		Code:  errs.Code,
		Msg:   errs.Msg,
		Seq:   c.SEQ,
		Trace: c.TraceID,
	}
	return c.SendJSON(resp)
}
