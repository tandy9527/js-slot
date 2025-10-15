package manager

import (
	"sync"

	"github.com/tandy9527/js-slot/core"
)

type GameManager struct {
	users map[int64]*core.User
	rooms map[string]*core.Room
	mu    sync.RWMutex
}

var manager *GameManager
var once sync.Once

func GetGameManager() *GameManager {
	once.Do(func() {
		manager = &GameManager{
			users: make(map[int64]*core.User),
			rooms: make(map[string]*core.Room),
		}
	})
	return manager
}

// ---------------- 玩家管理 ----------------

// 添加玩家
func (gm *GameManager) AddUser(u *core.User) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.users[u.UID] = u
}

// 获取玩家
func (gm *GameManager) GetUser(uid int64) *core.User {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	return gm.users[uid]
}

// 移除玩家
func (gm *GameManager) RemoveUser(uid int64) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	delete(gm.users, uid)
}

// ---------------- 房间管理 ----------------

// 添加房间
func (gm *GameManager) AddRoom(r *core.Room) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.rooms[r.ID] = r
}

// 获取房间
func (gm *GameManager) GetRoom(id string) *core.Room {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	return gm.rooms[id]
}

// 移除房间
func (gm *GameManager) RemoveRoom(id string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	delete(gm.rooms, id)
}
