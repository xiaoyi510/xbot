# XBot - OneBot 机器人框架

XBot 是一个基于 Go 语言开发的高性能 OneBot 机器人框架，提供简洁优雅的 API，支持多种驱动器和丰富的匹配器。

## ✨ 特性

- 🚀 **高性能**: 基于 Go 语言，支持高并发处理
- 🎯 **多种匹配器**: 命令、关键词、正则、前缀、后缀等多种匹配方式
- 🔌 **多驱动器支持**: WebSocket、反向 WebSocket、HTTP 等
- 🛡️ **完善的过滤器**: 群组、私聊、权限等多种过滤条件
- ⚡ **限流机制**: 内置限流器，支持 Redis 分布式限流
- 🎨 **中间件系统**: 灵活的中间件支持，可自定义处理流程
- 💾 **多种存储**: 支持内存、LevelDB、Redis 等存储方式
- 📦 **插件化**: 支持插件化开发，易于扩展
- 🔍 **AC 自动机**: 高性能关键词匹配，适用于敏感词过滤等场景
- 💬 **会话管理**: 内置会话管理器，支持多轮对话

## 📦 安装

```bash
go get -u github.com/xiaoyi510/xarr-onebot-xbot
```

## 🚀 快速开始

### 1. 创建配置文件

创建 `config.yaml` 配置文件：

```yaml
bot:
  nickname: ["机器人", "小助手"]
  super_users: [123456789]
  command_prefix: "/"

drivers:
  - type: ws_reverse
    url: "ws://127.0.0.1:8080"
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

### 2. 创建主程序

创建 `main.go` 文件：

```go
package main

import (
    "xbot"
    "xbot/logger"
)

func main() {
    // 加载配置文件
    cfg, err := xbot.LoadConfigFile("./config/config.yaml")
    if err != nil {
        panic(err)
    }

    // 运行并监听
    err = xbot.RunAndListen(cfg)
    if err != nil {
        panic(err)
    }
}
```

### 3. 创建插件

创建 `plugins/hello/hello.go` 文件：

```go
package hello

import "xbot"

func init() {
    engine := xbot.NewEngine()

    // 简单的命令响应
    engine.OnCommand("hello").Handle(func(ctx *xbot.Context) {
        ctx.Reply("Hello, World!")
    })

    // 带参数的命令
    engine.OnCommand("echo").Handle(func(ctx *xbot.Context) {
        args := ctx.GetArgs()
        if len(args) > 0 {
            ctx.Reply(args)
        } else {
            ctx.Reply("请输入要回复的内容")
        }
    })
}
```

### 4. 在主程序中引入插件

```go
package main

import (
    _ "yourproject/plugins/hello"  // 导入插件
    
    "xbot"
)

func main() {
    cfg, err := xbot.LoadConfigFile("./config/config.yaml")
    if err != nil {
        panic(err)
    }

    err = xbot.RunAndListen(cfg)
    if err != nil {
        panic(err)
    }
}
```

### 5. 运行

```bash
go run main.go
```

## 📖 核心概念

### Engine (引擎)

引擎是 XBot 的核心组件，负责注册匹配器和处理事件。每个插件通常创建一个独立的引擎。

```go
engine := xbot.NewEngine()
```

### Matcher (匹配器)

匹配器定义了何时触发事件处理。XBot 提供多种内置匹配器：

- **OnCommand**: 命令匹配
- **OnKeywords**: 关键词匹配
- **OnRegex**: 正则表达式匹配
- **OnPrefix**: 前缀匹配
- **OnSuffix**: 后缀匹配
- **OnFullMatch**: 完全匹配
- **OnMessage**: 所有消息事件
- **OnNotice**: 通知事件
- **OnRequest**: 请求事件
- **OnDFAKeywords**: DFA 关键词匹配（高性能）
- **OnACKeywords**: AC 自动机关键词匹配（更高性能）

### Context (上下文)

上下文包含了事件的所有信息以及便捷的操作方法。

```go
func handler(ctx *xbot.Context) {
    userID := ctx.GetUserID()
    groupID := ctx.GetGroupID()
    message := ctx.GetPlainText()
    
    ctx.Reply("收到消息：" + message)
}
```

### Filter (过滤器)

过滤器用于进一步筛选事件，只有通过所有过滤器的事件才会触发处理器。

```go
engine.OnCommand("admin", xbot.IsGroupAdmin(), xbot.IsInGroup(123456)).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("管理员命令")
    })
