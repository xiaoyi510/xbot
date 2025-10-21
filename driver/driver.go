package driver

import (
	"xbot/event"
	"xbot/types"
)

// Driver 驱动器接口
type Driver interface {
	// Connect 连接到 OneBot 实现
	Connect() error

	// CallAPI 调用 OneBot API
	CallAPI(action string, params map[string]interface{}) (*types.APIResponse, error)

	// SetEventHandler 设置事件处理器
	SetEventHandler(handler EventHandler)

	// Close 关闭连接
	Close() error

	// IsConnected 是否已连接
	IsConnected() bool
}

// EventHandler 事件处理器
type EventHandler func(event event.Event)

// Config 驱动器配置
type Config struct {
	Type              string // 驱动器类型
	URL               string // 连接 URL
	Host              string // 主机地址
	Port              int    // 端口
	AccessToken       string // 访问令牌
	ReconnectInterval int    // 重连间隔（秒）
	MaxReconnect      int    // 最大重连次数
	HeartbeatInterval int    // 心跳间隔（秒）
	Timeout           int    // API 调用超时（秒）
}
