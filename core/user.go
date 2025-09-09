package core

import (
	"sync"
	"time"

	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/errs"
	"github.com/tandy9527/js-slot/pkg/scripts"
	"github.com/tandy9527/js-slot/pkg/utils"
)

type User struct {
	UID         int64
	Conn        *Connection
	Balance     int64
	CurrentGame string
	CurrentRoom string
	LastActive  time.Time
	mu          sync.Mutex
	Session     string
}

func NewUser(uid int64, conn *Connection) *User {
	return &User{
		UID:        uid,
		Conn:       conn,
		Balance:    0,
		LastActive: time.Now(),
		Session:    utils.GetUserRedisKey(uid),
	}
}

func (u *User) SendResp(data any) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Conn.SendJSON(data)
}

func (u *User) UpdateBalance(amount int64) *errs.APIError {
	u.mu.Lock()
	defer u.mu.Unlock()
	result, err := GetDB("db0").ExecLua(scripts.RechargeLua, []string{u.Session}, amount)
	if err != nil || result == nil {
		return errs.ErrInternalServerError
	}
	resSlice := result.([]any)
	if len(resSlice) != 2 {
		return errs.ErrInternalServerError
	}
	u.Balance += amount
	u.LastActive = time.Now()
	return nil
}

func (u *User) Bet(bet int64) *errs.APIError {
	u.mu.Lock()
	defer u.mu.Unlock()
	result, err := GetDB("db0").ExecLua(scripts.BetLua, []string{u.Session}, bet)
	if err != nil {
		return errs.ErrInternalServerError
	}
	balance := result.(int64)

	switch balance {
	case -1: // 下注金额错误 下注金额等于或者小于0
		return errs.ErrWrongBetAmount
	case -2: // 余额不足
		return errs.ErrInsufficientBalance
	}
	BalanceChanges(&BalanceChangeData{
		UID:           u.UID,
		Time:          utils.CurrentTimestamp(),
		Amount:        bet,
		Type:          consts.TYPE_BET,
		BalanceBefore: balance + bet,
		BalanceAfter:  balance,
		GameCode:      GConf.GameCode,
	})
	u.Balance = balance
	return nil
}
