package xbot

import (
	"sync"
	"sync/atomic"
)

// KeywordProvider 关键词提供者接口（基础）
type KeywordProvider interface {
	GetKeywords() []string
}

// VersionedKeywordProvider 带版本号的关键词提供者（高性能）
type VersionedKeywordProvider interface {
	KeywordProvider
	GetVersion() int64
}

// ContextKeywordProvider 根据上下文返回关键词提供者
// 用于根据不同的群组或场景返回不同的关键词库
type ContextKeywordProvider interface {
	GetProvider(ctx *Context) VersionedKeywordProvider
}

// KeywordManager 完整的关键词管理器接口
type KeywordManager interface {
	VersionedKeywordProvider
	SetKeywords([]string)
	AddKeyword(string)
	RemoveKeyword(string) bool
	Clear()
	Count() int
}

// DefaultKeywordManager 默认的关键词管理器实现
type DefaultKeywordManager struct {
	keywords []string
	version  int64
	mu       sync.RWMutex
}

// NewKeywordManager 创建新的关键词管理器
func NewKeywordManager(keywords []string) *DefaultKeywordManager {
	return &DefaultKeywordManager{
		keywords: keywords,
		version:  1,
	}
}

// GetKeywords 获取关键词列表（返回副本，线程安全）
func (m *DefaultKeywordManager) GetKeywords() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]string, len(m.keywords))
	copy(result, m.keywords)
	return result
}

// GetVersion 获取版本号（O(1) 操作，高性能）
func (m *DefaultKeywordManager) GetVersion() int64 {
	return atomic.LoadInt64(&m.version)
}

// SetKeywords 设置关键词列表
func (m *DefaultKeywordManager) SetKeywords(keywords []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keywords = keywords
	atomic.AddInt64(&m.version, 1)
}

// AddKeyword 添加关键词
func (m *DefaultKeywordManager) AddKeyword(keyword string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keywords = append(m.keywords, keyword)
	atomic.AddInt64(&m.version, 1)
}

// RemoveKeyword 删除关键词
func (m *DefaultKeywordManager) RemoveKeyword(keyword string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, kw := range m.keywords {
		if kw == keyword {
			m.keywords = append(m.keywords[:i], m.keywords[i+1:]...)
			atomic.AddInt64(&m.version, 1)
			return true
		}
	}
	return false
}

// Clear 清空所有关键词
func (m *DefaultKeywordManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keywords = []string{}
	atomic.AddInt64(&m.version, 1)
}

// Count 获取关键词数量
func (m *DefaultKeywordManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.keywords)
}

// AtomicKeywordManager 使用 atomic.Value 的高性能实现（零锁读取）
type AtomicKeywordManager struct {
	keywords atomic.Value // []string
	version  int64
}

// NewAtomicKeywordManager 创建高性能关键词管理器
func NewAtomicKeywordManager(keywords []string) *AtomicKeywordManager {
	m := &AtomicKeywordManager{
		version: 1,
	}
	m.keywords.Store(keywords)
	return m
}

// GetKeywords 获取关键词列表（无锁读取）
func (m *AtomicKeywordManager) GetKeywords() []string {
	if kw := m.keywords.Load(); kw != nil {
		return kw.([]string)
	}
	return []string{}
}

// GetVersion 获取版本号
func (m *AtomicKeywordManager) GetVersion() int64 {
	return atomic.LoadInt64(&m.version)
}

// SetKeywords 设置关键词列表
func (m *AtomicKeywordManager) SetKeywords(keywords []string) {
	// 创建副本，避免外部修改
	kwCopy := make([]string, len(keywords))
	copy(kwCopy, keywords)
	m.keywords.Store(kwCopy)
	atomic.AddInt64(&m.version, 1)
}

// AddKeyword 添加关键词（copy-on-write）
func (m *AtomicKeywordManager) AddKeyword(keyword string) {
	current := m.GetKeywords()
	newKeywords := make([]string, len(current)+1)
	copy(newKeywords, current)
	newKeywords[len(current)] = keyword
	m.keywords.Store(newKeywords)
	atomic.AddInt64(&m.version, 1)
}

// RemoveKeyword 删除关键词（copy-on-write）
func (m *AtomicKeywordManager) RemoveKeyword(keyword string) bool {
	current := m.GetKeywords()
	for i, kw := range current {
		if kw == keyword {
			newKeywords := make([]string, 0, len(current)-1)
			newKeywords = append(newKeywords, current[:i]...)
			newKeywords = append(newKeywords, current[i+1:]...)
			m.keywords.Store(newKeywords)
			atomic.AddInt64(&m.version, 1)
			return true
		}
	}
	return false
}

// Clear 清空所有关键词
func (m *AtomicKeywordManager) Clear() {
	m.keywords.Store([]string{})
	atomic.AddInt64(&m.version, 1)
}

// Count 获取关键词数量
func (m *AtomicKeywordManager) Count() int {
	return len(m.GetKeywords())
}
