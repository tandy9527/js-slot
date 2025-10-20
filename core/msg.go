package core

import (
	"fmt"
	"strconv"

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
	Win  int64          `json:"-"` // 赢多少,会根据该值去结算
}

func (m *Message) GetMap() map[string]any {
	if m == nil || m.Data == nil {
		return nil
	}
	mp, ok := m.Data.(map[string]any)
	if !ok {
		return nil
	}
	return mp
}

func (m *Message) GetString(key string) string {
	v := m.GetMap()[key]
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%v", val)
	case int:
		return strconv.Itoa(val)
	default:
		return ""
	}
}

func (m *Message) GetInt(key string) int {
	v := m.GetMap()[key]
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}
func (m *Message) GetInt64(key string) int64 {
	v := m.GetMap()[key]
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case string:
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	default:
		return 0
	}
}
func (m *Message) GetBool(key string) bool {
	v := m.GetMap()[key]
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true" || val == "1"
	case float64:
		return val != 0
	default:
		return false
	}
}

func (m *Message) GetFloat64(key string) float64 {
	v := m.GetMap()[key]
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}
