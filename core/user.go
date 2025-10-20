package core

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/tandy9527/js-slot/pkg/consts"
	"github.com/tandy9527/js-slot/pkg/errs"
	"github.com/tandy9527/js-slot/pkg/scripts"
	"github.com/tandy9527/js-slot/pkg/utils"
	"github.com/tandy9527/js-util/cache"
)

type User struct {
	UID         int64
	Conn        *Connection
	Balance     int64
	CurrentGame string
	RoomID      string
	LastActive  time.Time
	mu          sync.Mutex
	Session     string
	Room        *Room
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
	// u.mu.Lock()
	// defer u.mu.Unlock()
	fmt.Println(data)
	u.Conn.SendJSON(data)
}

func (u *User) UpdateBalance(amount int64) *errs.APIError {
	u.mu.Lock()
	defer u.mu.Unlock()
	result, err := cache.GetDB("db0").ExecLua(scripts.RechargeLua, []string{u.Session}, amount)
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
	result, err := cache.GetDB("db0").ExecLua(scripts.BetLua, []string{u.Session}, bet)
	fmt.Println(result, err)
	if err != nil {
		return errs.ErrInternalServerError
	}
	balance := result.(int64)

	switch balance {
	case -1: // 下注金额错误 下注金额等于或者小于0
		//u.Conn.SendErr("", errs.ErrWrongBetAmount)
		return errs.ErrWrongBetAmount
	case -2: // 余额不足
		//u.Conn.SendErr("", errs.ErrInsufficientBalance)
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

// 余额有变动
func (u *User) BalanceChange() int64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	if balanceStr, err := cache.GetDB("db0").HGet(u.Session, "balance"); err == nil {
		u.Balance, _ = strconv.ParseInt(balanceStr, 10, 64)
	}
	fmt.Println(u.Balance)
	u.SendResp(RespMsg{Cmd: consts.RESP_CMD_BALANCE_CHANGE, Data: u.Balance, Code: errs.ErrSuccess.Code, Msg: errs.ErrSuccess.Msg, Seq: u.Conn.SEQ, Trace: u.Conn.TraceID})
	return u.Balance
}