```

## 🎯 匹配器使用

### 命令匹配

```go
// 单个命令
engine.OnCommand("help").Handle(func(ctx *xbot.Context) {
    ctx.Reply("帮助信息")
})

// 命令组
engine.OnCommandGroup([]string{"start", "begin"}).Handle(func(ctx *xbot.Context) {
    ctx.Reply("开始！")
})

// 带参数的命令
engine.OnCommand("say").Handle(func(ctx *xbot.Context) {
    args := ctx.GetArgs()
    ctx.Reply("你说：" + args)
})
```

### 关键词匹配

```go
// 简单关键词
engine.OnKeywords([]string{"你好", "hello"}).Handle(func(ctx *xbot.Context) {
    ctx.Reply("你好！")
})

// 前缀匹配
engine.OnPrefix("查询").Handle(func(ctx *xbot.Context) {
    query := ctx.GetPlainText()[6:] // 去掉"查询"两字
    ctx.Reply("正在查询：" + query)
})

// 后缀匹配
engine.OnSuffix("天气").Handle(func(ctx *xbot.Context) {
    ctx.Reply("天气查询功能")
})

// 完全匹配
engine.OnFullMatch("签到").Handle(func(ctx *xbot.Context) {
    ctx.Reply("签到成功！")
})
```

### 正则表达式匹配

```go
// 基础正则
engine.OnRegex(`^抽卡\s*(\d+)次$`).Handle(func(ctx *xbot.Context) {
    if ctx.RegexResult != nil && len(ctx.RegexResult.Groups) > 1 {
        times := ctx.RegexResult.Groups[1]
        ctx.Reply("抽卡 " + times + " 次")
    }
})

// 命名分组
engine.OnRegex(`^(?P<action>查询|搜索)\s+(?P<keyword>.+)$`).Handle(func(ctx *xbot.Context) {
    if ctx.RegexResult != nil {
        action := ctx.RegexResult.NamedGroups["action"]
        keyword := ctx.RegexResult.NamedGroups["keyword"]
        ctx.Reply(action + "：" + keyword)
    }
})
```

### 高性能关键词匹配

```go
// DFA 关键词匹配（适用于中等规模关键词库）
manager := xbot.NewKeywordManager([]string{"敏感词1", "敏感词2"})
engine.OnDFAKeywords(manager).Handle(func(ctx *xbot.Context) {
    ctx.Reply("检测到敏感词")
    ctx.Delete() // 撤回消息
})

// AC 自动机匹配（适用于大规模关键词库）
acManager := xbot.NewKeywordManager([]string{"违禁词1", "违禁词2"})
engine.OnACKeywords(acManager).Handle(func(ctx *xbot.Context) {
    ctx.Reply("检测到违禁词")
    ctx.Delete()
})

// 动态关键词匹配（根据上下文选择不同的关键词库）
provider := NewContextProvider() // 自定义实现
engine.OnACKeywordsWithContext(provider).Handle(func(ctx *xbot.Context) {
    ctx.Reply("匹配到关键词")
})
```

## 🛡️ 过滤器

### 内置过滤器

```go
// 群组过滤
engine.OnCommand("test", xbot.IsInGroup(123456, 789012)).Handle(handler)

// 私聊过滤
engine.OnCommand("test", xbot.IsPrivate()).Handle(handler)

// 群聊过滤
engine.OnCommand("test", xbot.IsGroup()).Handle(handler)

// 超级用户过滤
engine.OnCommand("admin", xbot.IsSuperUser()).Handle(handler)

// 群管理员过滤
engine.OnCommand("kick", xbot.IsGroupAdmin()).Handle(handler)

// 群主过滤
engine.OnCommand("transfer", xbot.IsGroupOwner()).Handle(handler)

