// author: tandy  2025.9.1
// router.go 文件实现了 GameRouter，用于管理游戏消息的注册、路由和异步处理。
// 核心功能包括：
// 1. 命令处理器注册（Register）。
// 2. 异步执行消息处理函数（WrapSyncHandler）。
// 3. 处理消息入口（HandleMessage），包含超时和错误处理。

package router

import (
	"context"
	"time"

	"github.com/tandy9527/js-slot/core"
	"github.com/tandy9527/js-slot/core/game"
	"github.com/tandy9527/js-slot/core/game/manager"
	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/errs"
	"github.com/tandy9527/js-slot/pkg/utils"
	"github.com/tandy9527/js-util/logger"
)

// type BaseGameInterface interface {
// 	HandleMessage(conn *core.Connection, msg core.Message) error
// }

var Router = NewRouter(3 * time.Second)

// GameHandlerFunc 异步处理消息
// user: 请求的用户
// gameinfo: 游戏数值配置
// msg: 请求的具体消息数据
type GameHandlerFunc func(ctx context.Context, user *core.User, gameinfo *game.GameInfo, msg core.Message) <-chan core.GameResult

type GameRouter struct {
	// 所有命令处理器
	handlers map[string]GameHandlerFunc
	// 超时时间
	Timeout time.Duration
}

func NewRouter(timeout time.Duration) *GameRouter {
	gr := &GameRouter{
		handlers: make(map[string]GameHandlerFunc),
		Timeout:  timeout,
	}
	gr.Register(consts.REQ_CMD_GET_BALANCE, WrapSyncHandler(core.GetBalance))
	return gr
}
func (g *GameRouter) GetHandler(cmd string) GameHandlerFunc {
	handler, ok := g.handlers[cmd]
	if !ok {
		return nil
	}
	return handler
}

// WrapSyncHandler 将同步函数包装为异步处理
func WrapSyncHandler(f func(*core.User, *game.GameInfo, core.Message) core.GameResult) GameHandlerFunc {
	return func(ctx context.Context, user *core.User, gameinfo *game.GameInfo, msg core.Message) <-chan core.GameResult {
		ch := make(chan core.GameResult, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// 遇到 panic 时通过 channel 发送错误（若 ctx 已取消则放弃）
					select {
					case <-ctx.Done():
						// 上游已取消，不发送
					default:
						ch <- core.GameResult{Err: errs.ErrInternalServerError}
					}
				}
				close(ch)
			}()

			res := f(user, gameinfo, msg)
			// 在发送结果前尊重 ctx，避免在上游已取消时阻塞
			select {
			case <-ctx.Done():
				// 上游已取消/超时，丢弃结果
			default:
				ch <- res
			}
		}()
		return ch
	}
}

// Register 注册命令处理器
func (g *GameRouter) Register(cmd string, handler GameHandlerFunc) {
	if g.handlers == nil {
		g.handlers = make(map[string]GameHandlerFunc)
	}
	g.handlers[cmd] = handler
}

// HandleMessage  处理消息入口
func (g *GameRouter) HandleMessage(conn *core.Connection, msg core.Message) error {
	startTime := utils.StartTime()
	handler, ok := g.handlers[msg.Cmd]
	logger.Infof("req -> cmd:[%s] uid:[%d],data:{%+v}", msg.Cmd, msg.UID, msg.Data)
	if !ok {
		return conn.SendErr(msg.Cmd, errs.ErrCmdNotFound)
	}
	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
	defer cancel()
	manager := manager.GetGameManager()
	if manager == nil {
		return conn.SendErr(msg.Cmd, errs.ErrInternalServerError)
	}
	user := manager.GetUser(msg.UID)
	if user == nil {
		logger.Errorf("cmd:%s user:%d not found", msg.Cmd, msg.UID)
		return conn.SendErr(msg.Cmd, errs.ErrPleaseLogIn)
	}
	gameinfo := user.Room.GetGameInfo()
	if gameinfo == nil {
		logger.Errorf("cmd:%s gameinfo not found", msg.Cmd)
		return conn.SendErr(msg.Cmd, errs.ErrInternalServerError)
	}
	if !bet(user, msg) {
		return nil
	}
	resultCh := handler(ctx, user, gameinfo, msg)

	select {
	case <-ctx.Done():
		logger.Errorf("resp <- cmd:[%s] timeout:[%v]ms", msg.Cmd, g.Timeout.Milliseconds())
		return conn.SendErr(msg.Cmd, errs.ErrTimeout)
	case res := <-resultCh:
		if res.Err != nil {
			return conn.SendErr(msg.Cmd, res.Err)
		}
		logger.Infof("resp <- cmd:[%s][%d]ms, uid:[%d],data:{%+v}", msg.Cmd, utils.RunTime(startTime), msg.UID, res.Data)
		//   去结算
		if res.Win > 0 {
			user.GameEnd(res.Win)
			dataMap := res.Data.(map[string]any)
			dataMap["balance"] = user.Balance
		}
		return conn.SendResp(msg.Cmd, res.Data)
	}
}

// spin  统一下注处理
func bet(user *core.User, msg core.Message) bool {
	if msg.Cmd == consts.REQ_CMD_SPIN {
		err := user.Bet(msg.GetInt64("bet"))
		if err != nil {
			user.Conn.SendErr(msg.Cmd, err)
			return false
		} else {
			user.BalanceChange()
		}
	}
	return true
}
