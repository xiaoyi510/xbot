package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xiaoyi510/xbot/event"
	"github.com/xiaoyi510/xbot/logger"
	"github.com/xiaoyi510/xbot/types"
	"github.com/xiaoyi510/xbot/utils"

	"github.com/gorilla/websocket"
)

// WebSocketDriver 正向 WebSocket 驱动器
// 客户端主动连接到 OneBot 实现的 WebSocket 服务器
type WebSocketDriver struct {
	config       Config
	conn         *websocket.Conn
	mu           sync.RWMutex
	eventHandler EventHandler
	apiResponses *utils.SafeMap[string, chan *types.APIResponse]
	connected    bool
	stopChan     chan struct{}
}

// NewWebSocketDriver 创建正向 WebSocket 驱动器
func NewWebSocketDriver(config Config) *WebSocketDriver {
	if config.ReconnectInterval == 0 {
		config.ReconnectInterval = 5
	}
	if config.Timeout == 0 {
		config.Timeout = 30
	}

	return &WebSocketDriver{
		config:       config,
		apiResponses: utils.NewSafeMap[string, chan *types.APIResponse](),
		stopChan:     make(chan struct{}),
	}
}

// Connect 连接到 OneBot 实现
func (d *WebSocketDriver) Connect() error {
	// 构建 WebSocket URL
	wsURL := d.config.URL
	if wsURL == "" {
		wsURL = fmt.Sprintf("ws://%s:%d", d.config.Host, d.config.Port)
	}

	// 设置请求头
	header := make(map[string][]string)
	if d.config.AccessToken != "" {
		header["Authorization"] = []string{"Bearer " + d.config.AccessToken}
	}

	// 连接 WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		return fmt.Errorf("WebSocket 连接失败: %w", err)
	}

	d.mu.Lock()
	d.conn = conn
	d.connected = true
	d.mu.Unlock()

	logger.Info("正向 WebSocket 连接成功", "url", wsURL)

	// 启动消息接收协程
	go d.receiveMessages()

	// 启动心跳协程
	if d.config.HeartbeatInterval > 0 {
		go d.heartbeatLoop()
	}

	// 启动重连协程
	go d.reconnectLoop()

	return nil
}

// receiveMessages 接收消息
func (d *WebSocketDriver) receiveMessages() {
	for {
		d.mu.RLock()
		conn := d.conn
		d.mu.RUnlock()

		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Warn("接收 WebSocket 消息失败", "error", err)
			d.mu.Lock()
			d.connected = false
			d.mu.Unlock()
			return
		}

		// 解析消息
		go d.handleMessage(message)
	}
}

// handleMessage 处理消息
func (d *WebSocketDriver) handleMessage(data []byte) {
	// 判断是事件还是 API 响应
	var base struct {
		Echo string `json:"echo"`
	}

	if err := json.Unmarshal(data, &base); err != nil {
		logger.Error("解析消息失败", "error", err)
		return
	}

	// 如果有 echo，说明是 API 响应
	if base.Echo != "" {
		var resp types.APIResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			logger.Error("解析 API 响应失败", "error", err)
			return
		}

		// 发送到对应的等待通道
		if ch, ok := d.apiResponses.Get(base.Echo); ok {
			select {
			case ch <- &resp:
			case <-time.After(time.Second):
				logger.Warn("API 响应通道阻塞", "echo", base.Echo)
			}
			d.apiResponses.Delete(base.Echo)
		}
		return
	}

	// 否则是事件
	evt, err := event.ParseEvent(data)
	if err != nil {
		logger.Error("解析事件失败", "error", err)
		return
	}

	// 调用事件处理器
	if d.eventHandler != nil {
		d.eventHandler(evt)
	}
}

// CallAPI 调用 OneBot API
func (d *WebSocketDriver) CallAPI(action string, params map[string]interface{}) (*types.APIResponse, error) {
	d.mu.RLock()
	conn := d.conn
	connected := d.connected
	d.mu.RUnlock()

	if !connected || conn == nil {
		return nil, errors.New("WebSocket 未连接")
	}

	// 生成 echo
	echo := utils.GenerateEcho()

	// 创建响应通道
	respChan := make(chan *types.APIResponse, 1)
	d.apiResponses.Set(echo, respChan)

	// 构建请求
	request := map[string]interface{}{
		"action": action,
		"params": params,
		"echo":   echo,
	}

	// 发送请求
	data, err := json.Marshal(request)
	if err != nil {
		d.apiResponses.Delete(echo)
		return nil, err
	}

	d.mu.RLock()
	err = conn.WriteMessage(websocket.TextMessage, data)
	d.mu.RUnlock()

	if err != nil {
		d.apiResponses.Delete(echo)
		return nil, err
	}

	// 等待响应
	timeout := time.Duration(d.config.Timeout) * time.Second
	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(timeout):
		d.apiResponses.Delete(echo)
		return nil, errors.New("API 调用超时")
	}
}

// SetEventHandler 设置事件处理器
func (d *WebSocketDriver) SetEventHandler(handler EventHandler) {
	d.eventHandler = handler
}

// Close 关闭连接
func (d *WebSocketDriver) Close() error {
	close(d.stopChan)

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.conn != nil {
		err := d.conn.Close()
		d.conn = nil
		d.connected = false
		return err
	}

	return nil
}

// IsConnected 是否已连接
func (d *WebSocketDriver) IsConnected() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.connected
}

// heartbeatLoop 心跳循环
func (d *WebSocketDriver) heartbeatLoop() {
	ticker := time.NewTicker(time.Duration(d.config.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopChan:
			return
		case <-ticker.C:
			if d.IsConnected() {
				d.mu.RLock()
				conn := d.conn
				d.mu.RUnlock()

				if conn != nil {
					// 发送 Ping 消息
					if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						logger.Warn("发送心跳失败", "error", err)
						d.mu.Lock()
						d.connected = false
						d.mu.Unlock()
					}
				}
			}
		}
	}
}

// reconnectLoop 重连循环
func (d *WebSocketDriver) reconnectLoop() {
	reconnectCount := 0
	ticker := time.NewTicker(time.Duration(d.config.ReconnectInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.stopChan:
			return
		case <-ticker.C:
			if !d.IsConnected() {
				// 检查是否超过最大重连次数
				if d.config.MaxReconnect > 0 && reconnectCount >= d.config.MaxReconnect {
					logger.Error("达到最大重连次数", "count", reconnectCount)
					return
				}

				logger.Info("尝试重新连接 WebSocket", "attempt", reconnectCount+1)
				if err := d.Connect(); err != nil {
					logger.Error("WebSocket 重连失败", "error", err)
					reconnectCount++
				} else {
					reconnectCount = 0
				}
			}
		}
	}
}
