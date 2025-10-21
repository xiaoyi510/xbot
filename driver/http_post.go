package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"xbot/event"
	"xbot/logger"
	"xbot/types"
	"xbot/utils"
)

// HTTPPostDriver 反向 HTTP POST 驱动器
// OneBot 实现通过 HTTP POST 推送事件到我们的服务器
// 我们通过 HTTP POST 调用 OneBot API
type HTTPPostDriver struct {
	config       Config
	client       *http.Client
	server       *http.Server
	mu           sync.RWMutex
	eventHandler EventHandler
	apiResponses *utils.SafeMap[string, chan *types.APIResponse]
	connected    bool
	stopChan     chan struct{}
}

// NewHTTPPostDriver 创建反向 HTTP POST 驱动器
func NewHTTPPostDriver(config Config) *HTTPPostDriver {
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	if config.Port == 0 {
		config.Port = 8080
	}

	return &HTTPPostDriver{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		apiResponses: utils.NewSafeMap[string, chan *types.APIResponse](),
		stopChan:     make(chan struct{}),
	}
}

// Connect 启动 HTTP 服务器接收事件
func (d *HTTPPostDriver) Connect() error {
	// 创建 HTTP 处理器
	mux := http.NewServeMux()
	mux.HandleFunc("/", d.handleEvent)

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", d.config.Host, d.config.Port)
	if d.config.Host == "" {
		addr = fmt.Sprintf(":%d", d.config.Port)
	}

	d.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	d.mu.Lock()
	d.connected = true
	d.mu.Unlock()

	// 启动 HTTP 服务器
	go func() {
		logger.Info("反向 HTTP 服务器启动", "addr", addr)
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP 服务器启动失败", "error", err)
			d.mu.Lock()
			d.connected = false
			d.mu.Unlock()
		}
	}()

	return nil
}

// handleEvent 处理接收到的事件
func (d *HTTPPostDriver) handleEvent(w http.ResponseWriter, r *http.Request) {
	// 验证 Access Token
	if d.config.AccessToken != "" {
		token := r.Header.Get("Authorization")
		expectedToken := "Bearer " + d.config.AccessToken
		if token != expectedToken {
			logger.Warn("无效的 Access Token", "token", token)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// 只接受 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("读取请求体失败", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 判断是事件还是 API 响应
	var base struct {
		Echo string `json:"echo"`
	}

	if err := json.Unmarshal(body, &base); err != nil {
		logger.Error("解析消息失败", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 如果有 echo，说明是 API 响应
	if base.Echo != "" {
		var resp types.APIResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			logger.Error("解析 API 响应失败", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
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

		// 返回空响应
		w.WriteHeader(http.StatusOK)
		return
	}

	// 否则是事件
	evt, err := event.ParseEvent(body)
	if err != nil {
		logger.Error("解析事件失败", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 调用事件处理器
	if d.eventHandler != nil {
		go d.eventHandler(evt)
	}

	// 返回空响应
	w.WriteHeader(http.StatusOK)
}

// CallAPI 调用 OneBot API
func (d *HTTPPostDriver) CallAPI(action string, params map[string]interface{}) (*types.APIResponse, error) {
	// 构建 API URL
	apiURL := d.config.URL
	if apiURL == "" {
		return nil, errors.New("未配置 OneBot API URL")
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

	reqBody, err := json.Marshal(request)
	if err != nil {
		d.apiResponses.Delete(echo)
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 发送 HTTP 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(reqBody))
	if err != nil {
		d.apiResponses.Delete(echo)
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if d.config.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+d.config.AccessToken)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		d.apiResponses.Delete(echo)
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		d.apiResponses.Delete(echo)
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		d.apiResponses.Delete(echo)
		return nil, fmt.Errorf("HTTP 请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 尝试立即解析响应（同步响应）
	var apiResp types.APIResponse
	if err := json.Unmarshal(body, &apiResp); err == nil {
		d.apiResponses.Delete(echo)
		return &apiResp, nil
	}

	// 如果不是同步响应，等待异步响应
	timeout := time.Duration(d.config.Timeout) * time.Second
	select {
	case apiResp := <-respChan:
		return apiResp, nil
	case <-time.After(timeout):
		d.apiResponses.Delete(echo)
		return nil, errors.New("API 调用超时")
	}
}

// SetEventHandler 设置事件处理器
func (d *HTTPPostDriver) SetEventHandler(handler EventHandler) {
	d.eventHandler = handler
}

// Close 关闭 HTTP 服务器
func (d *HTTPPostDriver) Close() error {
	close(d.stopChan)

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.server != nil {
		// 优雅关闭 HTTP 服务器
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := d.server.Shutdown(ctx); err != nil {
			logger.Error("关闭 HTTP 服务器失败", "error", err)
			return err
		}
		d.server = nil
	}

	d.connected = false
	return nil
}

// IsConnected 是否已连接
func (d *HTTPPostDriver) IsConnected() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.connected
}
