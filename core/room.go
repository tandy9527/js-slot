package core

import (
	"sync"

	"github.com/tandy9527/js-slot/core/game"
)

type Room struct {
	ID       string
	Users    map[int64]*User
	Status   string // waiting / playing
	mu       sync.RWMutex
	gameInfo *game.GameInfo //  数值配置
}

func (r *Room) GetGameInfo() *game.GameInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.gameInfo
}

func NewRoom(id string) *Room {
	return &Room{
		ID:       id,
		Users:    make(map[int64]*User),
		Status:   "waiting",
		gameInfo: game.GetGameInfo(),
	}
}

func (r *Room) AddUser(u *User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Users[u.UID] = u
	u.CurrentRoom = r.ID
}

func (r *Room) GetUser(uid int64) *User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Users[uid]
}

func (r *Room) RemoveUser(uid int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.Users[uid]; ok {
		u.CurrentRoom = ""
		delete(r.Users, uid)
	}
}

func (r *Room) GetUsers() []*User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*User, 0, len(r.Users))
	for _, u := range r.Users {
		list = append(list, u)
	}
	return list
}

func (r *Room) UserCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Users)
}

func (r *Room) SetStatus(status string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Status = status
}

func (r *Room) GetStatus() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Status
}

func (r *Room) SendAllUsers(data any) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.Users {
		u.Conn.SendJSON(data)
	}
}