// ToMe 过滤（被 @ 或私聊）
engine.OnCommand("test", xbot.ToMe()).Handle(handler)
```

### 组合过滤器

```go
// 与运算（所有条件都满足）
engine.OnCommand("test", 
    xbot.IsInGroup(123456),
    xbot.IsGroupAdmin(),
).Handle(handler)

// 或运算
engine.OnCommand("test",
    xbot.Or(xbot.IsPrivate(), xbot.IsSuperUser()),
).Handle(handler)

// 非运算
engine.OnCommand("test",
    xbot.Not(xbot.IsInGroup(123456)),
).Handle(handler)
```

### 自定义过滤器

```go
// 自定义过滤器函数
func IsVIP() xbot.Filter {
    return func(ctx *xbot.Context) bool {
        userID := ctx.GetUserID()
        // 检查用户是否为 VIP
        return checkVIP(userID)
    }
}

engine.OnCommand("vip", IsVIP()).Handle(handler)
```

## ⚡ 限流

### 内存限流

```go
// 限制每个用户 10 秒内只能调用 3 次
engine.OnCommand("check").
    Limit(10*time.Second, 3, func(ctx *xbot.Context) {
        ctx.Reply("调用过于频繁，请稍后再试")
    }).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("查询成功")
    })
```

### Redis 限流

```go
// 使用 Redis 实现分布式限流
limiter := xbot.NewRedisLimiter(
    ctx.Bot.Config.Redis,
    "command:check",
    10*time.Second,
    3,
    func(ctx *xbot.Context) {
        ctx.Reply("调用过于频繁")
    },
)

engine.OnCommand("check").
    LimitWithRedis(limiter).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("查询成功")
    })
```

## 🎨 中间件

### 使用内置中间件

```go
engine := xbot.NewEngine()

// 添加日志中间件
engine.UseLogger()

// 添加异常恢复中间件
engine.UseRecovery()

// 添加性能监控中间件
engine.UseMetrics()
```

### 自定义中间件

```go
// 自定义中间件
engine.Use(func(next func(*xbot.Context)) func(*xbot.Context) {
    return func(ctx *xbot.Context) {
        // 前置处理
        start := time.Now()
        
        // 调用下一个中间件或处理器
        next(ctx)
        
        // 后置处理
        duration := time.Since(start)
        ctx.Logger.Info("处理完成", "duration", duration)
    }
})
```

### 匹配器级别中间件

```go
// 只对特定匹配器应用中间件
engine.OnCommand("test").
    Use(func(next func(*xbot.Context)) func(*xbot.Context) {
        return func(ctx *xbot.Context) {
            // 自定义逻辑
            next(ctx)
        }
    }).
    Handle(handler)
```

## 💬 会话管理

会话管理用于实现多轮对话：

```go
engine.OnCommand("survey").Handle(func(ctx *xbot.Context) {
    ctx.Reply("请输入你的姓名：")
    
    // 等待用户输入姓名
    nameCtx := ctx.Session.Wait(ctx.GetUserID(), ctx.GetGroupID(), 30*time.Second)
    if nameCtx == nil {
        ctx.Reply("超时，调查已取消")
        return
    }
    name := nameCtx.GetPlainText()
    
    ctx.Reply("你好 " + name + "！请输入你的年龄：")
    
    // 等待用户输入年龄
    ageCtx := ctx.Session.Wait(ctx.GetUserID(), ctx.GetGroupID(), 30*time.Second)
    if ageCtx == nil {
        ctx.Reply("超时，调查已取消")
        return
    }
    age := ageCtx.GetPlainText()
    
    ctx.Reply("感谢参与！姓名：" + name + "，年龄：" + age)
})
```

## 💾 数据存储

### 使用插件存储

```go
var storage xbot.Storage

func init() {
    // 获取插件专用存储
    storage = xbot.GetStorage("myplugin")
    
    engine := xbot.NewEngine()
    // ... 注册处理器
}

// 存储数据
func saveData(key string, value string) error {
    return storage.Set([]byte(key), []byte(value))
}

