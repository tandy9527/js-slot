package core

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/logger"
)

// 游戏数据结构--BalanceChangeData 余额变动
type BalanceChangeData struct {
	UID           int64  // 用户ID
	Time          int64  // Unix 10位秒时间戳
	Amount        int64  // 变化金额
	BalanceBefore int64  // 变化前余额
	BalanceAfter  int64  // 变化后余额
	Type          int8   // 1 下注 2 游戏结算
	GameCode      string // 游戏代码
	UUID          string // UUID
}

// BalanceChanges 发送单条消息到 Redis,
func BalanceChanges(data *BalanceChangeData) error {
	// 自动生成 UUID
	data.UUID = uuid.New().String()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = GetDB("db8").LPush(consts.REDIS_DATA_QUEUE_PENDING, bytes)
	logger.Infof("PushDataRedis : %s", string(bytes))
	return err
}

// 游戏数据结构--GameLogData 游戏历史日志
type GameLogData struct {
	UID      int64  // 用户ID
	GameCode string // 游戏代码
}
