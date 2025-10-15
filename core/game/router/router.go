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
	"github.com/tandy9527/js-slot/pkg/errs"
)

type BaseGameInterface interface {
	HandleMessage(conn *core.Connection, msg core.Message) error
}

var Router = NewRouter(3 * time.Second)

// GameHandlerFunc 异步处理消息
type GameHandlerFunc func(ctx context.Context, msg core.Message) <-chan core.GameResult

type GameRouter struct {
	// 所有命令处理器
	handlers map[string]GameHandlerFunc
	// 超时时间
	Timeout time.Duration
}

func NewRouter(timeout time.Duration) *GameRouter {
	return &GameRouter{
		handlers: make(map[string]GameHandlerFunc),
		Timeout:  timeout,
	}
}
func (g *GameRouter) GetHandler(cmd string) GameHandlerFunc {
	handler, ok := g.handlers[cmd]
	if !ok {
		return nil
	}
	return handler
}

// WrapSyncHandler 将同步函数包装为异步处理
func WrapSyncHandler(f func(core.Message) core.GameResult) GameHandlerFunc {
	return func(ctx context.Context, msg core.Message) <-chan core.GameResult {
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

			res := f(msg)
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
	handler, ok := g.handlers[msg.Cmd]
	if !ok {
		return conn.SendErr(msg.Cmd, errs.ErrCmdNotFound)
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
	defer cancel()

	resultCh := handler(ctx, msg)

	select {
	case <-ctx.Done():
		return conn.SendErr(msg.Cmd, errs.ErrTimeout)
	case res := <-resultCh:
		if res.Err != nil {
			return conn.SendErr(msg.Cmd, res.Err)
		}
		return conn.SendResp(msg.Cmd, res.Data)
	}
}