// 读取数据
func getData(key string) (string, error) {
    data, err := storage.Get([]byte(key))
    if err != nil {
        return "", err
    }
    return string(data), nil
}
```

### 使用 Bot 存储

```go
engine.OnCommand("save").Handle(func(ctx *xbot.Context) {
    args := ctx.GetArgs()
    parts := strings.SplitN(args, " ", 2)
    
    if len(parts) != 2 {
        ctx.Reply("用法: /save <key> <value>")
        return
    }
    
    key := parts[0]
    value := parts[1]
    
    err := ctx.Storage.Set([]byte(key), []byte(value))
    if err != nil {
        ctx.Reply("保存失败：" + err.Error())
        return
    }
    
    ctx.Reply("保存成功")
})

engine.OnCommand("load").Handle(func(ctx *xbot.Context) {
    key := ctx.GetArgs()
    
    value, err := ctx.Storage.Get([]byte(key))
    if err != nil {
        ctx.Reply("读取失败：" + err.Error())
        return
    }
    
    ctx.Reply("值：" + string(value))
})
```

## 📤 消息发送

### 基础消息

```go
// 快速回复
ctx.Reply("Hello")

// 发送私聊消息
ctx.SendPrivateMessage(123456, "私聊消息")

// 发送群消息
ctx.SendGroupMessage(789012, "群消息")
```

### 消息构建器

```go
// 使用消息构建器
msg := message.NewBuilder().
    Text("欢迎 ").
    At(ctx.GetUserID()).
    Text(" 加入群聊！\n").
    Image("https://example.com/image.jpg").
    Build()

ctx.Reply(msg)
```

### CQ 码消息

```go
// 发送图片
ctx.Reply("[CQ:image,file=https://example.com/image.jpg]")

// 发送语音
ctx.Reply("[CQ:record,file=base64://...]")

// @某人
ctx.Reply(fmt.Sprintf("[CQ:at,qq=%d] 你好", userID))

// 多种元素组合
ctx.Reply(fmt.Sprintf(
    "[CQ:at,qq=%d] 你抽到了 [CQ:image,file=https://example.com/card.png]",
    ctx.GetUserID(),
))
```

## 🎮 Context API

### 获取事件信息

```go
// 获取用户 ID
userID := ctx.GetUserID()

// 获取群号
groupID := ctx.GetGroupID()

// 获取消息 ID
messageID := ctx.GetMessageID()

// 获取原始文本
text := ctx.GetPlainText()

// 获取命令参数
args := ctx.GetArgs()

// 获取消息对象
msg := ctx.GetMessage()
```

### 消息操作

```go
// 回复消息
ctx.Reply("回复内容")

// 撤回消息
ctx.Delete()

// 撤回指定消息
ctx.API.DeleteMsg(messageID)
```

### 群组操作

```go
// 踢出群成员
ctx.SetGroupKick(groupID, userID, false)

// 禁言
ctx.SetGroupBan(groupID, userID, 600) // 禁言 600 秒

// 设置群名片
ctx.SetGroupCard(groupID, userID, "新名片")

// 设置管理员
ctx.SetGroupAdmin(groupID, userID, true)
```

### 信息查询

```go
// 获取登录信息
info, _ := ctx.API.GetLoginInfo()

// 获取群成员信息
memberInfo, _ := ctx.API.GetGroupMemberInfo(groupID, userID, false)

// 获取群信息
groupInfo, _ := ctx.API.GetGroupInfo(groupID, false)

// 获取群成员列表
members, _ := ctx.API.GetGroupMemberList(groupID)
```

## 🔧 高级功能

### 设置优先级

```go
// 优先级越高越先执行
engine.OnCommand("high").
    Priority(100).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("高优先级")
    })

engine.OnCommand("low").
    Priority(1).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("低优先级")
    })
```

### 阻止后续匹配

```go
// 匹配后阻止其他处理器执行
engine.OnCommand("stop").
    SetBlock().
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("已阻止")
    })

// 不阻止
engine.OnCommand("continue").
    SetBlock(false).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("继续匹配")
    })
