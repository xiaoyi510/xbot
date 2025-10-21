# Driver 驱动器

XBot 支持多种 OneBot 11 通信驱动器，可以根据实际需求选择合适的驱动器类型。

## 驱动器类型

### 1. 反向 WebSocket (ws_reverse) ✅ 推荐

**适用场景**: 最常用的连接方式，Bot 作为 WebSocket 客户端连接到 OneBot 实现。

**特点**:
- ✅ 支持接收实时事件
- ✅ 支持调用 API
- ✅ 自动重连
- ✅ 双向通信

**配置示例**:
```yaml
drivers:
  - type: ws_reverse
    url: "ws://127.0.0.1:8080"
    access_token: "your_token"
    reconnect_interval: 5      # 重连间隔（秒）
    max_reconnect: 0           # 最大重连次数，0 表示无限
    timeout: 30                # API 调用超时（秒）
```

---

### 2. 正向 WebSocket (ws / websocket)

**适用场景**: OneBot 实现提供 WebSocket 服务器，Bot 主动连接。

**特点**:
- ✅ 支持接收实时事件
- ✅ 支持调用 API
- ✅ 自动重连
- ✅ 支持心跳检测
- ✅ 双向通信

**配置示例**:
```yaml
drivers:
  - type: ws
    url: "ws://127.0.0.1:5700"
    # 或者分开配置
    # host: "127.0.0.1"
    # port: 5700
    access_token: "your_token"
    reconnect_interval: 5      # 重连间隔（秒）
    max_reconnect: 0           # 最大重连次数，0 表示无限
    heartbeat_interval: 30     # 心跳间隔（秒），0 表示不发送心跳
    timeout: 30                # API 调用超时（秒）
```

**代码示例**:
```go
driver := driver.NewWebSocketDriver(driver.Config{
    URL:               "ws://127.0.0.1:5700",
    AccessToken:       "your_token",
    ReconnectInterval: 5,
    HeartbeatInterval: 30,
})
```

---

### 3. HTTP

**适用场景**: 只需要调用 API，不需要接收事件推送。

**特点**:
- ✅ 支持调用 API
- ❌ 不支持接收事件
- ✅ 实现简单
- ✅ 适合单向通信场景

**配置示例**:
```yaml
drivers:
  - type: http
    url: "http://127.0.0.1:5700"
    # 或者分开配置
    # host: "127.0.0.1"
    # port: 5700
    access_token: "your_token"
    timeout: 30                # API 调用超时（秒）
```

**代码示例**:
```go
driver := driver.NewHTTPDriver(driver.Config{
    URL:         "http://127.0.0.1:5700",
    AccessToken: "your_token",
    Timeout:     30,
})
```

**注意事项**:
- ⚠️ HTTP 驱动器**不会接收事件**，事件处理器不会被调用
- ⚠️ 适合只需要主动调用 API 的场景（如定时任务、Webhook 响应等）

---

### 4. 反向 HTTP POST (http_post)

**适用场景**: OneBot 实现通过 HTTP POST 推送事件，Bot 提供 HTTP 服务器接收。

**特点**:
- ✅ 支持接收实时事件
- ✅ 支持调用 API
- ✅ 适合穿透防火墙场景
- ✅ 双向通信

**配置示例**:
```yaml
drivers:
  - type: http_post
    # Bot 监听的地址和端口（接收事件）
    host: "0.0.0.0"
    port: 8080
    # OneBot API 地址（调用 API）
    url: "http://127.0.0.1:5700"
    access_token: "your_token"
    timeout: 30                # API 调用超时（秒）
```

**代码示例**:
```go
driver := driver.NewHTTPPostDriver(driver.Config{
    Host:        "0.0.0.0",
    Port:        8080,
    URL:         "http://127.0.0.1:5700",
    AccessToken: "your_token",
})
```

**工作流程**:
1. Bot 启动 HTTP 服务器监听 `host:port`
2. OneBot 实现通过 POST 请求推送事件到 Bot
3. Bot 通过 HTTP POST 调用 OneBot API

---

## 驱动器接口

所有驱动器都实现了 `Driver` 接口：

```go
type Driver interface {
    // 连接到 OneBot 实现
    Connect() error
    
    // 调用 OneBot API
    CallAPI(action string, params map[string]interface{}) (*types.APIResponse, error)
    
    // 设置事件处理器
    SetEventHandler(handler EventHandler)
    
    // 关闭连接
    Close() error
    
    // 是否已连接
    IsConnected() bool
}
```

---

## 配置参数说明

### 通用参数

| 参数 | 类型 | 说明 | 必填 |
|------|------|------|------|
| `type` | string | 驱动器类型 | ✅ |
| `access_token` | string | 访问令牌，用于鉴权 | ❌ |
| `timeout` | int | API 调用超时时间（秒），默认 30 | ❌ |

### WebSocket 相关参数

| 参数 | 类型 | 说明 | 适用驱动器 |
|------|------|------|-----------|
| `url` | string | WebSocket 连接地址 | ws_reverse, ws |
| `reconnect_interval` | int | 重连间隔（秒），默认 5 | ws_reverse, ws |
| `max_reconnect` | int | 最大重连次数，0 表示无限，默认 0 | ws_reverse, ws |
| `heartbeat_interval` | int | 心跳间隔（秒），0 表示不发送，默认 0 | ws |

