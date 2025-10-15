package core

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-util/cache"
	"github.com/tandy9527/js-util/logger"
)

// 游戏数据结构--BalanceChangeData 余额变动
type BalanceChangeData struct {
	UID           int64  `json:"uid,omitempty"`            // 用户ID
	Time          int64  `json:"time,omitempty"`           // Unix 10位秒时间戳
	Amount        int64  `json:"amount,omitempty"`         // 变化金额
	BalanceBefore int64  `json:"balance_before,omitempty"` // 变化前余额
	BalanceAfter  int64  `json:"balance_after,omitempty"`  // 变化后余额
	Type          int8   `json:"type,omitempty"`           // 1 下注 2 游戏结算
	GameCode      string `json:"game_code,omitempty"`      // 游戏代码
	UUID          string `json:"uuid,omitempty"`           // UUID
}

// BalanceChanges 发送单条消息到 Redis,
func BalanceChanges(data *BalanceChangeData) error {
	// 自动生成 UUID
	data.UUID = uuid.New().String()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = cache.GetDB("db8").LPush(consts.REDIS_DATA_QUEUE_PENDING, bytes)
	if err != nil {

	}
	logger.Infof("BalanceChanges-> Push Data Redis : %s", string(bytes))
	return err
}

// 游戏数据结构--GameLogData 游戏历史日志
type GameLogData struct {
	UID      int64  // 用户ID
	GameCode string // 游戏代码
}
