package router

import (
	"time"

	"github.com/tandy9527/js-slot/core"
	"github.com/tandy9527/js-slot/core/game/manager"
	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/errs"
	"github.com/tandy9527/js-slot/pkg/utils"
	"github.com/tandy9527/js-util/cache"
	"github.com/tandy9527/js-util/logger"
	"github.com/tandy9527/js-util/tools/jwt_tools"
	"github.com/tandy9527/js-util/tools/str_tools"
)

// 验证用户信息
func Login(conn *core.Connection, msg core.Message) core.GameResult {
	token := msg.GetString("t")
	if str_tools.IsEmpty(token) {
		return core.GameResult{Err: errs.ErrDataFormatError}
	}
	secret, err := cache.GetDB("db1").Get(consts.REDIS_SLOTS_JWT_KEY)
	if err != nil {
		logger.Errorf("auth error: %v", err)
		return core.GameResult{Err: errs.ErrInternalServerError}
	}
	if str_tools.IsEmpty(secret) {
		secret = str_tools.RandNumStr(50, 100)
		cache.GetDB("db1").Set(consts.REDIS_SLOTS_JWT_KEY, secret, 3*time.Hour)
	}
	//token
	uid, err := jwt_tools.ParseToken(token, secret, conn.IP)
	if err != nil {
		return core.GameResult{Err: errs.ErrTokenNull}
	}
	if uid < 0 {
		return core.GameResult{Err: errs.ErrTokenNull}
	}

	//  auth
	if ok, _ := cache.GetDB("db0").Exists(utils.GetUserRedisKey(uid)); !ok {
		return core.GameResult{Err: errs.ErrUserNotFound}
	}

	conn.UID = uid
	conn.AddGlobalConnManager()
	user := core.NewUser(uid, conn)
	room := core.NewRoom("room:888888")
	room.AddUser(user)
	m := manager.GetGameManager()
	m.AddUser(user)
	m.AddRoom(room)
	return core.GameResult{Data: map[string]string{"S": "login success"}}
}
