package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"xbot/event"
	"xbot/logger"
	"xbot/types"
	"xbot/utils"

	"github.com/gorilla/websocket"
)

// WSReverseDriver 反向 WebSocket 驱动器
type WSReverseDriver struct {
	config       Config
	conn         *websocket.Conn
	mu           sync.RWMutex
	eventHandler EventHandler
	apiResponses *utils.SafeMap[string, chan *types.APIResponse]
	connected    bool
	stopChan     chan struct{}
}

// NewWSReverseDriver 创建反向 WebSocket 驱动器
func NewWSReverseDriver(config Config) *WSReverseDriver {
	if config.ReconnectInterval == 0 {
		config.ReconnectInterval = 5
	}
	if config.Timeout == 0 {
		config.Timeout = 30
	}

	return &WSReverseDriver{
		config:       config,
		apiResponses: utils.NewSafeMap[string, chan *types.APIResponse](),
		stopChan:     make(chan struct{}),
	}
}

// Connect 连接到 OneBot 实现
func (d *WSReverseDriver) Connect() error {
	header := make(map[string][]string)
	if d.config.AccessToken != "" {
		header["Authorization"] = []string{"Bearer " + d.config.AccessToken}
	}

	conn, _, err := websocket.DefaultDialer.Dial(d.config.URL, header)
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}

	d.mu.Lock()
	d.conn = conn
	d.connected = true
	d.mu.Unlock()

	logger.Info("WebSocket 连接成功", "url", d.config.URL)

	// 启动消息接收协程
	go d.receiveMessages()

	// 启动重连协程
	go d.reconnectLoop()

	return nil
}

// receiveMessages 接收消息
func (d *WSReverseDriver) receiveMessages() {
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
			logger.Warn("接收消息失败", "error", err)
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
func (d *WSReverseDriver) handleMessage(data []byte) {
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
func (d *WSReverseDriver) CallAPI(action string, params map[string]interface{}) (*types.APIResponse, error) {
	d.mu.RLock()
	conn := d.conn
	connected := d.connected
	d.mu.RUnlock()

	if !connected || conn == nil {
		return nil, errors.New("未连接")
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

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
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
func (d *WSReverseDriver) SetEventHandler(handler EventHandler) {
	d.eventHandler = handler
}

// Close 关闭连接
func (d *WSReverseDriver) Close() error {
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
func (d *WSReverseDriver) IsConnected() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.connected
}

// reconnectLoop 重连循环
func (d *WSReverseDriver) reconnectLoop() {
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

				logger.Info("尝试重新连接", "attempt", reconnectCount+1)
				if err := d.Connect(); err != nil {
					logger.Error("重连失败", "error", err)
					reconnectCount++
				} else {
					reconnectCount = 0
				}
			}
		}
	}
}
