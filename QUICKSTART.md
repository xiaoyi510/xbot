# XBot 快速入门指南

本指南将带你在 5 分钟内快速上手 XBot 框架。

## 🎯 前置要求

- Go 1.18 或更高版本
- 一个 OneBot 实现（如 go-cqhttp、Lagrange 等）

## 📦 第一步：创建项目

```bash
# 创建项目目录
mkdir mybot
cd mybot

# 初始化 Go 模块
go mod init mybot

# 安装 XBot
go get -u github.com/xiaoyi510/xbot
```

## ⚙️ 第二步：配置文件

创建 `config/config.yaml`：

```bash
mkdir -p config
```

编辑 `config/config.yaml`：

```yaml
bot:
  nickname: ["机器人", "小助手"]
  super_users: [123456789]  # 替换为你的 QQ 号
  command_prefix: "/"

drivers:
  - type: ws_reverse
    url: "ws://127.0.0.1:8080"  # 替换为你的 OneBot 地址
    access_token: ""
    reconnect_interval: 5
    max_reconnect: 0
    timeout: 30

log:
  level: "info"
  file: "logs/bot.log"

storage:
  type: "leveldb"
```

## 🔌 第三步：创建插件

创建 `plugins/hello/hello.go`：

```bash
mkdir -p plugins/hello
```

编辑 `plugins/hello/hello.go`：

```go
package hello

import "github.com/xiaoyi510/xbot"

func init() {
    // 创建引擎
    engine := xbot.NewEngine()

    // Hello 命令
    engine.OnCommand("hello").Handle(func(ctx *xbot.Context) {
        ctx.Reply("Hello, World!")
    })

    // Echo 命令
    engine.OnCommand("echo").Handle(func(ctx *xbot.Context) {
        args := ctx.GetArgs()
        if args == "" {
            ctx.Reply("请输入要回复的内容")
            return
        }
        ctx.Reply(args)
    })

    // 关键词触发
    engine.OnKeywords([]string{"你好", "hello"}).Handle(func(ctx *xbot.Context) {
        ctx.Reply("你好！很高兴见到你～")
    })
}
```

## 🚀 第四步：创建主程序

创建 `main.go`：

```go
package main

import (
    _ "mybot/plugins/hello"  // 导入插件
    
    "xbot"
)

func main() {
    // 加载配置
    cfg, err := xbot.LoadConfigFile("./config/config.yaml")
    if err != nil {
        panic(err)
    }

    // 运行机器人
    err = xbot.RunAndListen(cfg)
    if err != nil {
        panic(err)
    }
}
```

## ▶️ 第五步：运行

```bash
# 安装依赖
go mod tidy

# 运行机器人
go run main.go
```

## 🎉 测试

在 QQ 中发送以下消息：

```
/hello
```

机器人会回复：`Hello, World!`

```
/echo 测试消息
```

机器人会回复：`测试消息`

```
你好
```

机器人会回复：`你好！很高兴见到你～`

## 📚 下一步

恭喜！你已经创建了第一个 XBot 机器人。现在你可以：

1. 查看 [README.md](./README.md) 了解更多功能
2. 查看 [EXAMPLES.md](./EXAMPLES.md) 学习实用示例
3. 查看 [API.md](./API.md) 了解完整的 API 文档

## 🔧 常见问题

### 1. 连接失败

**问题：** 机器人无法连接到 OneBot

**解决：**
- 检查 OneBot 实现是否正常运行
- 确认配置文件中的 URL 和端口正确
- 检查 access_token 是否匹配

### 2. 命令无响应

**问题：** 发送命令后机器人没有响应

**解决：**
- 确认命令前缀是否正确（默认为 `/`）
- 检查是否在正确的群组或私聊中测试
- 查看日志文件 `logs/bot.log` 确认是否收到消息

### 3. 导入错误

**问题：** `cannot find package "xbot"`

**解决：**
```bash
go mod tidy
go mod download
```

## 💡 快速示例

### 群管功能

```go
package admin

import "xbot"

func init() {
    engine := xbot.NewEngine()

    // 踢人（仅管理员可用）
    engine.OnCommand("kick", 
        xbot.IsGroup(),
        xbot.IsGroupAdmin(),
    ).Handle(func(ctx *xbot.Context) {
        atUsers := ctx.GetAtUsers()
        if len(atUsers) == 0 {
            ctx.Reply("请 @ 要踢出的成员")
            return
        }
        
        for _, userID := range atUsers {
            ctx.SetGroupKick(ctx.GetGroupID(), userID, false)
        }
        ctx.Reply("已踢出")
    })

    // 禁言
    engine.OnCommand("ban",
        xbot.IsGroup(),
        xbot.IsGroupAdmin(),
    ).Handle(func(ctx *xbot.Context) {
        atUsers := ctx.GetAtUsers()
        if len(atUsers) == 0 {
            ctx.Reply("请 @ 要禁言的成员")
            return
        }
        
        for _, userID := range atUsers {
            ctx.SetGroupBan(ctx.GetGroupID(), userID, 600) // 禁言10分钟
        }
        ctx.Reply("已禁言 10 分钟")
    })
}
```

### 数据存储

```go
package checkin

import (
    "encoding/json"
    "fmt"
    "time"
    "xbot"
)

var storage xbot.Storage

type CheckInData struct {
    Days  int   `json:"days"`
    Total int   `json:"total"`
}

func init() {
    storage = xbot.GetStorage("checkin")
    engine := xbot.NewEngine()

    engine.OnCommand("checkin").Handle(func(ctx *xbot.Context) {
        userID := ctx.GetUserID()
        key := fmt.Sprintf("user:%d", userID)
        
        data := &CheckInData{}
        if raw, err := storage.Get([]byte(key)); err == nil {
            json.Unmarshal(raw, data)
        }
        
        data.Days++
        data.Total++
        
        raw, _ := json.Marshal(data)
        storage.Set([]byte(key), raw)
        
        ctx.Reply(fmt.Sprintf("签到成功！连续 %d 天", data.Days))
    })
}
```

### 限流功能

```go
engine.OnCommand("query").
    Limit(10*time.Second, 3, func(ctx *xbot.Context) {
        ctx.Reply("调用过于频繁，请 10 秒后再试")
    }).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("查询成功")
    })
```

## 🎯 完整项目结构

```
mybot/
├── main.go                 # 主程序
├── config/
│   └── config.yaml        # 配置文件
├── plugins/               # 插件目录
│   ├── hello/
│   │   └── hello.go
│   ├── admin/
│   │   └── admin.go
│   └── checkin/
│       └── checkin.go
├── logs/                  # 日志目录
│   └── bot.log
├── data/                  # 数据目录
│   └── checkin/          # 插件数据
└── go.mod                # Go 模块文件
```

## 🔗 相关链接

- **OneBot 实现：**
  - [go-cqhttp](https://github.com/Mrs4s/go-cqhttp)
  - [Lagrange](https://github.com/LagrangeDev/Lagrange.Core)
  
- **文档：**
  - [XBot README](./README.md)
  - [示例代码](./EXAMPLES.md)
  - [API 文档](./API.md)

---

开始你的 XBot 开发之旅吧！🚀

