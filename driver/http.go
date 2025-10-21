package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"xbot/logger"
	"xbot/types"
)

// HTTPDriver HTTP 驱动器
// 通过 HTTP 请求调用 OneBot API，不接收事件推送
// 适用于只需要调用 API 不需要接收事件的场景
type HTTPDriver struct {
	config       Config
	client       *http.Client
	mu           sync.RWMutex
	eventHandler EventHandler
	connected    bool
}

// NewHTTPDriver 创建 HTTP 驱动器
func NewHTTPDriver(config Config) *HTTPDriver {
	if config.Timeout == 0 {
		config.Timeout = 30
	}

	return &HTTPDriver{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		connected: false,
	}
}

// Connect 连接到 OneBot 实现（HTTP 驱动器不需要持久连接）
func (d *HTTPDriver) Connect() error {
	// HTTP 驱动器不需要建立持久连接，测试一下连通性即可
	baseURL := d.config.URL
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://%s:%d", d.config.Host, d.config.Port)
	}

	// 尝试调用 get_version 测试连接
	_, err := d.CallAPI("get_version", nil)
	if err != nil {
		logger.Warn("HTTP 驱动器连接测试失败", "error", err, "url", baseURL)
		// 不返回错误，因为 HTTP 驱动器可以在后续调用时再连接
	} else {
		logger.Info("HTTP 驱动器连接成功", "url", baseURL)
	}

	d.mu.Lock()
	d.connected = true
	d.mu.Unlock()

	return nil
}

// CallAPI 调用 OneBot API
func (d *HTTPDriver) CallAPI(action string, params map[string]interface{}) (*types.APIResponse, error) {
	// 构建 API URL
	baseURL := d.config.URL
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://%s:%d", d.config.Host, d.config.Port)
	}
	apiURL := fmt.Sprintf("%s/%s", baseURL, action)

	// 构建请求
	var reqBody []byte
	var err error
	if params != nil {
		reqBody, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("序列化请求参数失败: %w", err)
		}
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建 HTTP 请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if d.config.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+d.config.AccessToken)
	}

	// 发送请求
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var apiResp types.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &apiResp, nil
}

// SetEventHandler 设置事件处理器
// 注意：HTTP 驱动器不会接收事件推送
func (d *HTTPDriver) SetEventHandler(handler EventHandler) {
	d.eventHandler = handler
	logger.Warn("HTTP 驱动器不支持接收事件推送，事件处理器将不会被调用")
}

// Close 关闭连接
func (d *HTTPDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.connected = false
	return nil
}

// IsConnected 是否已连接
func (d *HTTPDriver) IsConnected() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.connected
}
