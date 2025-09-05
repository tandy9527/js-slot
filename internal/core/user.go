package core

import (
	"js-slot/pkg/errs"
	"js-slot/pkg/scripts"
	"js-slot/pkg/utils"
	"sync"
	"time"
)

type User struct {
	UID         int64
	Conn        *Connection
	Balance     int64
	CurrentGame string
	CurrentRoom string
	LastActive  time.Time
	mu          sync.Mutex
}

func NewUser(uid int64, conn *Connection) *User {
	return &User{
		UID:        uid,
		Conn:       conn,
		Balance:    0,
		LastActive: time.Now(),
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
	result, err := GetDB("db0").ExecLua(scripts.RechargeLua, []string{utils.GetUserRedisKey(u.UID)}, amount)
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