### HTTP 相关参数

| 参数 | 类型 | 说明 | 适用驱动器 |
|------|------|------|-----------|
| `url` | string | HTTP API 地址 | http, http_post |
| `host` | string | 主机地址 | ws, http, http_post |
| `port` | int | 端口号 | ws, http, http_post |

---

## 如何选择驱动器

### 场景 1: 常规 Bot 开发 → `ws_reverse`
最推荐的方式，适合大多数场景。

### 场景 2: OneBot 实现提供 WebSocket 服务 → `ws`
如 go-cqhttp 的正向 WebSocket 模式。

### 场景 3: 只需调用 API，不需要事件 → `http`
例如：定时任务 Bot、Webhook 响应器等。

### 场景 4: 需要穿透防火墙接收事件 → `http_post`
OneBot 实现主动推送事件到公网服务器。

---

## 多驱动器支持

XBot 支持同时使用多个驱动器：

```yaml
drivers:
  # 主要连接：反向 WebSocket
  - type: ws_reverse
    url: "ws://127.0.0.1:8080"
    
  # 备用连接：HTTP（仅 API 调用）
  - type: http
    url: "http://127.0.0.1:5700"
```

**注意事项**:
- 多驱动器会接收重复的事件（如果多个都支持事件接收）
- API 调用会使用第一个可用的驱动器
- 建议只配置一个主要驱动器

---

## 最佳实践

### 1. 生产环境推荐配置

```yaml
drivers:
  - type: ws_reverse
    url: "ws://127.0.0.1:8080"
    access_token: "${ACCESS_TOKEN}"  # 使用环境变量
    reconnect_interval: 5
    max_reconnect: 0                  # 无限重连
    timeout: 30
```

### 2. 启用 Access Token

所有驱动器都支持 Access Token 鉴权：

```yaml
drivers:
  - type: ws_reverse
    url: "ws://127.0.0.1:8080"
    access_token: "your_secure_token_here"
```

驱动器会自动在请求头中添加：
```
Authorization: Bearer your_secure_token_here
```

### 3. 合理设置超时时间

- **快速响应场景**: `timeout: 10`
- **常规场景**: `timeout: 30`（默认）
- **发送大文件等耗时操作**: `timeout: 60` 或更长

### 4. 配置重连策略

```yaml
reconnect_interval: 5      # 每 5 秒重试一次
max_reconnect: 10          # 最多重试 10 次（之后放弃）
# 或
max_reconnect: 0           # 无限重试（推荐）
```

---

## 故障排查

### WebSocket 连接失败

**问题**: `WebSocket 连接失败: dial tcp ...`

**解决方案**:
1. 检查 OneBot 实现是否已启动
2. 检查 URL 是否正确（`ws://` 而不是 `http://`）
3. 检查防火墙设置
4. 检查 Access Token 是否正确

### HTTP 请求失败

**问题**: `HTTP 请求失败，状态码: 401`

**解决方案**:
1. 检查 Access Token 是否配置正确
2. 确认 OneBot 实现的鉴权配置

### API 调用超时

**问题**: `API 调用超时`

**解决方案**:
1. 增加 `timeout` 配置
2. 检查网络连接
3. 检查 OneBot 实现是否正常响应

### 反向 HTTP 端口占用

**问题**: `HTTP 服务器启动失败: bind: address already in use`

**解决方案**:
1. 更换端口号
2. 检查是否有其他程序占用该端口
3. 使用 `netstat -an | grep 8080` 查看端口占用情况

---

## 开发指南

### 自定义驱动器

如果需要实现自定义驱动器，只需实现 `Driver` 接口：

```go
type MyCustomDriver struct {
    config       Config
    eventHandler EventHandler
}

func (d *MyCustomDriver) Connect() error {
    // 实现连接逻辑
    return nil
}

func (d *MyCustomDriver) CallAPI(action string, params map[string]interface{}) (*types.APIResponse, error) {
    // 实现 API 调用逻辑
    return nil, nil
}

func (d *MyCustomDriver) SetEventHandler(handler EventHandler) {
    d.eventHandler = handler
}

func (d *MyCustomDriver) Close() error {
    return nil
}

func (d *MyCustomDriver) IsConnected() bool {
    return true
}
```

然后在 `bot.go` 中添加对应的创建逻辑。

---

## 参考资料

- [OneBot 11 标准](https://github.com/botuniverse/onebot-11)
- [go-cqhttp 文档](https://docs.go-cqhttp.org/)
- [LLOneBot 文档](https://llonebot.github.io/zh-CN/)

---

## 更新日志

### v1.0.0 (2025-10-21)
- ✅ 新增正向 WebSocket 驱动器 (`ws`)
- ✅ 新增 HTTP 驱动器 (`http`)
- ✅ 新增反向 HTTP POST 驱动器 (`http_post`)
- ✅ 改进驱动器配置加载逻辑
- ✅ 完善文档和示例