```

### 中止事件传播

```go
engine.OnCommand("abort").Handle(func(ctx *xbot.Context) {
    ctx.Reply("中止后续所有匹配器")
    ctx.Abort() // 中止后续匹配器
})
```

## 🔌 驱动器配置

### 反向 WebSocket（推荐）

```yaml
drivers:
  - type: ws_reverse
    url: "ws://127.0.0.1:8080"
    access_token: "your_token"
    reconnect_interval: 5
    max_reconnect: 0
    timeout: 30
```

### 正向 WebSocket

```yaml
drivers:
  - type: ws
    url: "ws://127.0.0.1:5700"
    access_token: "your_token"
    reconnect_interval: 5
    heartbeat_interval: 30
```

### HTTP

```yaml
drivers:
  - type: http
    url: "http://127.0.0.1:5700"
    access_token: "your_token"
    timeout: 30
```

### HTTP POST

```yaml
drivers:
  - type: http_post
    host: "0.0.0.0"
    port: 8080
    url: "http://127.0.0.1:5700"
    access_token: "your_token"
```

## 📝 完整示例

### 天气查询插件

```go
package weather

import (
    "fmt"
    "strings"
    "time"
    
    "xbot"
)

var storage xbot.Storage

func init() {
    storage = xbot.GetStorage("weather")
    engine := xbot.NewEngine()
    
    // 注册中间件
    engine.UseRecovery()
    engine.UseLogger()
    
    // 天气查询命令
    engine.OnCommand("weather").
        Limit(10*time.Second, 1, func(ctx *xbot.Context) {
            ctx.Reply("查询过于频繁，请稍后再试")
        }).
        Handle(handleWeather)
    
    // 设置默认城市
    engine.OnCommand("setcity").
        Handle(handleSetCity)
}

func handleWeather(ctx *xbot.Context) {
    city := ctx.GetArgs()
    
    // 如果没有提供城市，使用默认城市
    if city == "" {
        key := fmt.Sprintf("city:%d", ctx.GetUserID())
        data, err := storage.Get([]byte(key))
        if err != nil {
            ctx.Reply("请提供城市名称：/weather 北京")
            return
        }
        city = string(data)
    }
    
    // 查询天气（这里是模拟）
    weather := queryWeather(city)
    ctx.Reply(weather)
}

func handleSetCity(ctx *xbot.Context) {
    city := ctx.GetArgs()
    if city == "" {
        ctx.Reply("请提供城市名称：/setcity 北京")
        return
    }
    
    key := fmt.Sprintf("city:%d", ctx.GetUserID())
    err := storage.Set([]byte(key), []byte(city))
    if err != nil {
        ctx.Reply("设置失败：" + err.Error())
        return
    }
    
    ctx.Reply("默认城市已设置为：" + city)
}

func queryWeather(city string) string {
    // 这里应该调用实际的天气 API
    return fmt.Sprintf("%s的天气：晴，温度 25℃", city)
}
```

### 群管插件

```go
package admin

import (
    "time"
    
    "xbot"
)

