package ws

import (
	"encoding/json"

	"github.com/tandy9527/js-slot/core"
	"github.com/tandy9527/js-slot/pkg/errs"

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
	m.UID = c.UID
	router.Router.HandleMessage(c, m)
}

// onCloseHandler 连接关闭后的清理逻辑
func onCloseHandler(c *core.Connection) {
	// TODO
	//logger.Infof("[onCloseHandler] User %d disconnected\n", c.UID)
	c.OnClose()
}
