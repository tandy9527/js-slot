package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tandy9527/js-slot/core"
	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/errs"
	"github.com/tandy9527/js-slot/pkg/utils"
	"github.com/tandy9527/js-util/cache"
	"github.com/tandy9527/js-util/logger"
	"github.com/tandy9527/js-util/tools/jwt_tools"
	"github.com/tandy9527/js-util/tools/str_tools"

	"github.com/tandy9527/js-slot/core/game/manager"
	"github.com/tandy9527/js-slot/core/game/router"

	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
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
	//sendAES(conn)
}

// func sendAES(c *core.Connection) error {
// 	// TODO:  后面处理
// 	c.SendJSON(core.Message{
// 		Cmd:  "aes",
// 		Data: "ok",
// 	})
// 	return nil
// }

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

func auth(token, ip string) (int64, *errs.APIError) {
	if str_tools.IsEmpty(t) {
		return -1, errs.ErrTokenNull
	}
	//token := str_tools.Base64Decode(t)
	secret, err := cache.GetDB("db1").Get(consts.REDIS_SLOTS_JWT_KEY)
	if err != nil {
		logger.Errorf("auth error: %v", err)
		return -1, errs.ErrInternalServerError
	}
	if str_tools.IsEmpty(secret) {
		secret = str_tools.RandNumStr(50, 100)
		cache.GetDB("db1").Set(consts.REDIS_SLOTS_JWT_KEY, secret, 3*time.Hour)
	}
	//token
	uid, err := jwt_tools.ParseToken(token, secret, ip)
	if err != nil {
		return -1, errs.ErrTokenErr
	}
	if uid < 0 {
		return -1, errs.ErrTokenErr
	}
	return uid, nil

}

func login(c *core.Connection, token string) {
	// TODO:  后面会改成各种鉴权方式
	// GRPC 进行鉴权
	// uid, _ := strconv.ParseInt(token, 10, 64)
	uid, err := auth(token, "")
	if err != nil {
		c.SendErr("", err)
		return
	}
	if ok, _ := cache.GetDB("db0").Exists(utils.GetUserRedisKey(uid)); !ok {
		c.SendErr("", errs.ErrUserNotFound)
		return
	}

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
}

// onCloseHandler 连接关闭后的清理逻辑
func onCloseHandler(c *core.Connection) {
	// TODO
	//logger.Infof("[onCloseHandler] User %d disconnected\n", c.UID)
	c.OnClose()
}
