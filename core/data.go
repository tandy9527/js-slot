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
	Ctime         int64  `json:"ctime,omitempty"`          // Unix 10位秒时间戳
	Amount        int64  `json:"amount,omitempty"`         // 变化金额
	BalanceBefore int64  `json:"balance_before,omitempty"` // 变化前余额
	BalanceAfter  int64  `json:"balance_after,omitempty"`  // 变化后余额
	Type          int8   `json:"type,omitempty"`           // 1 下注 2 游戏结算
	GameCode      string `json:"game_code,omitempty"`      // 游戏代码
	GameID        int    `json:"game_id,omitempty"`        // 游戏ID

}

// 游戏数据结构--GameLogData 游戏历史日志
type GameLogData struct {
	UID      int64  // 用户ID
	GameCode string `json:"game_code,omitempty"` // 游戏代码
	GameID   int    `json:"game_id,omitempty"`   // 游戏ID
	Bet      int64  `json:"bet,omitempty"`       // 下注
	Win      int64  `json:"win,omitempty"`       // 赢
	Balance  int64  `json:"balance,omitempty"`   // 余额
	Ctime    int64  `json:"ctime,omitempty"`     // Unix 10位秒时间戳
	Data     any    `json:"data,omitempty"`      // 游戏数据

}

// 落库的数据
type Data struct {
	Data any    `json:"data"`           // 存储的数据
	Type string `json:"cmd"`            // 命令
	UUID string `json:"uuid,omitempty"` // UUID
}

// 数据持久化
func Persistent(data Data) error {
	// 自动生成 UUID
	data.UUID = uuid.New().String()

	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Persistent error: %v", err)
		return err
	}
	_, err = cache.GetDB("db8").LPush(consts.REDIS_DATA_QUEUE_PENDING, bytes)
	if err != nil {
		logger.Errorf("Persistent error: %v", err)
		return err
	}
	logger.Infof("Persistent-> [%s] Push Data Redis : %s", data.Type, string(bytes))
	return nil
}

// 余额变动
func UpdataBalance(uid int64, amount int64, balanceBefore int64, blanceAfter int64, typ int8, time int64) {
	d := BalanceChangeData{
		UID:           uid,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  blanceAfter,
		Ctime:         time,
		Type:          typ,
		GameCode:      GConf.GameCode,
		GameID:        GConf.GameID,
	}
	Persistent(Data{
		Data: d,
		Type: consts.DATA_PERSISTENT_TYPE_UPDATE,
	})
}

func SlotSpinRecord(uid int64, bet int64, win int64, balance int64, data any, time int64) {
	d := GameLogData{
		UID:      uid,
		GameCode: GConf.GameCode,
		GameID:   GConf.GameID,
		Bet:      bet,
		Win:      win,
		Balance:  balance,
		Ctime:    time,
		Data:     data,
	}
	Persistent(Data{
		Data: d,
		Type: consts.DATA_PERSISTENT_TYPE_SPIN,
	})
}
