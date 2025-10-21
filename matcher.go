package xbot

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/dlclark/regexp2"
)

// Matcher 匹配器
type Matcher struct {
	priority    int
	filters     []Filter
	limiter     Limiter
	handler     interface{}
	middlewares []func(next func(*Context)) func(*Context)
	matchFunc   func(*Context) bool
	block       bool // 是否阻止继续匹配，默认为false（继续匹配）
}

// newMatcher 创建匹配器
func newMatcher(matchFunc func(*Context) bool, filters ...Filter) *Matcher {
	return &Matcher{
		priority:    0,
		filters:     filters,
		matchFunc:   matchFunc,
		middlewares: make([]func(next func(*Context)) func(*Context), 0),
		block:       false, // 默认继续匹配
	}
}

// Filter 添加过滤器
func (m *Matcher) Filter(filters ...Filter) *Matcher {
	m.filters = append(m.filters, filters...)
	return m
}

// Limit 设置限流
func (m *Matcher) Limit(duration time.Duration, count int, onExceed func(*Context)) *Matcher {
	// 默认使用内存限流器
	m.limiter = NewMemoryLimiter(duration, count, onExceed)
	return m
}

// LimitWithRedis 使用 Redis 限流
func (m *Matcher) LimitWithRedis(limiter Limiter) *Matcher {
	m.limiter = limiter
	return m
}

// Handle 设置处理函数
func (m *Matcher) Handle(handler interface{}) *Matcher {
	m.handler = handler
	return m
}

// Priority 设置优先级
func (m *Matcher) Priority(p int) *Matcher {
	m.priority = p
	return m
}

// Use 添加中间件
func (m *Matcher) Use(middlewares ...func(next func(*Context)) func(*Context)) *Matcher {
	m.middlewares = append(m.middlewares, middlewares...)
	return m
}

// SetBlock 设置是否阻止继续匹配，默认为false（继续匹配）
// 无参数调用时默认设置为true，阻止继续匹配下一个处理器
func (m *Matcher) SetBlock(block ...bool) *Matcher {
	if len(block) == 0 {
		m.block = true // 无参数时默认为true
	} else {
		m.block = block[0]
	}
	return m
}

// Match 判断是否匹配
func (m *Matcher) Match(ctx *Context) bool {
	// 先检查匹配函数
	if m.matchFunc != nil && !m.matchFunc(ctx) {
		return false
	}

	// 再检查过滤器
	for _, filter := range m.filters {
		if !filter(ctx) {
			return false
		}
	}

	// 检查限流
	if m.limiter != nil {
		userID := ctx.GetUserID()
		groupID := ctx.GetGroupID()
		key := generateLimiterKey(userID, groupID)

		if !m.limiter.Allow(key) {
			// 触发限流回调
			if limiter, ok := m.limiter.(*SlidingWindowLimiter); ok && limiter.onExceed != nil {
				limiter.onExceed(ctx)
			} else if limiter, ok := m.limiter.(*MemoryLimiter); ok && limiter.onExceed != nil {
				limiter.onExceed(ctx)
			}
			return false
		}
	}

	return true
}

// Execute 执行处理函数
func (m *Matcher) Execute(ctx *Context) {
	if m.handler == nil {
		return
	}

	// 应用中间件
	handler := func(ctx *Context) {
		// 使用反射调用处理函数
		handlerValue := reflect.ValueOf(m.handler)
		handlerValue.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}

	// 从后向前应用中间件
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		handler = m.middlewares[i](handler)
	}

	// 执行
	handler(ctx)
}

// generateLimiterKey 生成限流 key
func generateLimiterKey(userID, groupID int64) string {
	if groupID == 0 {
		return fmt.Sprintf("limiter:user:%d", userID)
	}
	return fmt.Sprintf("limiter:group:%d:user:%d", groupID, userID)
}

// 辅助函数：创建命令匹配器
func commandMatcher(command string, prefix string) func(*Context) bool {
	return func(ctx *Context) bool {
		text := ctx.GetPlainText()
		if !strings.HasPrefix(text, prefix) {
			return false
		}
		text = strings.TrimPrefix(text, prefix)
		text = strings.TrimSpace(text)

		// 检查命令
		if text == command || strings.HasPrefix(text, command+" ") {
			return true
		}
		return false
	}
}

// 辅助函数：创建命令组匹配器
func commandGroupMatcher(commands []string, prefix string) func(*Context) bool {
	return func(ctx *Context) bool {
		text := ctx.GetPlainText()
		if !strings.HasPrefix(text, prefix) {
			return false
		}
		text = strings.TrimPrefix(text, prefix)
		text = strings.TrimSpace(text)

		// 检查任一命令
		for _, command := range commands {
			if text == command || strings.HasPrefix(text, command+" ") {
				return true
			}
		}
		return false
	}
}

// 辅助函数：创建关键词匹配器
func keywordsMatcher(keywords []string) func(*Context) bool {
	return func(ctx *Context) bool {
		text := ctx.GetPlainText()
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				return true
			}
		}
		return false
	}
}

