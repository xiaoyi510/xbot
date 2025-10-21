package xbot

import (
	"sync"
	"testing"
)

// TestDefaultKeywordManager 测试默认关键词管理器
func TestDefaultKeywordManager(t *testing.T) {
	manager := NewKeywordManager([]string{"关键词1", "关键词2"})

	// 测试 GetKeywords
	keywords := manager.GetKeywords()
	if len(keywords) != 2 {
		t.Errorf("Expected 2 keywords, got %d", len(keywords))
	}

	// 测试 GetVersion
	version1 := manager.GetVersion()
	if version1 != 1 {
		t.Errorf("Expected version 1, got %d", version1)
	}

	// 测试 AddKeyword
	manager.AddKeyword("关键词3")
	keywords = manager.GetKeywords()
	if len(keywords) != 3 {
		t.Errorf("Expected 3 keywords after add, got %d", len(keywords))
	}

	version2 := manager.GetVersion()
	if version2 != 2 {
		t.Errorf("Expected version 2 after add, got %d", version2)
	}

	// 测试 RemoveKeyword
	if !manager.RemoveKeyword("关键词1") {
		t.Error("Failed to remove existing keyword")
	}

	keywords = manager.GetKeywords()
	if len(keywords) != 2 {
		t.Errorf("Expected 2 keywords after remove, got %d", len(keywords))
	}

	version3 := manager.GetVersion()
	if version3 != 3 {
		t.Errorf("Expected version 3 after remove, got %d", version3)
	}

	// 测试删除不存在的关键词
	if manager.RemoveKeyword("不存在") {
		t.Error("Should not remove non-existent keyword")
	}

	// 测试 Clear
	manager.Clear()
	keywords = manager.GetKeywords()
	if len(keywords) != 0 {
		t.Errorf("Expected 0 keywords after clear, got %d", len(keywords))
	}

	version4 := manager.GetVersion()
	if version4 != 4 {
		t.Errorf("Expected version 4 after clear, got %d", version4)
	}

	// 测试 Count
	manager.AddKeyword("A")
	manager.AddKeyword("B")
	if manager.Count() != 2 {
		t.Errorf("Expected count 2, got %d", manager.Count())
	}
}

// TestAtomicKeywordManager 测试原子关键词管理器
func TestAtomicKeywordManager(t *testing.T) {
	manager := NewAtomicKeywordManager([]string{"关键词1", "关键词2"})

	// 测试 GetKeywords
	keywords := manager.GetKeywords()
	if len(keywords) != 2 {
		t.Errorf("Expected 2 keywords, got %d", len(keywords))
	}

	// 测试 GetVersion
	version1 := manager.GetVersion()
	if version1 != 1 {
		t.Errorf("Expected version 1, got %d", version1)
	}

	// 测试 AddKeyword
	manager.AddKeyword("关键词3")
	keywords = manager.GetKeywords()
	if len(keywords) != 3 {
		t.Errorf("Expected 3 keywords after add, got %d", len(keywords))
	}

	version2 := manager.GetVersion()
	if version2 != 2 {
		t.Errorf("Expected version 2 after add, got %d", version2)
	}

	// 测试 RemoveKeyword
	if !manager.RemoveKeyword("关键词1") {
		t.Error("Failed to remove existing keyword")
	}

	keywords = manager.GetKeywords()
	if len(keywords) != 2 {
		t.Errorf("Expected 2 keywords after remove, got %d", len(keywords))
	}

	// 测试 Clear
	manager.Clear()
	keywords = manager.GetKeywords()
	if len(keywords) != 0 {
		t.Errorf("Expected 0 keywords after clear, got %d", len(keywords))
	}

	// 测试 Count
	manager.SetKeywords([]string{"A", "B", "C"})
	if manager.Count() != 3 {
		t.Errorf("Expected count 3, got %d", manager.Count())
	}
}

// TestConcurrency 测试并发安全性
func TestConcurrency(t *testing.T) {
	manager := NewKeywordManager([]string{"初始"})

	var wg sync.WaitGroup
	workers := 100
	operations := 10

	// 并发读写
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				// 读取
				_ = manager.GetKeywords()
				_ = manager.GetVersion()

				// 写入
				if j%2 == 0 {
					manager.AddKeyword("关键词")
				} else {
					manager.RemoveKeyword("关键词")
				}
			}
		}(i)
	}

	wg.Wait()

	// 验证版本号递增
	finalVersion := manager.GetVersion()
	expectedMin := int64(1 + workers*operations)
	if finalVersion < expectedMin {
		t.Errorf("Expected version >= %d, got %d", expectedMin, finalVersion)
	}
}

// TestVersionIncrement 测试版本号递增
func TestVersionIncrement(t *testing.T) {
	manager := NewKeywordManager([]string{})

	v1 := manager.GetVersion()
	manager.AddKeyword("A")
	v2 := manager.GetVersion()

	if v2 != v1+1 {
		t.Errorf("Expected version %d, got %d", v1+1, v2)
	}

	manager.SetKeywords([]string{"B", "C"})
	v3 := manager.GetVersion()

	if v3 != v2+1 {
		t.Errorf("Expected version %d, got %d", v2+1, v3)
	}

	manager.RemoveKeyword("B")
	v4 := manager.GetVersion()

	if v4 != v3+1 {
		t.Errorf("Expected version %d, got %d", v3+1, v4)
	}

	manager.Clear()
	v5 := manager.GetVersion()

	if v5 != v4+1 {
		t.Errorf("Expected version %d, got %d", v4+1, v5)
	}
}

// BenchmarkDefaultManager_GetVersion 基准测试：获取版本号
func BenchmarkDefaultManager_GetVersion(b *testing.B) {
	manager := NewKeywordManager([]string{"关键词"})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetVersion()
	}
}

// BenchmarkDefaultManager_GetKeywords 基准测试：获取关键词列表
func BenchmarkDefaultManager_GetKeywords(b *testing.B) {
	keywords := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keywords[i] = "关键词" + string(rune(i))
	}
	manager := NewKeywordManager(keywords)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetKeywords()
	}
}

// BenchmarkAtomicManager_GetVersion 基准测试：原子管理器获取版本号
func BenchmarkAtomicManager_GetVersion(b *testing.B) {
	manager := NewAtomicKeywordManager([]string{"关键词"})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetVersion()
	}
}

// BenchmarkAtomicManager_GetKeywords 基准测试：原子管理器获取关键词
func BenchmarkAtomicManager_GetKeywords(b *testing.B) {
	keywords := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keywords[i] = "关键词" + string(rune(i))
	}
	manager := NewAtomicKeywordManager(keywords)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetKeywords()
	}
}
