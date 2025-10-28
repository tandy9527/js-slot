// author: tandy  2025.9.1
// router.go 文件实现了 GameRouter，用于管理游戏消息的注册、路由和异步处理。
// 核心功能包括：
// 1. 命令处理器注册（Register）。
// 2. 异步执行消息处理函数（WrapSyncHandler）。
// 3. 处理消息入口（HandleMessage），包含超时和错误处理。

package router

import (
	"context"
	"runtime/debug"
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
	//gr.Register(consts.REQ_CMD_LOGIN, WrapSyncHandler(core.Login))
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
					logger.Errorf("[panic recovered] %v\n%s", r, debug.Stack())
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
	logger.Infof("req -> cmd:[%s] uid:[%d], data:{%+v}", msg.Cmd, msg.UID, msg.Data)
	// login
	if msg.Cmd == consts.REQ_CMD_LOGIN {
		// ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
		// defer cancel()

		// resultCh := core.Login(conn, msg)
		// select {
		// case <-ctx.Done():
		// 	return conn.SendErr(msg.Cmd, errs.ErrTimeout)
		// case res := <-resultCh:
		// 	if res.Err != nil {
		// 		return conn.SendErr(msg.Cmd, res.Err)
		// 	}
		// 	return conn.SendResp(msg.Cmd, res.Data)
		// }
		res := Login(conn, msg)
		if res.Err != nil {
			return conn.SendErr(msg.Cmd, res.Err)
		}
		logger.Infof("resp <- cmd:[%s][%d]ms, data:{%+v}", msg.Cmd, utils.RunTime(startTime), res.Data)
		return conn.SendResp(msg.Cmd, res.Data)
	}
	// 非 login 类型必须有 user
	handler, ok := g.handlers[msg.Cmd]
	if !ok {
		return conn.SendErr(msg.Cmd, errs.ErrCmdNotFound)
	}

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

	isBet, bet := SpinBet(user, msg)
	if !isBet {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
	defer cancel()
	resultCh := handler(ctx, user, gameinfo, msg)

	select {
	case <-ctx.Done():
		logger.Errorf("resp <- cmd:[%s] timeout:[%v]ms", msg.Cmd, g.Timeout.Milliseconds())
		return conn.SendErr(msg.Cmd, errs.ErrTimeout)
	case res := <-resultCh:
		if res.Err != nil {
			return conn.SendErr(msg.Cmd, res.Err)
		}
		logger.Infof("resp <- cmd:[%s][%d]ms, uid:[%d], data:{%+v}", msg.Cmd, utils.RunTime(startTime), msg.UID, res.Data)

		// 结算逻辑
		SpinRecord(msg.Cmd, user, bet, res.Win, user.Balance, res.Data)
		if res.Win > 0 {
			user.GameEnd(res.Win)
		}

		return conn.SendByBalance(msg.Cmd, res.Data, user.Balance)
	}
}

// spin  统一下注处理
func SpinBet(user *core.User, msg core.Message) (bool, int64) {
	if msg.Cmd == consts.REQ_CMD_SPIN {
		err := user.Bet(msg.GetInt64("bet"))
		if err != nil {
			user.Conn.SendErr(msg.Cmd, err)
			return false, 0
		} else {
			user.BalanceChange()
		}
	}
	return true, msg.GetInt64("bet")
}

func SpinRecord(cmd string, user *core.User, bet int64, win int64, balance int64, data any) {
	if cmd == consts.REQ_CMD_SPIN {
		core.SlotSpinRecord(user.UID, bet, win, balance, data, utils.CurrentTimestamp())
	}
}
