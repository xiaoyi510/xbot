package storage

import (
	"strings"
	"sync"
)

// MemoryStorage 内存存储实现
type MemoryStorage struct {
	data sync.Map
}

// NewMemoryStorage 创建内存存储
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

// Get 获取值
func (m *MemoryStorage) Get(key string) ([]byte, error) {
	if value, ok := m.data.Load(key); ok {
		return value.([]byte), nil
	}
	return nil, nil
}

// Set 设置值
func (m *MemoryStorage) Set(key string, value []byte) error {
	// 复制数据以避免外部修改
	data := make([]byte, len(value))
	copy(data, value)
	m.data.Store(key, data)
	return nil
}

// Delete 删除值
func (m *MemoryStorage) Delete(key string) error {
	m.data.Delete(key)
	return nil
}

// Has 判断是否存在
func (m *MemoryStorage) Has(key string) bool {
	_, ok := m.data.Load(key)
	return ok
}

// Keys 获取所有以 prefix 开头的 key
func (m *MemoryStorage) Keys(prefix string) ([]string, error) {
	var keys []string

	m.data.Range(func(key, value interface{}) bool {
		k := key.(string)
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
		return true
	})

	return keys, nil
}

// Close 关闭存储（内存存储无需关闭）
func (m *MemoryStorage) Close() error {
	return nil
}