// 辅助函数：创建正则匹配器（使用 regexp2 支持高级特性）
func regexMatcher(pattern string) func(*Context) bool {
	regex, err := regexp2.Compile(pattern, 0)
	if err != nil {
		// 如果编译失败，返回永不匹配的函数
		return func(ctx *Context) bool {
			return false
		}
	}

	return func(ctx *Context) bool {
		text := ctx.GetPlainText()

		// 使用 regexp2 进行匹配
		match, err := regex.FindStringMatch(text)
		if err != nil || match == nil {
			return false
		}

		// 存储匹配结果到 context
		regexResult := &RegexMatch{
			Match:       match.String(),
			Groups:      make([]string, 0),
			NamedGroups: make(map[string]string),
		}

		// 提取所有分组（包括完整匹配）
		groups := match.Groups()
		for i, group := range groups {
			if len(group.Captures) > 0 {
				captureValue := group.Captures[0].String()
				regexResult.Groups = append(regexResult.Groups, captureValue)

				// 如果有命名分组，添加到 NamedGroups
				if group.Name != "" && i > 0 { // 跳过索引 0 的完整匹配
					regexResult.NamedGroups[group.Name] = captureValue
				}
			}
		}

		// 将结果存储到 context
		ctx.RegexResult = regexResult

		return true
	}
}

// 辅助函数：创建前缀匹配器
func prefixMatcher(prefix string) func(*Context) bool {
	return func(ctx *Context) bool {
		text := ctx.GetPlainText()
		return strings.HasPrefix(text, prefix)
	}
}

// 辅助函数：创建后缀匹配器
func suffixMatcher(suffix string) func(*Context) bool {
	return func(ctx *Context) bool {
		text := ctx.GetPlainText()
		return strings.HasSuffix(text, suffix)
	}
}

// 辅助函数：创建完全匹配器
func fullMatchMatcher(text string) func(*Context) bool {
	return func(ctx *Context) bool {
		return ctx.GetPlainText() == text
	}
}

// 辅助函数：创建 DFA 匹配器（高性能版本，基于版本号）
// provider 必须实现 VersionedKeywordProvider 接口
func dfaMatcher(provider VersionedKeywordProvider) func(*Context) bool {
	// DFA 树节点
	type dfaNode struct {
		children map[rune]*dfaNode
		isEnd    bool
	}

	// DFA 缓存结构
	type dfaCache struct {
		root        *dfaNode
		lastVersion int64
		mu          sync.RWMutex
	}

	cache := &dfaCache{
		lastVersion: -1,
	}

	// 构建 DFA 树
	buildDFA := func(words []string) *dfaNode {
		root := &dfaNode{children: make(map[rune]*dfaNode)}
		for _, word := range words {
			if word == "" {
				continue
			}
			node := root
			for _, char := range []rune(word) {
				if _, ok := node.children[char]; !ok {
					node.children[char] = &dfaNode{children: make(map[rune]*dfaNode)}
				}
				node = node.children[char]
			}
			node.isEnd = true
		}
		return root
	}

	// DFA 匹配
	matchDFA := func(root *dfaNode, text string) bool {
		if root == nil {
			return false
		}
		runes := []rune(text)
		for i := 0; i < len(runes); i++ {
			node := root
			j := i
			for j < len(runes) {
				if child, ok := node.children[runes[j]]; ok {
					node = child
					j++
					if node.isEnd {
						return true
					}
				} else {
					break
				}
			}
		}
		return false
	}

	return func(ctx *Context) bool {
		text := ctx.GetPlainText()

		// 获取版本号（O(1) 操作）
		version := provider.GetVersion()

		// 检查是否需要更新 DFA 树
		cache.mu.RLock()
		needUpdate := cache.root == nil || cache.lastVersion != version
		currentRoot := cache.root
		cache.mu.RUnlock()

		// 需要更新时重建 DFA 树
		if needUpdate {
			words := provider.GetKeywords()
			if len(words) == 0 {
				return false
			}

			cache.mu.Lock()
			// 双重检查，避免多个 goroutine 重复构建
			if cache.root == nil || cache.lastVersion != version {
				cache.root = buildDFA(words)
				cache.lastVersion = version
				currentRoot = cache.root
			} else {
				currentRoot = cache.root
			}
			cache.mu.Unlock()
		}

		// 使用缓存的 DFA 树进行匹配
		return matchDFA(currentRoot, text)
	}
}

// 辅助函数：创建 AC 自动机匹配器（高性能版本，基于版本号）
// provider 必须实现 VersionedKeywordProvider 接口
func acMatcher(provider VersionedKeywordProvider) func(*Context) bool {
	// AC自动机缓存结构
	type acCache struct {
		machine     *ACMachine
		lastVersion int64
		mu          sync.RWMutex
	}

	cache := &acCache{
		machine:     NewACMachine(),
		lastVersion: -1,
	}

	return func(ctx *Context) bool {
		text := ctx.GetPlainText()

		// 获取版本号（O(1) 操作）
		version := provider.GetVersion()

		// 检查是否需要更新 AC 自动机
		cache.mu.RLock()
		needUpdate := cache.lastVersion != version
		cache.mu.RUnlock()

		// 需要更新时重建 AC 自动机
		if needUpdate {
			words := provider.GetKeywords()
			if len(words) == 0 {
				return false
			}

			cache.mu.Lock()
			// 双重检查，避免多个 goroutine 重复构建
			if cache.lastVersion != version {
				cache.machine.Build(words)
				cache.lastVersion = version
			}
			cache.mu.Unlock()
		}

		// 使用 AC 自动机进行匹配
		return cache.machine.Match(text)
	}
}

