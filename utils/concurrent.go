package utils

import (
	"sync"
)

// SafeMap 并发安全的 Map
type SafeMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewSafeMap 创建并发安全的 Map
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: make(map[K]V),
	}
}

// Set 设置值
func (sm *SafeMap[K, V]) Set(key K, value V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

// Get 获取值
func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	v, ok := sm.m[key]
	return v, ok
}

// Delete 删除值
func (sm *SafeMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
}

// Has 判断是否存在
func (sm *SafeMap[K, V]) Has(key K) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.m[key]
	return ok
}

// Len 返回长度
func (sm *SafeMap[K, V]) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.m)
}

// Keys 返回所有 key
func (sm *SafeMap[K, V]) Keys() []K {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	keys := make([]K, 0, len(sm.m))
	for k := range sm.m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有 value
func (sm *SafeMap[K, V]) Values() []V {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	values := make([]V, 0, len(sm.m))
	for _, v := range sm.m {
		values = append(values, v)
	}
	return values
}

// Range 遍历
func (sm *SafeMap[K, V]) Range(f func(key K, value V) bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for k, v := range sm.m {
		if !f(k, v) {
			break
		}
	}
}

// Clear 清空
func (sm *SafeMap[K, V]) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m = make(map[K]V)
}
