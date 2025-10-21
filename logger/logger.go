package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Level 日志级别
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String 返回日志级别字符串
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel 解析日志级别
func ParseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return LevelDebug
	case "info", "INFO":
		return LevelInfo
	case "warn", "WARN", "warning", "WARNING":
		return LevelWarn
	case "error", "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	WithField(key string, value interface{}) Logger
	SetLevel(level Level)
}

var (
	defaultLogger Logger
	mu            sync.RWMutex
)

func init() {
	defaultLogger = NewDefaultLogger(os.Stdout, LevelInfo)
}

// SetDefault 设置默认日志记录器
func SetDefault(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	defaultLogger = l
}

// GetDefault 获取默认日志记录器
func GetDefault() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return defaultLogger
}

// Debug 记录 Debug 级别日志
func Debug(msg string, fields ...interface{}) {
	GetDefault().Debug(msg, fields...)
}

// Info 记录 Info 级别日志
func Info(msg string, fields ...interface{}) {
	GetDefault().Info(msg, fields...)
}

// Warn 记录 Warn 级别日志
func Warn(msg string, fields ...interface{}) {
	GetDefault().Warn(msg, fields...)
}

// Error 记录 Error 级别日志
func Error(msg string, fields ...interface{}) {
	GetDefault().Error(msg, fields...)
}

// WithField 添加字段
func WithField(key string, value interface{}) Logger {
	return GetDefault().WithField(key, value)
}

// SetLevel 设置日志级别
func SetLevel(level Level) {
	GetDefault().SetLevel(level)
}

// Debugf 格式化 Debug 日志
func Debugf(format string, args ...interface{}) {
	Debug(fmt.Sprintf(format, args...))
}

// Infof 格式化 Info 日志
func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

// Warnf 格式化 Warn 日志
func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...))
}

// Errorf 格式化 Error 日志
func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}

// NewMultiWriter 创建多输出日志记录器
func NewMultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}

// MustCreateFile 创建日志文件，失败则 panic
func MustCreateFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法创建日志文件: %v", err)
	}
	return file
}
