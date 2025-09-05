package core

import (
	"github.com/tandy9527/js-slot/pkg/errs"
)

// Message represents a WebSocket message.
type Message struct {
	Cmd  string `json:"cmd"`  // 命令
	Data any    `json:"data"` // 数据
	UID  int64  `json:"-"`    // 用户ID, 仅服务器内部使用
}

type RespMsg struct {
	Data  any       `json:"data,omitempty"`
	Seq   int64     `json:"seq,omitempty"` // 消息序列号
	Code  errs.Code `json:"code,omitempty"`
	Msg   string    `json:"msg,omitempty"`
	Cmd   string    `json:"cmd"`
	Trace string    `json:"trace,omitempty"` //  方便追踪错误
}

// GameResult
type GameResult struct {
	Data any            `json:"data,omitempty"`
	Err  *errs.APIError `json:"err,omitempty"`
}
