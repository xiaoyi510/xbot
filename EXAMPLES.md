# XBot 实用示例

本文档提供了 XBot 框架的实用示例代码，帮助你快速上手开发。

## 📋 目录

- [基础示例](#基础示例)
- [消息处理](#消息处理)
- [群管功能](#群管功能)
- [数据存储](#数据存储)
- [会话管理](#会话管理)
- [定时任务](#定时任务)
- [API 调用](#api-调用)
- [高级功能](#高级功能)

## 基础示例

### Hello World

```go
package hello

import "xbot"

func init() {
    engine := xbot.NewEngine()
    
    engine.OnCommand("hello").Handle(func(ctx *xbot.Context) {
        ctx.Reply("Hello, World!")
    })
}
```

### 带参数的命令

```go
engine.OnCommand("say").Handle(func(ctx *xbot.Context) {
    args := ctx.GetArgs()
    if args == "" {
        ctx.Reply("请输入要说的话")
        return
    }
    ctx.Reply(args)
})
```

### 关键词触发

```go
engine.OnKeywords([]string{"你好", "hello", "hi"}).Handle(func(ctx *xbot.Context) {
    ctx.Reply("你好！很高兴见到你～")
})
```

### 正则表达式

```go
// 匹配 "抽卡 X 次"
engine.OnRegex(`^抽卡\s*(\d+)次$`).Handle(func(ctx *xbot.Context) {
    if ctx.RegexResult != nil && len(ctx.RegexResult.Groups) > 1 {
        times := ctx.RegexResult.Groups[1]
        ctx.Reply(fmt.Sprintf("正在为你抽卡 %s 次...", times))
    }
})

// 使用命名分组
engine.OnRegex(`^(?P<action>查询|搜索)\s+(?P<keyword>.+)$`).Handle(func(ctx *xbot.Context) {
    if ctx.RegexResult != nil {
        action := ctx.RegexResult.NamedGroups["action"]
        keyword := ctx.RegexResult.NamedGroups["keyword"]
        ctx.Reply(fmt.Sprintf("正在%s：%s", action, keyword))
    }
})
```

## 消息处理

### 发送不同类型的消息

```go
engine.OnCommand("msg").Handle(func(ctx *xbot.Context) {
    // 纯文本
    ctx.Reply("这是纯文本消息")
    
    // 图片
    ctx.Reply("[CQ:image,file=https://example.com/image.jpg]")
    
    // @某人
    userID := ctx.GetUserID()
    ctx.Reply(fmt.Sprintf("[CQ:at,qq=%d] 你好", userID))
    
    // 组合消息
    msg := fmt.Sprintf(
        "[CQ:at,qq=%d] 欢迎！[CQ:image,file=https://example.com/welcome.jpg]",
        userID,
    )
    ctx.Reply(msg)
})
```

### 使用消息构建器

```go
import "github.com/xiaoyi510/xbot/message"

engine.OnCommand("welcome").Handle(func(ctx *xbot.Context) {
    msg := message.NewBuilder().
        Text("欢迎 ").
        At(ctx.GetUserID()).
        Text(" 来到本群！\n").
        Image("https://example.com/welcome.jpg").
        Build()
    
    ctx.Reply(msg)
})
```

### 转发消息

```go
engine.OnCommand("forward").Handle(func(ctx *xbot.Context) {
    // 转发到指定群
    targetGroupID := int64(123456789)
    message := ctx.GetPlainText()
    
    ctx.SendGroupMessage(targetGroupID, fmt.Sprintf(
        "来自 %d 的消息：%s",
        ctx.GetUserID(),
        message,
    ))
    
    ctx.Reply("消息已转发")
})
```

### 消息撤回

```go
engine.OnCommand("recall").Handle(func(ctx *xbot.Context) {
    ctx.Reply("这条消息将在 5 秒后撤回")
    
    time.Sleep(5 * time.Second)
    ctx.Delete()
})
```

## 群管功能

### 禁言功能

```go
engine.OnCommand("ban", 
    xbot.IsGroup(),
    xbot.IsGroupAdmin(),
).Handle(func(ctx *xbot.Context) {
    atUsers := ctx.GetAtUsers()
    if len(atUsers) == 0 {
        ctx.Reply("请 @ 要禁言的成员")
        return
    }
    
    // 获取禁言时长（默认 10 分钟）
    duration := 600
    args := ctx.GetArgs()
    if args != "" {
        if d, err := strconv.Atoi(args); err == nil {
            duration = d * 60 // 转换为秒
        }
    }
    
    groupID := ctx.GetGroupID()
    for _, userID := range atUsers {
        ctx.SetGroupBan(groupID, userID, int64(duration))
    }
    
    ctx.Reply(fmt.Sprintf("已禁言 %d 分钟", duration/60))
})

// 解除禁言
engine.OnCommand("unban",
    xbot.IsGroup(),
    xbot.IsGroupAdmin(),
).Handle(func(ctx *xbot.Context) {
    atUsers := ctx.GetAtUsers()
    if len(atUsers) == 0 {
        ctx.Reply("请 @ 要解除禁言的成员")
        return
    }
    
    groupID := ctx.GetGroupID()
    for _, userID := range atUsers {
        ctx.SetGroupBan(groupID, userID, 0)
    }
    
    ctx.Reply("已解除禁言")
})
```

### 踢人功能

```go
engine.OnCommand("kick",
    xbot.IsGroup(),
    xbot.IsGroupAdmin(),
).Handle(func(ctx *xbot.Context) {
    atUsers := ctx.GetAtUsers()
    if len(atUsers) == 0 {
        ctx.Reply("请 @ 要踢出的成员")
        return
    }
    
    groupID := ctx.GetGroupID()
    for _, userID := range atUsers {
        // 不拒绝再次加群
        ctx.SetGroupKick(groupID, userID, false)
    }
    
    ctx.Reply("已踢出指定成员")
})
```

### 设置群名片

```go
engine.OnCommand("card",
    xbot.IsGroup(),
    xbot.IsGroupAdmin(),
).Handle(func(ctx *xbot.Context) {
    atUsers := ctx.GetAtUsers()
    if len(atUsers) == 0 {
        ctx.Reply("请 @ 要修改名片的成员")
        return
    }
    
    card := ctx.GetArgs()
    if card == "" {
        ctx.Reply("请输入新名片")
        return
    }
    
    groupID := ctx.GetGroupID()
    for _, userID := range atUsers {
        ctx.SetGroupCard(groupID, userID, card)
    }
    
    ctx.Reply("名片已修改")
})
```

### 敏感词检测

```go
// 初始化敏感词库
var sensitiveWords = xbot.NewKeywordManager([]string{
    "违禁词1",
    "违禁词2",
    "敏感词",
})

engine.OnACKeywords(sensitiveWords, xbot.IsGroup()).Handle(func(ctx *xbot.Context) {
    // 撤回消息
    ctx.Delete()
    
    // 警告用户
    ctx.Reply(fmt.Sprintf(
        "[CQ:at,qq=%d] 请勿发送违规内容",
        ctx.GetUserID(),
    ))
    
    // 如果不是管理员，则禁言
    if !ctx.IsAdmin() {
        ctx.SetGroupBan(ctx.GetGroupID(), ctx.GetUserID(), 60)
        ctx.Reply("已禁言 1 分钟")
    }
})
```

### 欢迎新成员

```go
engine.OnNotice().Handle(func(ctx *xbot.Context) {
    // 判断是否为群成员增加事件
    if notice, ok := ctx.Event.(*event.GroupIncreaseNoticeEvent); ok {
        welcomeMsg := message.NewBuilder().
            At(notice.UserID).
            Text(" 欢迎加入本群！\n").
            Text("发送 /help 查看帮助信息").
            Build()
        
        ctx.SendGroupMessage(notice.GroupID, welcomeMsg)
    }
})
```

## 数据存储

### 签到系统

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
    LastCheckIn int64 `json:"last_checkin"`
    Days        int   `json:"days"`
    Total       int   `json:"total"`
}

func init() {
    storage = xbot.GetStorage("checkin")
    engine := xbot.NewEngine()
    
    // 签到命令
    engine.OnCommand("checkin").Handle(handleCheckIn)
    
    // 查询签到信息
    engine.OnCommand("mycheck").Handle(handleMyCheckIn)
}

func handleCheckIn(ctx *xbot.Context) {
    userID := ctx.GetUserID()
    key := fmt.Sprintf("user:%d", userID)
    
    // 获取签到数据
    data := &CheckInData{}
    if raw, err := storage.Get([]byte(key)); err == nil {
        json.Unmarshal(raw, data)
    }
    
    // 检查是否已签到
    today := time.Now().Format("2006-01-02")
    lastCheckIn := time.Unix(data.LastCheckIn, 0).Format("2006-01-02")
    
    if today == lastCheckIn {
        ctx.Reply("你今天已经签到过了！")
        return
    }
    
    // 更新签到数据
    now := time.Now().Unix()
    yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
    
    if lastCheckIn == yesterday {
        data.Days++ // 连续签到
    } else {
        data.Days = 1 // 重置连续签到
    }
    
    data.LastCheckIn = now
    data.Total++
    
    // 保存数据
    raw, _ := json.Marshal(data)
    storage.Set([]byte(key), raw)
    
    ctx.Reply(fmt.Sprintf(
        "签到成功！\n连续签到：%d 天\n累计签到：%d 天",
        data.Days,
        data.Total,
    ))
}

func handleMyCheckIn(ctx *xbot.Context) {
    userID := ctx.GetUserID()
    key := fmt.Sprintf("user:%d", userID)
    
    data := &CheckInData{}
    raw, err := storage.Get([]byte(key))
    if err != nil {
        ctx.Reply("你还没有签到过")
        return
    }
    
    json.Unmarshal(raw, data)
    
    ctx.Reply(fmt.Sprintf(
        "签到信息：\n连续签到：%d 天\n累计签到：%d 天",
        data.Days,
        data.Total,
    ))
}
```

### 积分系统

```go
package points

import (
    "encoding/json"
    "fmt"
    "xbot"
)

var storage xbot.Storage

func init() {
    storage = xbot.GetStorage("points")
    engine := xbot.NewEngine()
    
    engine.OnCommand("points").Handle(handlePoints)
    engine.OnCommand("give", xbot.IsSuperUser()).Handle(handleGive)
}

func getPoints(userID int64) int {
    key := fmt.Sprintf("points:%d", userID)
    data, err := storage.Get([]byte(key))
    if err != nil {
        return 0
    }
    
    var points int
    json.Unmarshal(data, &points)
    return points
}

func setPoints(userID int64, points int) error {
    key := fmt.Sprintf("points:%d", userID)
    data, _ := json.Marshal(points)
    return storage.Set([]byte(key), data)
}

func handlePoints(ctx *xbot.Context) {
    userID := ctx.GetUserID()
    points := getPoints(userID)
    ctx.Reply(fmt.Sprintf("你的积分：%d", points))
}

func handleGive(ctx *xbot.Context) {
    atUsers := ctx.GetAtUsers()
    if len(atUsers) == 0 {
        ctx.Reply("请 @ 要赠送积分的用户")
        return
    }
    
    args := ctx.GetArgs()
    points, err := strconv.Atoi(args)
    if err != nil || points <= 0 {
        ctx.Reply("请输入正确的积分数量")
        return
    }
    
    for _, userID := range atUsers {
        current := getPoints(userID)
        setPoints(userID, current+points)
    }
    
    ctx.Reply(fmt.Sprintf("已赠送 %d 积分", points))
}
```

## 会话管理

### 问卷调查

```go
engine.OnCommand("survey").Handle(func(ctx *xbot.Context) {
    ctx.Reply("开始问卷调查\n请输入你的姓名：")
    
    // 等待姓名
    nameCtx := ctx.Session.Wait(ctx.GetUserID(), ctx.GetGroupID(), 30*time.Second)
    if nameCtx == nil {
        ctx.Reply("超时，调查已取消")
        return
    }
    name := nameCtx.GetPlainText()
    
    ctx.Reply("请输入你的年龄：")
    
    // 等待年龄
    ageCtx := ctx.Session.Wait(ctx.GetUserID(), ctx.GetGroupID(), 30*time.Second)
    if ageCtx == nil {
        ctx.Reply("超时，调查已取消")
        return
    }
    age := ageCtx.GetPlainText()
    
    ctx.Reply("请输入你的建议：")
    
    // 等待建议
    suggestionCtx := ctx.Session.Wait(ctx.GetUserID(), ctx.GetGroupID(), 60*time.Second)
    if suggestionCtx == nil {
        ctx.Reply("超时，调查已取消")
        return
    }
    suggestion := suggestionCtx.GetPlainText()
    
    // 保存调查结果
    result := fmt.Sprintf(
        "问卷调查结果：\n姓名：%s\n年龄：%s\n建议：%s",
        name, age, suggestion,
    )
    
    ctx.Reply("感谢参与问卷调查！")
    
    // 将结果发送给管理员
    ctx.SendPrivateMessage(ctx.Bot.Config.SuperUsers[0], result)
})
```

### 猜数字游戏

```go
engine.OnCommand("guess").Handle(func(ctx *xbot.Context) {
    // 生成随机数
    target := rand.Intn(100) + 1
    attempts := 0
    maxAttempts := 5
    
    ctx.Reply(fmt.Sprintf(
        "猜数字游戏开始！\n我想了一个 1-100 的数字\n你有 %d 次机会",
        maxAttempts,
    ))
    
    for attempts < maxAttempts {
        // 等待用户输入
        guessCtx := ctx.Session.Wait(
            ctx.GetUserID(),
            ctx.GetGroupID(),
            30*time.Second,
        )
        
        if guessCtx == nil {
            ctx.Reply("超时，游戏结束")
            return
        }
        
        attempts++
        guess, err := strconv.Atoi(guessCtx.GetPlainText())
        if err != nil {
            ctx.Reply("请输入一个数字")
            continue
        }
        
        if guess == target {
            ctx.Reply(fmt.Sprintf(
                "🎉 恭喜你猜对了！\n答案是 %d\n用了 %d 次机会",
                target, attempts,
            ))
            return
        }
        
        remaining := maxAttempts - attempts
        if guess < target {
            ctx.Reply(fmt.Sprintf(
                "太小了！还有 %d 次机会",
                remaining,
            ))
        } else {
            ctx.Reply(fmt.Sprintf(
                "太大了！还有 %d 次机会",
                remaining,
            ))
        }
    }
    
    ctx.Reply(fmt.Sprintf(
        "游戏结束！答案是 %d",
        target,
    ))
})
```

## 定时任务

### 定时提醒

```go
package reminder

import (
    "time"
    "xbot"
)

func init() {
    engine := xbot.NewEngine()
    
    // 启动定时任务
    go startScheduledTasks(engine)
    
    engine.OnCommand("remind").Handle(handleRemind)
}

func startScheduledTasks(engine *xbot.Engine) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        // 每小时执行一次任务
        checkReminders()
    }
}

func handleRemind(ctx *xbot.Context) {
    // 实现提醒功能
    ctx.Reply("提醒已设置")
}
```

### 每日任务

```go
func init() {
    engine := xbot.NewEngine()
    
    // 每天早上 8 点发送早报
    go func() {
        for {
            now := time.Now()
            next := time.Date(
                now.Year(), now.Month(), now.Day(),
                8, 0, 0, 0,
                now.Location(),
            )
            
            if now.After(next) {
                next = next.Add(24 * time.Hour)
            }
            
            time.Sleep(time.Until(next))
            
            // 发送早报到所有群
            sendMorningNews(engine)
        }
    }()
}

func sendMorningNews(engine *xbot.Engine) {
    // 获取新闻内容
    news := fetchNews()
    
    // 发送到所有群（需要维护群列表）
    groupIDs := []int64{123456, 789012}
    
    for _, groupID := range groupIDs {
        // 这里需要通过 Bot 实例发送
        // 可以通过全局变量或其他方式获取 Bot 实例
    }
}
```

## API 调用

### 获取群成员列表

```go
engine.OnCommand("members", xbot.IsGroup(), xbot.IsGroupAdmin()).Handle(func(ctx *xbot.Context) {
    groupID := ctx.GetGroupID()
    
    members, err := ctx.API.GetGroupMemberList(groupID)
    if err != nil {
        ctx.Reply("获取失败：" + err.Error())
        return
    }
    
    ctx.Reply(fmt.Sprintf("群成员数量：%d", len(members)))
})
```

### 获取群信息

```go
engine.OnCommand("groupinfo", xbot.IsGroup()).Handle(func(ctx *xbot.Context) {
    groupID := ctx.GetGroupID()
    
    info, err := ctx.API.GetGroupInfo(groupID, false)
    if err != nil {
        ctx.Reply("获取失败：" + err.Error())
        return
    }
    
    msg := fmt.Sprintf(
        "群信息：\n群号：%d\n群名：%s\n成员数：%d",
        info.GroupID,
        info.GroupName,
        info.MemberCount,
    )
    
    ctx.Reply(msg)
})
```

### 获取用户信息

```go
engine.OnCommand("userinfo").Handle(func(ctx *xbot.Context) {
    atUsers := ctx.GetAtUsers()
    var userID int64
    
    if len(atUsers) > 0 {
        userID = atUsers[0]
    } else {
        userID = ctx.GetUserID()
    }
    
    if ctx.GetGroupID() != 0 {
        // 群成员信息
        info, err := ctx.API.GetGroupMemberInfo(
            ctx.GetGroupID(),
            userID,
            false,
        )
        if err != nil {
            ctx.Reply("获取失败：" + err.Error())
            return
        }
        
        msg := fmt.Sprintf(
            "用户信息：\nQQ：%d\n昵称：%s\n群名片：%s\n角色：%s",
            info.UserID,
            info.Nickname,
            info.Card,
            info.Role,
        )
        ctx.Reply(msg)
    } else {
        // 好友信息
        ctx.Reply("私聊暂不支持查询用户信息")
    }
})
```

## 高级功能

### 消息限流

```go
// 全局限流：每个用户 10 秒内只能调用 3 次
engine.OnCommand("api").
    Limit(10*time.Second, 3, func(ctx *xbot.Context) {
        ctx.Reply("调用过于频繁，请稍后再试")
    }).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("API 调用成功")
    })
```

### 多级权限系统

```go
// 定义权限等级
const (
    PermissionUser  = 0
    PermissionVIP   = 1
    PermissionAdmin = 2
    PermissionOwner = 3
)

// 自定义权限过滤器
func HasPermission(level int) xbot.Filter {
    return func(ctx *xbot.Context) bool {
        userID := ctx.GetUserID()
        userLevel := getUserPermission(userID)
        return userLevel >= level
    }
}

// 使用权限过滤器
engine.OnCommand("vip", HasPermission(PermissionVIP)).Handle(func(ctx *xbot.Context) {
    ctx.Reply("VIP 专属功能")
})

engine.OnCommand("admin", HasPermission(PermissionAdmin)).Handle(func(ctx *xbot.Context) {
    ctx.Reply("管理员功能")
})
```

### 自定义中间件

```go
// 日志中间件
func LogMiddleware() func(next func(*xbot.Context)) func(*xbot.Context) {
    return func(next func(*xbot.Context)) func(*xbot.Context) {
        return func(ctx *xbot.Context) {
            start := time.Now()
            
            ctx.Logger.Info("开始处理",
                "user", ctx.GetUserID(),
                "group", ctx.GetGroupID(),
            )
            
            next(ctx)
            
            ctx.Logger.Info("处理完成",
                "duration", time.Since(start),
            )
        }
    }
}

// 使用中间件
engine.Use(LogMiddleware())
```

### 错误处理

```go
engine.OnCommand("test").Handle(func(ctx *xbot.Context) {
    defer func() {
        if err := recover(); err != nil {
            ctx.Logger.Error("处理出错", "error", err)
            ctx.Reply("操作失败，请稍后重试")
        }
    }()
    
    // 可能出错的操作
    result := doSomething()
    ctx.Reply(result)
})
```

### 优先级控制

```go
// 高优先级（先执行）
engine.OnCommand("urgent").
    Priority(100).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("紧急命令")
    })

// 普通优先级
engine.OnCommand("normal").
    Priority(50).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("普通命令")
    })

// 低优先级（后执行）
engine.OnCommand("lazy").
    Priority(10).
    Handle(func(ctx *xbot.Context) {
        ctx.Reply("低优先级命令")
    })
```

---

更多示例请参考 [README.md](./README.md) 和项目中的插件代码。

