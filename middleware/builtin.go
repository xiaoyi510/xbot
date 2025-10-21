package middleware

import (
	"fmt"
	"github.com/xiaoyi510/xbot/logger"
	"reflect"
	"runtime/debug"
	"time"
)

// Recovery 恢复中间件 - 捕获 panic
func Recovery() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx interface{}) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(fmt.Sprintf("Panic recovered: %v", err))
					logger.Error(string(debug.Stack()))
				}
			}()
			next(ctx)
		}
	}
}

// Logger 日志中间件 - 记录每个事件的处理及其类型
// 只有当事件匹配到规则时才记录处理完成日志，避免记录无用的日志
func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx interface{}) {
			start := time.Now()
			next(ctx)
			duration := time.Since(start)

			// 使用反射检查是否有 matched 字段，避免循环依赖
			shouldLog := false
			eventType := "未知"

			if ctx != nil {
				v := reflect.ValueOf(ctx)
				if v.Kind() == reflect.Ptr && !v.IsNil() {
					elem := v.Elem()
					if elem.Kind() == reflect.Struct {
						// 检查是否有 matched 字段
						matchedField := elem.FieldByName("matched")
						if matchedField.IsValid() && matchedField.Kind() == reflect.Bool {
							shouldLog = matchedField.Bool()
						} else {
							// 如果没有 matched 字段，保持原有行为
							shouldLog = true
						}

						// 尝试获取 Event 字段来确定事件类型
						eventField := elem.FieldByName("Event")
						if eventField.IsValid() && !eventField.IsNil() {
							eventType = fmt.Sprintf("%T", eventField.Interface())
						} else {
							eventType = fmt.Sprintf("%T", ctx)
						}
					} else {
						// 不是结构体，保持原有行为
						shouldLog = true
						eventType = fmt.Sprintf("%T", ctx)
					}
				}
			}

			// 只有匹配到规则或者是其他类型的上下文时才记录日志
			if shouldLog {
				logger.Debug("事件处理完成", "event_type", eventType, "duration", duration)
			}
		}
	}
}

// MessageLogger 消息日志中间件 - 详细记录消息信息
func MessageLogger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx interface{}) {
			// 这里需要类型断言，但由于中间件接口限制，我们先执行
			next(ctx)
		}
	}
}

// Metrics 性能统计中间件
func Metrics() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx interface{}) {
			start := time.Now()
			next(ctx)
			duration := time.Since(start)

			// 可以在这里收集性能指标
			if duration > time.Second {
				logger.Warn("事件处理耗时较长", "duration", duration)
			}
		}
	}
}

// Timeout 超时中间件
func Timeout(timeout time.Duration) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx interface{}) {
			done := make(chan struct{})

			go func() {
				next(ctx)
				close(done)
			}()

			select {
			case <-done:
				// 正常完成
			case <-time.After(timeout):
				logger.Warn("事件处理超时", "timeout", timeout)
			}
		}
	}
}

// Concurrency 并发控制中间件
func Concurrency(max int) Middleware {
	sem := make(chan struct{}, max)

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx interface{}) {
			sem <- struct{}{}
			defer func() { <-sem }()
			next(ctx)
		}
	}
}