// 辅助函数：创建 AC 自动机动态匹配器（支持 ContextKeywordProvider）
// provider 根据上下文返回不同的关键词提供者
func acMatcherWithContext(provider ContextKeywordProvider) func(*Context) bool {
	// 为每个唯一的 provider 维护独立的 AC 自动机缓存
	type providerCache struct {
		machine     *ACMachine
		lastVersion int64
		mu          sync.RWMutex
	}

	cacheMap := sync.Map{} // map[VersionedKeywordProvider]*providerCache

	return func(ctx *Context) bool {
		text := ctx.GetPlainText()

		// 根据上下文获取对应的关键词提供者
		versionedProvider := provider.GetProvider(ctx)
		if versionedProvider == nil {
			return false
		}

		// 获取或创建该 provider 的缓存
		cacheInterface, _ := cacheMap.LoadOrStore(versionedProvider, &providerCache{
			machine:     NewACMachine(),
			lastVersion: -1,
		})
		cache := cacheInterface.(*providerCache)

		// 获取版本号
		version := versionedProvider.GetVersion()

		// 检查是否需要更新 AC 自动机
		cache.mu.RLock()
		needUpdate := cache.lastVersion != version
		cache.mu.RUnlock()

		// 需要更新时重建 AC 自动机
		if needUpdate {
			words := versionedProvider.GetKeywords()
			if len(words) == 0 {
				return false
			}

			cache.mu.Lock()
			// 双重检查
			if cache.lastVersion != version {
				cache.machine.Build(words)
				cache.lastVersion = version
			}
			cache.mu.Unlock()
		}

		// 使用 AC 自动机进行匹配
		return cache.machine.Match(text)
	}
}

// 辅助函数：创建 DFA 动态匹配器（支持 ContextKeywordProvider）
// provider 根据上下文返回不同的关键词提供者
func dfaMatcherWithContext(provider ContextKeywordProvider) func(*Context) bool {
	// DFA 树节点
	type dfaNode struct {
		children map[rune]*dfaNode
		isEnd    bool
	}

	// 为每个唯一的 provider 维护独立的 DFA 缓存
	type providerCache struct {
		root        *dfaNode
		lastVersion int64
		mu          sync.RWMutex
	}

	cacheMap := sync.Map{} // map[VersionedKeywordProvider]*providerCache

	// 构建 DFA 树
	buildDFA := func(words []string) *dfaNode {
		root := &dfaNode{children: make(map[rune]*dfaNode)}
		for _, word := range words {
			if word == "" {
				continue
			}
			node := root
			for _, char := range []rune(word) {
				if _, ok := node.children[char]; !ok {
					node.children[char] = &dfaNode{children: make(map[rune]*dfaNode)}
				}
				node = node.children[char]
			}
			node.isEnd = true
		}
		return root
	}

	// DFA 匹配
	matchDFA := func(root *dfaNode, text string) bool {
		if root == nil {
			return false
		}
		runes := []rune(text)
		for i := 0; i < len(runes); i++ {
			node := root
			j := i
			for j < len(runes) {
				if child, ok := node.children[runes[j]]; ok {
					node = child
					j++
					if node.isEnd {
						return true
					}
				} else {
					break
				}
			}
		}
		return false
	}

	return func(ctx *Context) bool {
		text := ctx.GetPlainText()

		// 根据上下文获取对应的关键词提供者
		versionedProvider := provider.GetProvider(ctx)
		if versionedProvider == nil {
			return false
		}

		// 获取或创建该 provider 的缓存
		cacheInterface, _ := cacheMap.LoadOrStore(versionedProvider, &providerCache{
			root:        nil,
			lastVersion: -1,
		})
		cache := cacheInterface.(*providerCache)

		// 获取版本号
		version := versionedProvider.GetVersion()

		// 检查是否需要更新 DFA 树
		cache.mu.RLock()
		needUpdate := cache.root == nil || cache.lastVersion != version
		currentRoot := cache.root
		cache.mu.RUnlock()

		// 需要更新时重建 DFA 树
		if needUpdate {
			words := versionedProvider.GetKeywords()
			if len(words) == 0 {
				return false
			}

			cache.mu.Lock()
			// 双重检查
			if cache.root == nil || cache.lastVersion != version {
				cache.root = buildDFA(words)
				cache.lastVersion = version
				currentRoot = cache.root
			} else {
				currentRoot = cache.root
			}
			cache.mu.Unlock()
		}

		// 使用缓存的 DFA 树进行匹配
		return matchDFA(currentRoot, text)
	}
}