func init() {
    engine := xbot.NewEngine()
    
    // 踢人命令（仅管理员可用）
    engine.OnCommand("kick",
        xbot.IsGroup(),
        xbot.IsGroupAdmin(),
    ).Handle(func(ctx *xbot.Context) {
        // 获取被 @ 的用户
        atUsers := ctx.GetAtUsers()
        if len(atUsers) == 0 {
            ctx.Reply("请 @ 要踢出的成员")
            return
        }
        
        for _, userID := range atUsers {
            ctx.SetGroupKick(ctx.GetGroupID(), userID, false)
        }
        
        ctx.Reply("已踢出指定成员")
    })
    
    // 禁言命令
    engine.OnCommand("ban",
        xbot.IsGroup(),
        xbot.IsGroupAdmin(),
    ).Handle(func(ctx *xbot.Context) {
        atUsers := ctx.GetAtUsers()
        if len(atUsers) == 0 {
            ctx.Reply("请 @ 要禁言的成员")
            return
        }
        
        // 默认禁言 10 分钟
        duration := 600
        
        for _, userID := range atUsers {
            ctx.SetGroupBan(ctx.GetGroupID(), userID, int64(duration))
        }
        
        ctx.Reply("已禁言指定成员 10 分钟")
    })
    
    // 敏感词检测
    keywords := xbot.NewKeywordManager([]string{
        "违禁词1", "违禁词2", "敏感词",
    })
    
    engine.OnACKeywords(keywords, xbot.IsGroup()).
        Handle(func(ctx *xbot.Context) {
            // 撤回消息
            ctx.Delete()
            
            // 如果不是管理员，则禁言
            if !ctx.IsAdmin() {
                ctx.SetGroupBan(
                    ctx.GetGroupID(),
                    ctx.GetUserID(),
                    60, // 禁言 1 分钟
                )
                ctx.Reply("检测到违规内容，已禁言 1 分钟")
            }
        })
}
```

## 📚 API 参考

### Engine 方法

| 方法 | 说明 |
|------|------|
| `NewEngine()` | 创建新引擎 |
| `Use(middlewares...)` | 添加全局中间件 |
| `OnCommand(cmd, filters...)` | 命令匹配 |
| `OnKeywords(keywords, filters...)` | 关键词匹配 |
| `OnRegex(pattern, filters...)` | 正则匹配 |
| `OnPrefix(prefix, filters...)` | 前缀匹配 |
| `OnSuffix(suffix, filters...)` | 后缀匹配 |
| `OnFullMatch(text, filters...)` | 完全匹配 |
| `OnMessage(filters...)` | 消息事件 |
| `OnNotice(filters...)` | 通知事件 |
| `OnRequest(filters...)` | 请求事件 |
| `OnDFAKeywords(provider, filters...)` | DFA 关键词匹配 |
| `OnACKeywords(provider, filters...)` | AC 自动机匹配 |

### Matcher 方法

| 方法 | 说明 |
|------|------|
| `Handle(handler)` | 设置处理函数 |
| `Filter(filters...)` | 添加过滤器 |
| `Limit(duration, count, onExceed)` | 设置限流 |
| `Priority(p)` | 设置优先级 |
| `Use(middlewares...)` | 添加中间件 |
| `SetBlock(block...)` | 阻止后续匹配 |

### Context 方法

| 方法 | 说明 |
|------|------|
| `GetUserID()` | 获取用户 ID |
| `GetGroupID()` | 获取群号 |
| `GetMessageID()` | 获取消息 ID |
| `GetPlainText()` | 获取纯文本 |
| `GetArgs()` | 获取命令参数 |
| `GetMessage()` | 获取消息对象 |
| `GetAtUsers()` | 获取被 @ 的用户列表 |
| `Reply(msg)` | 快速回复 |
| `Delete()` | 撤回消息 |
| `SendPrivateMessage(userID, msg)` | 发送私聊消息 |
| `SendGroupMessage(groupID, msg)` | 发送群消息 |
| `SetGroupKick(groupID, userID, reject)` | 踢出群成员 |
| `SetGroupBan(groupID, userID, duration)` | 禁言 |
| `SetGroupCard(groupID, userID, card)` | 设置群名片 |
| `IsAdmin()` | 是否为管理员 |
| `IsOwner()` | 是否为群主 |
| `Abort()` | 中止后续匹配 |

### 内置过滤器

| 过滤器 | 说明 |
|--------|------|
| `IsPrivate()` | 私聊消息 |
| `IsGroup()` | 群聊消息 |
| `IsInGroup(groupIDs...)` | 指定群组 |
| `IsSuperUser()` | 超级用户 |
| `IsGroupAdmin()` | 群管理员 |
| `IsGroupOwner()` | 群主 |
| `ToMe()` | 被 @ 或私聊 |
| `Or(filters...)` | 或运算 |
| `And(filters...)` | 与运算 |
| `Not(filter)` | 非运算 |

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🔗 相关链接

- [OneBot 标准](https://github.com/botuniverse/onebot)
- [Go-CQHTTP](https://github.com/Mrs4s/go-cqhttp)
- [Lagrange](https://github.com/LagrangeDev/Lagrange.Core)

---

Made with ❤️ by XArr Team

