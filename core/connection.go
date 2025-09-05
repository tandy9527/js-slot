package core

import (
	"encoding/json"
	"js-slot/pkg/consts"
	"js-slot/pkg/errs"
	"js-slot/pkg/logger"
	"js-slot/pkg/utils"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Connection 包含 websocket 连接和用户ID
type Connection struct {
	SEQ     int64 // 每条消息的唯一ID，自增
	UID     int64
	Ws      *websocket.Conn
	TraceID string // 方便出问题日志

	sendChan chan []byte
	mu       sync.Mutex

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PingInterval time.Duration
}

// 所有连接
var (
	connections = make(map[int64]*Connection)
	connMu      sync.RWMutex
)

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

// Add 注册连接
func Add(c *Connection) {
	connMu.Lock()
	defer connMu.Unlock()
	if connections[c.UID] == nil {
		GetDB("db1").Incr(consts.REDIS_CONN_TOTAL)
	}
	connections[c.UID] = c
}

// Del 注销连接
func Del(c *Connection) {
	connMu.Lock()
	defer connMu.Unlock()
	if _, ok := connections[c.UID]; ok {
		delete(connections, c.UID)
		close(c.sendChan)
		c.Ws.Close()
		GetDB("db1").Decr(consts.REDIS_CONN_TOTAL)
	}
}

// SendJSON 将消息推送到发送队列（非阻塞）
func (c *Connection) SendJSON(v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	select {
	case c.sendChan <- b:
		return nil
	default:
		return nil // 可改成返回 error 或丢弃
	}
}

// ReadPump 读取消息
func (c *Connection) ReadPump(onMessage func(c *Connection, msg []byte), onClose func(c *Connection)) {
	defer func() {
		Del(c)
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

// SendToUser 直接给指定用户发送消息
func SendToUser(uid int64, msg []byte) bool {
	connMu.RLock()
	conn, ok := connections[uid]
	connMu.RUnlock()
	if !ok {
		return false
	}

	select {
	case conn.sendChan <- msg:
		return true
	default:
		logger.Errorf("[SendToUser] User %d send queue full, message dropped", uid)
		return false
	}
}

// Broadcast 广播给所有用户
func Broadcast(msg []byte) {
	connMu.RLock()
	defer connMu.RUnlock()
	for _, c := range connections {
		select {
		case c.sendChan <- msg:
		default:
			logger.Errorf("[Broadcast] User %d send queue full, message dropped", c.UID)
		}
	}
}

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
	if err := c.SendJSON(resp); err != nil {
		return err
	}
	return nil
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
	if err := c.SendJSON(resp); err != nil {
		return err
	}
	return nil
}
