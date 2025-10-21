package xbot

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter 限流器接口
type Limiter interface {
	Allow(key string) bool
	Reset(key string) error
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	redis    *redis.Client
	duration time.Duration
	maxCount int
	onExceed func(ctx *Context)
}

// NewRateLimiter 创建滑动窗口限流器
func NewRateLimiter(rdb *redis.Client, duration time.Duration, maxCount int, onExceed func(*Context)) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		redis:    rdb,
		duration: duration,
		maxCount: maxCount,
		onExceed: onExceed,
	}
}

// Allow 判断是否允许请求
func (l *SlidingWindowLimiter) Allow(key string) bool {
	if l.redis == nil {
		// 如果没有 Redis，使用内存限流
		return true
	}

	ctx := context.Background()
	now := time.Now()
	windowStart := now.Add(-l.duration)

	// 使用 Lua 脚本保证原子性
	script := `
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local window_start = tonumber(ARGV[2])
		local max_count = tonumber(ARGV[3])
		local ttl = tonumber(ARGV[4])
		
		-- 移除过期的记录
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start)
		
		-- 获取当前窗口内的计数
		local count = redis.call('ZCARD', key)
		
		if count < max_count then
			-- 添加当前请求
			redis.call('ZADD', key, now, now)
			redis.call('EXPIRE', key, ttl)
			return 1
		else
			return 0
		end
	`

	result, err := l.redis.Eval(ctx, script, []string{key},
		now.UnixNano(),
		windowStart.UnixNano(),
		l.maxCount,
		int(l.duration.Seconds())+1,
	).Int()

	if err != nil {
		// 出错时允许请求
		return true
	}

	return result == 1
}

// Reset 重置限流计数
func (l *SlidingWindowLimiter) Reset(key string) error {
	if l.redis == nil {
		return nil
	}

	ctx := context.Background()
	return l.redis.Del(ctx, key).Err()
}

// GetRemaining 获取剩余请求次数
func (l *SlidingWindowLimiter) GetRemaining(key string) int {
	if l.redis == nil {
		return l.maxCount
	}

	ctx := context.Background()
	now := time.Now()
	windowStart := now.Add(-l.duration)

	// 移除过期的记录
	l.redis.ZRemRangeByScore(ctx, key, "0", fmt.Sprint(windowStart.UnixNano()))

	// 获取当前计数
	count, err := l.redis.ZCard(ctx, key).Result()
	if err != nil {
		return l.maxCount
	}

	remaining := l.maxCount - int(count)
	if remaining < 0 {
		return 0
	}

	return remaining
}

// MemoryLimiter 内存限流器（简单实现，适用于单机）
type MemoryLimiter struct {
	records  map[string][]time.Time
	duration time.Duration
	maxCount int
	onExceed func(ctx *Context)
}

// NewMemoryLimiter 创建内存限流器
func NewMemoryLimiter(duration time.Duration, maxCount int, onExceed func(*Context)) *MemoryLimiter {
	limiter := &MemoryLimiter{
		records:  make(map[string][]time.Time),
		duration: duration,
		maxCount: maxCount,
		onExceed: onExceed,
	}

	// 定期清理过期记录
	go limiter.cleanup()

	return limiter
}

// Allow 判断是否允许请求
func (l *MemoryLimiter) Allow(key string) bool {
	now := time.Now()
	windowStart := now.Add(-l.duration)

	// 清理过期记录
	if records, ok := l.records[key]; ok {
		var valid []time.Time
		for _, t := range records {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		l.records[key] = valid
	}

	// 检查是否超过限制
	if len(l.records[key]) >= l.maxCount {
		return false
	}

	// 添加当前请求
	l.records[key] = append(l.records[key], now)
	return true
}

// Reset 重置限流计数
func (l *MemoryLimiter) Reset(key string) error {
	delete(l.records, key)
	return nil
}

// cleanup 清理过期记录
func (l *MemoryLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		windowStart := now.Add(-l.duration)

		for key, records := range l.records {
			var valid []time.Time
			for _, t := range records {
				if t.After(windowStart) {
					valid = append(valid, t)
				}
			}

			if len(valid) == 0 {
				delete(l.records, key)
			} else {
				l.records[key] = valid
			}
		}
	}
}
