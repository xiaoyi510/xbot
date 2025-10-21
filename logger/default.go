package logger

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// 颜色代码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
)

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	mu     sync.Mutex
	out    io.Writer
	level  Level
	fields map[string]interface{}
}

// NewDefaultLogger 创建默认日志记录器
func NewDefaultLogger(out io.Writer, level Level) *DefaultLogger {
	return &DefaultLogger{
		out:    out,
		level:  level,
		fields: make(map[string]interface{}),
	}
}

// SetLevel 设置日志级别
func (l *DefaultLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Debug 记录 Debug 级别日志
func (l *DefaultLogger) Debug(msg string, fields ...interface{}) {
	l.log(LevelDebug, msg, fields...)
}

// Info 记录 Info 级别日志
func (l *DefaultLogger) Info(msg string, fields ...interface{}) {
	l.log(LevelInfo, msg, fields...)
}

// Warn 记录 Warn 级别日志
func (l *DefaultLogger) Warn(msg string, fields ...interface{}) {
	l.log(LevelWarn, msg, fields...)
}

// Error 记录 Error 级别日志
func (l *DefaultLogger) Error(msg string, fields ...interface{}) {
	l.log(LevelError, msg, fields...)
}

// WithField 添加字段
func (l *DefaultLogger) WithField(key string, value interface{}) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &DefaultLogger{
		out:    l.out,
		level:  l.level,
		fields: newFields,
	}
}

// log 记录日志
func (l *DefaultLogger) log(level Level, msg string, fields ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	// 格式化时间
	now := time.Now().Format("2006-01-02 15:04:05")

	// 选择颜色
	color := colorReset
	switch level {
	case LevelDebug:
		color = colorGray
	case LevelInfo:
		color = colorGreen
	case LevelWarn:
		color = colorYellow
	case LevelError:
		color = colorRed
	}

	// 构建日志消息
	logMsg := fmt.Sprintf("%s[%s]%s [%s] %s",
		color, level.String(), colorReset, now, msg)

	// 添加默认字段
	if len(l.fields) > 0 {
		logMsg += " |"
		for k, v := range l.fields {
			logMsg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	// 添加传入的字段
	if len(fields) > 0 {
		if len(l.fields) == 0 {
			logMsg += " |"
		}
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logMsg += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}

	logMsg += "\n"

	// 写入日志
	l.out.Write([]byte(logMsg))
}
