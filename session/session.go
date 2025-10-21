package session

import (
	"errors"
	"github.com/xiaoyi510/xbot/utils"
	"sync"
	"time"
)

// Session 会话
type Session struct {
	ID        string
	UserID    int64
	GroupID   int64
	Data      map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
	ch        chan interface{}
}

// Manager 会话管理器
type Manager struct {
	store   Store
	ttl     time.Duration
	mu      sync.RWMutex
	waiting *utils.SafeMap[string, *Session]
}

// NewManager 创建会话管理器
func NewManager(store Store, ttl time.Duration) *Manager {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}

	return &Manager{
		store:   store,
		ttl:     ttl,
		waiting: utils.NewSafeMap[string, *Session](),
	}
}

// Get 获取会话
func (m *Manager) Get(userID, groupID int64) (*Session, bool) {
	key := utils.GenerateSessionKey(userID, groupID)
	session, err := m.store.Get(key)
	if err != nil || session == nil {
		return nil, false
	}
	return session, true
}

// Set 设置会话
func (m *Manager) Set(session *Session) error {
	session.UpdatedAt = time.Now()
	key := utils.GenerateSessionKey(session.UserID, session.GroupID)
	return m.store.Set(key, session, m.ttl)
}

// Delete 删除会话
func (m *Manager) Delete(userID, groupID int64) error {
	key := utils.GenerateSessionKey(userID, groupID)
	return m.store.Delete(key)
}

// CreateWaitSession 创建等待会话
func (m *Manager) CreateWaitSession(userID, groupID int64, timeout time.Duration) *Session {
	key := utils.GenerateSessionKey(userID, groupID)

	session := &Session{
		ID:        key,
		UserID:    userID,
		GroupID:   groupID,
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ch:        make(chan interface{}, 1),
	}

	m.waiting.Set(key, session)

	// 超时自动删除
	go func() {
		time.Sleep(timeout)
		m.waiting.Delete(key)
	}()

	return session
}

// NotifyWaitSession 通知等待会话
func (m *Manager) NotifyWaitSession(userID, groupID int64, data interface{}) bool {
	key := utils.GenerateSessionKey(userID, groupID)

	session, ok := m.waiting.Get(key)
	if !ok {
		return false
	}

	select {
	case session.ch <- data:
		m.waiting.Delete(key)
		return true
	default:
		return false
	}
}

// Wait 等待会话响应
func (s *Session) Wait(timeout time.Duration) (interface{}, error) {
	select {
	case data := <-s.ch:
		return data, nil
	case <-time.After(timeout):
		return nil, errors.New("等待超时")
	}
}
