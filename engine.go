package xbot

import (
	"sort"
	"sync"
	"time"

	"github.com/xiaoyi510/xbot/event"
	"github.com/xiaoyi510/xbot/logger"
	"github.com/xiaoyi510/xbot/middleware"
)

// Engine 引擎
type Engine struct {
	matchers    []*Matcher
	mu          sync.RWMutex
	bot         *Bot
	middlewares []func(next func(*Context)) func(*Context)
}

// NewEngine 创建引擎
func NewEngine() *Engine {
	engine := &Engine{
		matchers:    make([]*Matcher, 0),
		middlewares: make([]func(next func(*Context)) func(*Context), 0),
	}

	// 自动注册到全局
	RegisterEngine(engine)

	return engine
}

// Use 添加全局中间件
func (e *Engine) Use(middlewares ...func(next func(*Context)) func(*Context)) *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.middlewares = append(e.middlewares, middlewares...)
	return e
}

// OnCommand 命令匹配
func (e *Engine) OnCommand(command string, filters ...Filter) *Matcher {
	prefix := ""
	if e.bot != nil {
		prefix = e.bot.Config.CommandPrefix
	} else {
		prefix = "/" // 默认前缀
	}

	matcher := newMatcher(commandMatcher(command, prefix), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnCommandGroup 命令组匹配
func (e *Engine) OnCommandGroup(commands []string, filters ...Filter) *Matcher {
	prefix := ""
	if e.bot != nil {
		prefix = e.bot.Config.CommandPrefix
	} else {
		prefix = "/"
	}

	matcher := newMatcher(commandGroupMatcher(commands, prefix), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnKeywords 关键词匹配
func (e *Engine) OnKeywords(keywords []string, filters ...Filter) *Matcher {
	matcher := newMatcher(keywordsMatcher(keywords), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnRegex 正则表达式匹配
func (e *Engine) OnRegex(pattern string, filters ...Filter) *Matcher {
	matcher := newMatcher(regexMatcher(pattern), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnPrefix 前缀匹配
func (e *Engine) OnPrefix(prefix string, filters ...Filter) *Matcher {
	matcher := newMatcher(prefixMatcher(prefix), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnSuffix 后缀匹配
func (e *Engine) OnSuffix(suffix string, filters ...Filter) *Matcher {
	matcher := newMatcher(suffixMatcher(suffix), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnFullMatch 完全匹配
func (e *Engine) OnFullMatch(text string, filters ...Filter) *Matcher {
	matcher := newMatcher(fullMatchMatcher(text), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnDFAKeywords DFA 关键词匹配器（支持动态更新关键词列表）
// provider 必须实现 VersionedKeywordProvider 接口
//
// 使用版本号检测变化，O(1) 性能，零内存分配
//
// 推荐使用 NewKeywordManager 或 NewAtomicKeywordManager 创建关键词管理器
func (e *Engine) OnDFAKeywords(provider VersionedKeywordProvider, filters ...Filter) *Matcher {
	matcher := newMatcher(dfaMatcher(provider), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnDFAKeywordsWithContext DFA 动态关键词匹配器（支持根据上下文选择关键词库）
// provider 必须实现 ContextKeywordProvider 接口
//
// 根据上下文（如群ID）动态选择不同的关键词库
// 适用场景：不同群组使用不同的敏感词库
func (e *Engine) OnDFAKeywordsWithContext(provider ContextKeywordProvider, filters ...Filter) *Matcher {
	matcher := newMatcher(dfaMatcherWithContext(provider), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnACKeywords AC自动机关键词匹配器（支持动态更新关键词列表）
// provider 必须实现 VersionedKeywordProvider 接口
//
// AC自动机（Aho-Corasick算法）比DFA性能更好，时间复杂度 O(n)
// 推荐在关键词数量较多（>100）时使用
//
// 使用版本号检测变化，O(1) 性能，零内存分配
func (e *Engine) OnACKeywords(provider VersionedKeywordProvider, filters ...Filter) *Matcher {
	matcher := newMatcher(acMatcher(provider), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnACKeywordsWithContext AC自动机动态关键词匹配器（支持根据上下文选择关键词库）
// provider 必须实现 ContextKeywordProvider 接口
//
// 结合AC自动机的高性能和动态Provider的灵活性
// 根据上下文（如群ID）动态选择不同的关键词库
// 适用场景：不同群组使用不同的敏感词库，且关键词数量较多
func (e *Engine) OnACKeywordsWithContext(provider ContextKeywordProvider, filters ...Filter) *Matcher {
	matcher := newMatcher(acMatcherWithContext(provider), filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnMessage 消息事件
func (e *Engine) OnMessage(filters ...Filter) *Matcher {
	matcher := newMatcher(func(ctx *Context) bool {
		_, isPrivate := ctx.Event.(*event.PrivateMessageEvent)
		_, isGroup := ctx.Event.(*event.GroupMessageEvent)
		return isPrivate || isGroup
	}, filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnNotice 通知事件
func (e *Engine) OnNotice(filters ...Filter) *Matcher {
	matcher := newMatcher(func(ctx *Context) bool {
		return ctx.Event.GetPostType() == "notice"
	}, filters...)
	e.addMatcher(matcher)
	return matcher
}

// OnRequest 请求事件
func (e *Engine) OnRequest(filters ...Filter) *Matcher {
	matcher := newMatcher(func(ctx *Context) bool {
		return ctx.Event.GetPostType() == "request"
	}, filters...)
	e.addMatcher(matcher)
	return matcher
}

// addMatcher 添加匹配器
func (e *Engine) addMatcher(matcher *Matcher) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.matchers = append(e.matchers, matcher)

	// 按优先级排序
	sort.Slice(e.matchers, func(i, j int) bool {
		return e.matchers[i].priority > e.matchers[j].priority
	})
}

// HandleEvent 处理事件
func (e *Engine) HandleEvent(evt event.Event) {
	// 创建上下文
	ctx := NewContext(evt, e.bot)

	// 使用中间件包装处理流程
	handler := func(ctx *Context) {
		e.handleEventInternal(ctx)
	}

	// 应用引擎级中间件
	e.mu.RLock()
	mws := e.middlewares
	e.mu.RUnlock()

	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}

	// 执行
	handler(ctx)
}

// handleEventInternal 内部事件处理
func (e *Engine) handleEventInternal(ctx *Context) {
	e.mu.RLock()
	matchers := make([]*Matcher, len(e.matchers))
	copy(matchers, e.matchers)
	e.mu.RUnlock()

	// 遍历匹配器
	for _, matcher := range matchers {
		// 检查是否已中止
		if ctx.IsAborted() {
			break
		}

		if matcher.Match(ctx) {
			// 使用 goroutine 处理事件，避免阻塞
			go func(m *Matcher) {
				defer func() {
					if err := recover(); err != nil {
						logger.Error("处理事件时发生错误", "error", err)
					}
				}()

				m.Execute(ctx)
			}(matcher)

			// 标记已匹配
			ctx.matched = true

			// 检查是否应该阻止继续匹配
			if matcher.block {
				break
			}
		}
	}
}

// SetBot 设置 Bot
func (e *Engine) SetBot(bot *Bot) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.bot = bot
}

// 全局引擎注册
var (
	globalEngines []*Engine
	engineMu      sync.Mutex
)

// RegisterEngine 注册引擎
func RegisterEngine(engine *Engine) {
	engineMu.Lock()
	defer engineMu.Unlock()
	globalEngines = append(globalEngines, engine)
}

// GetEngines 获取所有引擎
func GetEngines() []*Engine {
	engineMu.Lock()
	defer engineMu.Unlock()
	engines := make([]*Engine, len(globalEngines))
	copy(engines, globalEngines)
	return engines
}

// 添加默认中间件的便捷方法
func (e *Engine) UseRecovery() *Engine {
	return e.Use(func(next func(*Context)) func(*Context) {
		return func(ctx *Context) {
			middleware.Recovery()(func(c interface{}) {
				next(c.(*Context))
			})(ctx)
		}
	})
}

func (e *Engine) UseLogger() *Engine {
	return e.Use(func(next func(*Context)) func(*Context) {
		return func(ctx *Context) {
			middleware.Logger()(func(c interface{}) {
				next(c.(*Context))
			})(ctx)
		}
	})
}

func (e *Engine) UseMetrics() *Engine {
	return e.Use(func(next func(*Context)) func(*Context) {
		return func(ctx *Context) {
			middleware.Metrics()(func(c interface{}) {
				next(c.(*Context))
			})(ctx)
		}
	})
}

// formatTimestamp 格式化时间戳
func formatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}
