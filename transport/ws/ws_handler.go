package ws

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tandy9527/js-slot/core"
	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/errs"

	"github.com/tandy9527/js-slot/core/game/manager"
	"github.com/tandy9527/js-slot/core/game/router"

	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}
var jwtSecret = []byte("9f2c4b6d8a1e4f5b7c9d2e3f1a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b")

// Claims 定义
type Claims struct {
	jwt.RegisteredClaims
	UID int64 `json:"uid"`
}

func parseToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return -1, err
	}
	if !token.Valid {
		return -1, nil
	}
	claims := token.Claims.(*Claims)
	if claims == nil {
		return 0, err
	}
	return claims.UID, nil
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	//  TODO 先这样处理，后面会改成各种鉴权方式
	// uid := r.URL.Query().Get("uid")

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// uidInt, _ := strconv.ParseInt(uid, 10, 64)
	conn := core.New(wsConn)
	// 连接建立后，启动读写协程
	go conn.ReadPump(onMessageHandler, onCloseHandler)
	go conn.WritePump()

	// user := core.NewUser(uidInt, conn)
	// room := core.NewRoom("room:1")
	// room.AddUser(user)
	// manager := game.GetGameManager()
	// manager.AddUser(user)
	// manager.AddRoom(room)
	sendAES(conn)
}

func sendAES(c *core.Connection) error {
	// TODO:  后面处理
	c.SendJSON(core.Message{
		Cmd:  "aes",
		Data: "ok",
	})
	return nil
}

// 收到消息后的回调函数
func onMessageHandler(c *core.Connection, msg []byte) {
	var m core.Message
	err := json.Unmarshal(msg, &m)
	if err != nil {
		c.SendErr("", errs.ErrDataFormatError)
		return
	}
	if m.Cmd == consts.LOGIN_CMD {
		dataMap, ok := m.Data.(map[string]any)
		if !ok {
			c.SendErr("", errs.ErrDataFormatError)
			return
		}
		if val, ok := dataMap["token"]; ok {
			token, ok := val.(string)
			if ok {
				login(c, token)
			} else {
				c.SendErr("", errs.ErrDataFormatError)
			}

		} else {
			c.SendErr("", errs.ErrDataFormatError)
		}
		return
	}
	m.UID = c.UID
	router.Router.HandleMessage(c, m)
}

func login(c *core.Connection, token string) error {
	// TODO:  后面会改成各种鉴权方式
	// GRPC 进行鉴权
	uid, _ := strconv.ParseInt(token, 10, 64)

	c.UID = uid
	c.AddGlobalConnManager()
	user := core.NewUser(uid, c)
	room := core.NewRoom("room:1")
	room.AddUser(user)
	manager := manager.GetGameManager()
	manager.AddUser(user)
	manager.AddRoom(room)
	c.SendJSON(core.Message{
		Cmd:  "login",
		Data: fmt.Sprintf("user %d login success", uid),
	})
	return nil
}

// onCloseHandler 连接关闭后的清理逻辑
func onCloseHandler(c *core.Connection) {
	// TODO
	//logger.Infof("[onCloseHandler] User %d disconnected\n", c.UID)
	c.OnClose()
}
