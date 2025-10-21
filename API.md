# XBot API 参考文档

本文档详细介绍 XBot 框架的所有 API 接口。

## 📋 目录

- [核心 API](#核心-api)
- [引擎 API](#引擎-api)
- [匹配器 API](#匹配器-api)
- [上下文 API](#上下文-api)
- [消息 API](#消息-api)
- [过滤器 API](#过滤器-api)
- [存储 API](#存储-api)
- [会话 API](#会话-api)
- [OneBot API](#onebot-api)

---

## 核心 API

### Run

运行机器人并返回管理器。

```go
func Run(cfg *Config) (*BotManager, error)
```

**参数：**
- `cfg`: 机器人配置

**返回：**
- `*BotManager`: 机器人管理器
- `error`: 错误信息

**示例：**
```go
manager, err := xbot.Run(cfg)
if err != nil {
    panic(err)
}
```

### RunAndListen

运行机器人并阻塞等待退出信号。

```go
func RunAndListen(cfg *Config) error
```

**参数：**
- `cfg`: 机器人配置

**返回：**
- `error`: 错误信息

**示例：**
```go
err := xbot.RunAndListen(cfg)
if err != nil {
    panic(err)
}
```

### LoadConfigFile

从文件加载配置。

```go
func LoadConfigFile(path string) (*Config, error)
```

**参数：**
- `path`: 配置文件路径

**返回：**
- `*Config`: 配置对象
- `error`: 错误信息

**示例：**
```go
cfg, err := xbot.LoadConfigFile("./config/config.yaml")
if err != nil {
    panic(err)
}
```

### GetStorage

获取插件专用存储。

```go
func GetStorage(pluginName string) storage.Storage
```

**参数：**
- `pluginName`: 插件名称

**返回：**
- `storage.Storage`: 存储实例

**示例：**
```go
storage := xbot.GetStorage("myplugin")
```

---

## 引擎 API

### NewEngine

创建新引擎。

```go
func NewEngine() *Engine
```

**返回：**
- `*Engine`: 引擎实例

**示例：**
```go
engine := xbot.NewEngine()
```

### Use

添加全局中间件。

```go
func (e *Engine) Use(middlewares ...func(next func(*Context)) func(*Context)) *Engine
```

**参数：**
- `middlewares`: 中间件函数列表

**返回：**
- `*Engine`: 引擎实例（支持链式调用）

**示例：**
```go
engine.Use(middleware1, middleware2)
```

### OnCommand

命令匹配器。

```go
func (e *Engine) OnCommand(command string, filters ...Filter) *Matcher
```

**参数：**
- `command`: 命令名称（不含前缀）
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnCommand("help").Handle(handler)
engine.OnCommand("admin", xbot.IsSuperUser()).Handle(handler)
```

### OnCommandGroup

命令组匹配器。

```go
func (e *Engine) OnCommandGroup(commands []string, filters ...Filter) *Matcher
```

**参数：**
- `commands`: 命令名称列表
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnCommandGroup([]string{"start", "begin"}).Handle(handler)
```

### OnKeywords

关键词匹配器。

```go
func (e *Engine) OnKeywords(keywords []string, filters ...Filter) *Matcher
```

**参数：**
- `keywords`: 关键词列表
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnKeywords([]string{"你好", "hello"}).Handle(handler)
```

### OnRegex

正则表达式匹配器。

```go
func (e *Engine) OnRegex(pattern string, filters ...Filter) *Matcher
```

**参数：**
- `pattern`: 正则表达式
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnRegex(`^抽卡\s*(\d+)次$`).Handle(handler)
```

### OnPrefix

前缀匹配器。

```go
func (e *Engine) OnPrefix(prefix string, filters ...Filter) *Matcher
```

**参数：**
- `prefix`: 前缀字符串
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnPrefix("查询").Handle(handler)
```

### OnSuffix

后缀匹配器。

```go
func (e *Engine) OnSuffix(suffix string, filters ...Filter) *Matcher
```

**参数：**
- `suffix`: 后缀字符串
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnSuffix("天气").Handle(handler)
```

### OnFullMatch

完全匹配器。

```go
func (e *Engine) OnFullMatch(text string, filters ...Filter) *Matcher
```

**参数：**
- `text`: 要匹配的完整文本
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnFullMatch("签到").Handle(handler)
```

### OnDFAKeywords

DFA 关键词匹配器（高性能）。

```go
func (e *Engine) OnDFAKeywords(provider VersionedKeywordProvider, filters ...Filter) *Matcher
```

**参数：**
- `provider`: 关键词提供者
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
manager := xbot.NewKeywordManager([]string{"敏感词"})
engine.OnDFAKeywords(manager).Handle(handler)
```

### OnACKeywords

AC 自动机关键词匹配器（更高性能）。

```go
func (e *Engine) OnACKeywords(provider VersionedKeywordProvider, filters ...Filter) *Matcher
```

**参数：**
- `provider`: 关键词提供者
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
manager := xbot.NewKeywordManager([]string{"违禁词"})
engine.OnACKeywords(manager).Handle(handler)
```

### OnMessage

消息事件匹配器。

```go
func (e *Engine) OnMessage(filters ...Filter) *Matcher
```

**参数：**
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnMessage(xbot.IsPrivate()).Handle(handler)
```

### OnNotice

通知事件匹配器。

```go
func (e *Engine) OnNotice(filters ...Filter) *Matcher
```

**参数：**
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnNotice().Handle(handler)
```

### OnRequest

请求事件匹配器。

```go
func (e *Engine) OnRequest(filters ...Filter) *Matcher
```

**参数：**
- `filters`: 可选的过滤器列表

**返回：**
- `*Matcher`: 匹配器实例

**示例：**
```go
engine.OnRequest().Handle(handler)
```

---

## 匹配器 API

### Handle

设置处理函数。

```go
func (m *Matcher) Handle(handler interface{}) *Matcher
```

**参数：**
- `handler`: 处理函数，类型为 `func(*Context)`

**返回：**
- `*Matcher`: 匹配器实例（支持链式调用）

**示例：**
```go
matcher.Handle(func(ctx *xbot.Context) {
    ctx.Reply("处理完成")
})
```

### Filter

添加过滤器。

```go
func (m *Matcher) Filter(filters ...Filter) *Matcher
```

**参数：**
- `filters`: 过滤器列表

**返回：**
- `*Matcher`: 匹配器实例（支持链式调用）

**示例：**
```go
matcher.Filter(xbot.IsGroup(), xbot.IsGroupAdmin())
```

### Limit

设置内存限流。

```go
func (m *Matcher) Limit(duration time.Duration, count int, onExceed func(*Context)) *Matcher
```

**参数：**
- `duration`: 时间窗口
- `count`: 最大调用次数
- `onExceed`: 超限回调函数

**返回：**
- `*Matcher`: 匹配器实例（支持链式调用）

**示例：**
```go
matcher.Limit(10*time.Second, 3, func(ctx *xbot.Context) {
    ctx.Reply("调用过于频繁")
})
```

### LimitWithRedis

使用 Redis 限流。

```go
func (m *Matcher) LimitWithRedis(limiter Limiter) *Matcher
```

**参数：**
- `limiter`: Redis 限流器实例

**返回：**
- `*Matcher`: 匹配器实例（支持链式调用）

**示例：**
```go
limiter := xbot.NewRedisLimiter(redis, "key", 10*time.Second, 3, onExceed)
matcher.LimitWithRedis(limiter)
```

### Priority

设置优先级。

```go
func (m *Matcher) Priority(p int) *Matcher
```

**参数：**
- `p`: 优先级（数值越大优先级越高）

**返回：**
- `*Matcher`: 匹配器实例（支持链式调用）

**示例：**
```go
matcher.Priority(100)
```

### Use

添加匹配器级别中间件。

```go
func (m *Matcher) Use(middlewares ...func(next func(*Context)) func(*Context)) *Matcher
```

**参数：**
- `middlewares`: 中间件函数列表

**返回：**
- `*Matcher`: 匹配器实例（支持链式调用）

**示例：**
```go
matcher.Use(middleware1, middleware2)
```

### SetBlock

设置是否阻止后续匹配器。

```go
func (m *Matcher) SetBlock(block ...bool) *Matcher
```

**参数：**
- `block`: 是否阻止（可选，默认 true）

**返回：**
- `*Matcher`: 匹配器实例（支持链式调用）

**示例：**
```go
// 阻止后续匹配
matcher.SetBlock()

// 不阻止
matcher.SetBlock(false)
```

---

## 上下文 API

### 信息获取

#### GetUserID

获取用户 ID。

```go
func (ctx *Context) GetUserID() int64
```

**返回：**
- `int64`: 用户 ID

#### GetGroupID

获取群号。

```go
func (ctx *Context) GetGroupID() int64
```

**返回：**
- `int64`: 群号（私聊时为 0）

#### GetMessageID

获取消息 ID。

```go
func (ctx *Context) GetMessageID() int64
```

**返回：**
- `int64`: 消息 ID

#### GetPlainText

获取纯文本消息（不含 CQ 码）。

```go
func (ctx *Context) GetPlainText() string
```

**返回：**
- `string`: 纯文本

#### GetRawMessage

获取原始消息（含 CQ 码）。

```go
func (ctx *Context) GetRawMessage() string
```

**返回：**
- `string`: 原始消息

#### GetArgs

获取命令参数（去除命令本身）。

```go
func (ctx *Context) GetArgs() string
```

**返回：**
- `string`: 命令参数

**示例：**
```go
// 消息: "/echo hello world"
args := ctx.GetArgs() // "hello world"
```

#### GetMessage

获取消息对象。

```go
func (ctx *Context) GetMessage() *message.Message
```

**返回：**
- `*message.Message`: 消息对象

#### GetAtUsers

获取被 @ 的用户列表。

```go
func (ctx *Context) GetAtUsers() []int64
```

**返回：**
- `[]int64`: 用户 ID 列表

### 消息发送

#### Reply

快速回复消息。

```go
func (ctx *Context) Reply(msg interface{})
```

**参数：**
- `msg`: 消息内容（string 或 *message.Message）

**示例：**
```go
ctx.Reply("Hello")
ctx.Reply(message.NewBuilder().Text("Hello").Build())
```

#### SendPrivateMessage

发送私聊消息。

```go
func (ctx *Context) SendPrivateMessage(userID int64, msg interface{})
```

**参数：**
- `userID`: 用户 ID
- `msg`: 消息内容

#### SendGroupMessage

发送群消息。

```go
func (ctx *Context) SendGroupMessage(groupID int64, msg interface{})
```

**参数：**
- `groupID`: 群号
- `msg`: 消息内容

### 消息操作

#### Delete

撤回当前消息。

```go
func (ctx *Context) Delete()
```

### 群组操作

#### SetGroupKick

踢出群成员。

```go
func (ctx *Context) SetGroupKick(groupID, userID int64, rejectAddRequest bool)
```

**参数：**
- `groupID`: 群号
- `userID`: 用户 ID
- `rejectAddRequest`: 是否拒绝再次加群

#### SetGroupBan

禁言群成员。

```go
func (ctx *Context) SetGroupBan(groupID, userID int64, duration int64)
```

**参数：**
- `groupID`: 群号
- `userID`: 用户 ID
- `duration`: 禁言时长（秒，0 表示解除禁言）

#### SetGroupWholeBan

全体禁言。

```go
func (ctx *Context) SetGroupWholeBan(groupID int64, enable bool)
```

**参数：**
- `groupID`: 群号
- `enable`: 是否启用

#### SetGroupCard

设置群名片。

```go
func (ctx *Context) SetGroupCard(groupID, userID int64, card string)
```

**参数：**
- `groupID`: 群号
- `userID`: 用户 ID
- `card`: 群名片

#### SetGroupAdmin

设置管理员。

```go
func (ctx *Context) SetGroupAdmin(groupID, userID int64, enable bool)
```

**参数：**
- `groupID`: 群号
- `userID`: 用户 ID
- `enable`: 是否设置为管理员

### 权限判断

#### IsAdmin

判断当前用户是否为管理员或群主。

```go
func (ctx *Context) IsAdmin() bool
```

**返回：**
- `bool`: 是否为管理员

#### IsOwner

判断当前用户是否为群主。

```go
func (ctx *Context) IsOwner() bool
```

**返回：**
- `bool`: 是否为群主

### 流程控制

#### Abort

中止后续所有匹配器。

```go
func (ctx *Context) Abort()
```

#### IsAborted

判断是否已中止。

```go
func (ctx *Context) IsAborted() bool
```

**返回：**
- `bool`: 是否已中止

---

## 消息 API

### NewBuilder

创建消息构建器。

```go
func NewBuilder() *Builder
```

**返回：**
- `*Builder`: 消息构建器

**示例：**
```go
msg := message.NewBuilder().
    Text("Hello ").
    At(123456).
    Image("https://example.com/image.jpg").
    Build()
```

### Builder 方法

#### Text

添加文本。

```go
func (b *Builder) Text(text string) *Builder
```

#### At

@ 某人。

```go
func (b *Builder) At(userID int64) *Builder
```

#### AtAll

@ 全体成员。

```go
func (b *Builder) AtAll() *Builder
```

#### Image

添加图片。

```go
func (b *Builder) Image(file string) *Builder
```

#### Record

添加语音。

```go
func (b *Builder) Record(file string) *Builder
```

#### Video

添加视频。

```go
func (b *Builder) Video(file string) *Builder
```

#### Face

添加表情。

```go
func (b *Builder) Face(id int) *Builder
```

#### Reply

回复消息。

```go
func (b *Builder) Reply(messageID int64) *Builder
```

#### Build

构建消息。

```go
func (b *Builder) Build() *Message
```

---

## 过滤器 API

### IsPrivate

私聊消息过滤器。

```go
func IsPrivate() Filter
```

### IsGroup

群聊消息过滤器。

```go
func IsGroup() Filter
```

### IsInGroup

指定群组过滤器。

```go
func IsInGroup(groupIDs ...int64) Filter
```

### IsSuperUser

超级用户过滤器。

```go
func IsSuperUser() Filter
```

### IsGroupAdmin

群管理员过滤器（含群主）。

```go
func IsGroupAdmin() Filter
```

### IsGroupOwner

群主过滤器。

```go
func IsGroupOwner() Filter
```

### ToMe

被 @ 或私聊过滤器。

```go
func ToMe() Filter
```

### Or

或运算过滤器。

```go
func Or(filters ...Filter) Filter
```

### And

与运算过滤器。

```go
func And(filters ...Filter) Filter
```

### Not

非运算过滤器。

```go
func Not(filter Filter) Filter
```

---

## 存储 API

### Set

存储数据。

```go
func (s Storage) Set(key, value []byte) error
```

### Get

获取数据。

```go
func (s Storage) Get(key []byte) ([]byte, error)
```

### Delete

删除数据。

```go
func (s Storage) Delete(key []byte) error
```

### Has

判断键是否存在。

```go
func (s Storage) Has(key []byte) (bool, error)
```

### Close

关闭存储。

```go
func (s Storage) Close() error
```

---

## 会话 API

### Wait

等待用户输入。

```go
func (m *Manager) Wait(userID, groupID int64, timeout time.Duration) *Context
```

**参数：**
- `userID`: 用户 ID
- `groupID`: 群号（私聊时为 0）
- `timeout`: 超时时间

**返回：**
- `*Context`: 上下文（超时返回 nil）

**示例：**
```go
ctx.Reply("请输入内容：")
inputCtx := ctx.Session.Wait(ctx.GetUserID(), ctx.GetGroupID(), 30*time.Second)
if inputCtx == nil {
    ctx.Reply("超时")
    return
}
input := inputCtx.GetPlainText()
```

---

## OneBot API

通过 `ctx.API` 访问 OneBot API。

### 消息相关

#### SendPrivateMsg

发送私聊消息。

```go
func (c *Client) SendPrivateMsg(userID int64, message string) error
```

#### SendGroupMsg

发送群消息。

```go
func (c *Client) SendGroupMsg(groupID int64, message string) error
```

#### DeleteMsg

撤回消息。

```go
func (c *Client) DeleteMsg(messageID int64) error
```

### 群组相关

#### SetGroupKick

踢出群成员。

```go
func (c *Client) SetGroupKick(groupID, userID int64, rejectAddRequest bool) error
```

#### SetGroupBan

禁言。

```go
func (c *Client) SetGroupBan(groupID, userID int64, duration int64) error
```

#### SetGroupWholeBan

全体禁言。

```go
func (c *Client) SetGroupWholeBan(groupID int64, enable bool) error
```

#### SetGroupCard

设置群名片。

```go
func (c *Client) SetGroupCard(groupID, userID int64, card string) error
```

#### SetGroupAdmin

设置管理员。

```go
func (c *Client) SetGroupAdmin(groupID, userID int64, enable bool) error
```

### 信息获取

#### GetLoginInfo

获取登录信息。

```go
func (c *Client) GetLoginInfo() (*types.LoginInfo, error)
```

#### GetGroupInfo

获取群信息。

```go
func (c *Client) GetGroupInfo(groupID int64, noCache bool) (*types.GroupInfo, error)
```

#### GetGroupMemberInfo

获取群成员信息。

```go
func (c *Client) GetGroupMemberInfo(groupID, userID int64, noCache bool) (*types.GroupMemberInfo, error)
```

#### GetGroupMemberList

获取群成员列表。

```go
func (c *Client) GetGroupMemberList(groupID int64) ([]*types.GroupMemberInfo, error)
```

---

## 关键词管理 API

### NewKeywordManager

创建关键词管理器。

```go
func NewKeywordManager(keywords []string) *KeywordManager
```

**参数：**
- `keywords`: 初始关键词列表

**返回：**
- `*KeywordManager`: 关键词管理器

### AddKeyword

添加关键词。

```go
func (km *KeywordManager) AddKeyword(keyword string)
```

### RemoveKeyword

移除关键词。

```go
func (km *KeywordManager) RemoveKeyword(keyword string)
```

### GetKeywords

获取所有关键词。

```go
func (km *KeywordManager) GetKeywords() []string
```

### GetVersion

获取版本号。

```go
func (km *KeywordManager) GetVersion() int64
```

---

完整的 API 列表请参考源代码和 [README.md](./README.md)。

