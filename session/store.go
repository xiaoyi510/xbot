package session

import (
	"sync"
	"time"
)

// Store 会话存储接口
type Store interface {
	Get(key string) (*Session, error)
	Set(key string, session *Session, ttl time.Duration) error
	Delete(key string) error
}

// MemoryStore 内存存储实现
type MemoryStore struct {
	data sync.Map
}

type memoryItem struct {
	session   *Session
	expiresAt time.Time
}

// NewMemoryStore 创建内存存储
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{}

	// 启动清理协程
	go store.cleanup()

	return store
}

// Get 获取会话
func (s *MemoryStore) Get(key string) (*Session, error) {
	value, ok := s.data.Load(key)
	if !ok {
		return nil, nil
	}

	item := value.(*memoryItem)

	// 检查是否过期
	if time.Now().After(item.expiresAt) {
		s.data.Delete(key)
		return nil, nil
	}

	return item.session, nil
}

// Set 设置会话
func (s *MemoryStore) Set(key string, session *Session, ttl time.Duration) error {
	item := &memoryItem{
		session:   session,
		expiresAt: time.Now().Add(ttl),
	}

	s.data.Store(key, item)
	return nil
}

// Delete 删除会话
func (s *MemoryStore) Delete(key string) error {
	s.data.Delete(key)
	return nil
}

// cleanup 清理过期会话
func (s *MemoryStore) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		s.data.Range(func(key, value interface{}) bool {
			item := value.(*memoryItem)
			if now.After(item.expiresAt) {
				s.data.Delete(key)
			}
			return true
		})
	}
}
