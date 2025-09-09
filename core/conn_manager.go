package core

import (
	"sync"

	"github.com/tandy9527/js-slot/pkg/consts"
)

type ConnManager struct {
	mu    sync.RWMutex
	conns map[int64]*Connection
}

var (
	manager     *ConnManager
	managerOnce sync.Once
)

// GlobalManager 单例
func GlobalConnManager() *ConnManager {
	managerOnce.Do(func() {
		manager = &ConnManager{
			conns: make(map[int64]*Connection),
		}
	})
	return manager
}

// Add 添加连接
func (m *ConnManager) Add(conn *Connection) {
	if conn == nil || conn.UID == consts.USER_STATUS_NOLOGIN {
		return
	}

	m.mu.Lock()
	old, ok := m.conns[conn.UID]
	m.conns[conn.UID] = conn
	m.mu.Unlock()

	if ok && old != nil && old != conn {
		// 异步关闭旧连接，避免阻塞 Add 调用
		go func(c *Connection) {
			c.OnClose()
		}(old)
	}
}

// Remove 删除连接
func (m *ConnManager) Remove(uid int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.conns, uid)
}

// Get 获取连接
func (m *ConnManager) Get(uid int64) (*Connection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, ok := m.conns[uid]
	return conn, ok
}

// CloseAll 关闭所有连接
func (m *ConnManager) CloseAll() {
	m.mu.Lock()
	conns := make([]*Connection, 0, len(m.conns))
	for _, conn := range m.conns {
		conns = append(conns, conn)
	}
	m.conns = make(map[int64]*Connection) // 清空 map
	m.mu.Unlock()
	// // 并发关闭
	// for _, conn := range conns {
	// 	go conn.onClose()
	// }
	var wg sync.WaitGroup
	for _, conn := range conns {
		wg.Add(1)
		go func(c *Connection) {
			defer wg.Done()
			c.OnClose()
		}(conn)
	}
	wg.Wait() // 等待所有连接关闭，保证日志打印
}

func (m *ConnManager) SendByteAll(msg []byte) {
	m.mu.RLock()
	conns := make([]*Connection, 0, len(m.conns))
	for _, conn := range m.conns {
		if conn != nil {
			conns = append(conns, conn)
		}
	}
	m.mu.RUnlock()

	for _, conn := range conns {
		conn.SendByte(msg)
	}
}
func (m *ConnManager) SendJsonAll(msg string) {
	m.mu.RLock()
	conns := make([]*Connection, 0, len(m.conns))
	for _, conn := range m.conns {
		if conn != nil {
			conns = append(conns, conn)
		}
	}
	m.mu.RUnlock()

	for _, conn := range conns {
		conn.SendJSON(msg)
	}
}
