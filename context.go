package xbot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xiaoyi510/xbot/event"
	"github.com/xiaoyi510/xbot/logger"
	"github.com/xiaoyi510/xbot/message"
	"github.com/xiaoyi510/xbot/session"
	"github.com/xiaoyi510/xbot/storage"
)

// RegexMatch 正则匹配结果
type RegexMatch struct {
	// Match 匹配到的完整字符串
	Match string
	// Groups 按索引存储的分组（索引0是完整匹配，1开始是分组）
	Groups []string
	// NamedGroups 按名称存储的分组
	NamedGroups map[string]string
}

// Context 上下文
type Context struct {
	Event   event.Event
	Bot     *Bot
	State   map[string]interface{}
	Logger  logger.Logger
	Storage storage.Storage
	Session *session.Manager

	// RegexResult 正则匹配结果
	RegexResult *RegexMatch

	matched        bool
	aborted        bool // 是否中止后续匹配器
	shouldContinue bool // 是否显式调用Next()继续
}

// NewContext 创建上下文
func NewContext(evt event.Event, bot *Bot) *Context {
	return &Context{
		Event:          evt,
		Bot:            bot,
		State:          make(map[string]interface{}),
		Logger:         logger.GetDefault().WithField("selfID", evt.GetSelfID()),
		Storage:        bot.Storage,
		Session:        bot.SessionManager,
		matched:        false,
		aborted:        false,
		shouldContinue: true, // 默认继续执行后续匹配器
	}
}

// GetUserID 获取用户 ID
func (ctx *Context) GetUserID() int64 {
	switch evt := ctx.Event.(type) {
	case *event.PrivateMessageEvent:
		return evt.UserID
	case *event.GroupMessageEvent:
		return evt.UserID
	case *event.FriendRequestEvent:
		return evt.UserID
	case *event.GroupRequestEvent:
		return evt.UserID
	default:
		return 0
	}
}

// GetGroupID 获取群号
func (ctx *Context) GetGroupID() int64 {
	switch evt := ctx.Event.(type) {
	case *event.GroupMessageEvent:
		return evt.GroupID
	case *event.GroupRequestEvent:
		return evt.GroupID
	default:
		return 0
	}
}

// GetMessageID 获取消息 ID
func (ctx *Context) GetMessageID() int64 {
	switch evt := ctx.Event.(type) {
	case *event.PrivateMessageEvent:
		return evt.MessageID
	case *event.GroupMessageEvent:
		return evt.MessageID
	default:
		return 0
	}
}

// GetMessage 获取消息
func (ctx *Context) GetMessage() *message.Message {
	switch evt := ctx.Event.(type) {
	case *event.PrivateMessageEvent:
		return &evt.ParsedMessage
	case *event.GroupMessageEvent:
		return &evt.ParsedMessage
	default:
		return nil
	}
}

// GetPlainText 获取纯文本消息
func (ctx *Context) GetPlainText() string {
	msg := ctx.GetMessage()
	if msg == nil {
		return ""
	}
	return msg.GetPlainText()
}

// GetRawMessage 获取原始消息（含 CQ 码）
func (ctx *Context) GetRawMessage() string {
	msg := ctx.GetMessage()
	if msg == nil {
		return ""
	}
	return msg.GetRawMessage()
}

// GetArgs 获取命令参数
// 返回去除命令前缀和命令本身后的参数部分
// 例如："/echo hello world" 返回 "hello world"
func (ctx *Context) GetArgs() string {
	text := ctx.GetPlainText()
	if text == "" {
		return ""
	}

	// 如果有命令前缀，尝试去除
	prefix := ctx.Bot.Config.CommandPrefix
	if prefix != "" && len(text) > len(prefix) && text[:len(prefix)] == prefix {
		text = text[len(prefix):]
	}

	// 查找第一个空格，去除命令本身
	for i, ch := range text {
		if ch == ' ' || ch == '\t' || ch == '\n' {
			// 跳过连续的空白字符
			for i < len(text) && (text[i] == ' ' || text[i] == '\t' || text[i] == '\n') {
				i++
			}
			if i < len(text) {
				return text[i:]
			}
			return ""
		}
	}

	// 没有找到空格，说明没有参数
	return ""
}

// GetAtUsers 获取被 @ 的用户列表
func (ctx *Context) GetAtUsers() []int64 {
	msg := ctx.GetMessage()
	if msg == nil {
		return nil
	}

	var users []int64
	for _, seg := range *msg {
		if seg.Type == "at" {
			if qq, ok := seg.Data["qq"].(string); ok {
				// 将字符串转为 int64
				var userID int64
				fmt.Sscanf(qq, "%d", &userID)
				if userID > 0 {
					users = append(users, userID)
				}
			} else if qq, ok := seg.Data["qq"].(float64); ok {
				users = append(users, int64(qq))
			} else if qq, ok := seg.Data["qq"].(int64); ok {
				users = append(users, qq)
			}
		}
	}
	return users
}

// Reply 回复消息
// 返回值：消息ID, 错误
// 消息ID可用于后续操作（如撤回、设置精华等）
func (ctx *Context) Reply(msg interface{}) (int64, error) {
	var messageData interface{}

	// 转换消息类型
	switch m := msg.(type) {
	case string:
		messageData = []message.MessageSegment{message.Text(m)}
	case message.Message:
		messageData = m
	case []message.MessageSegment:
		messageData = m
	case message.MessageSegment:
		messageData = []message.MessageSegment{m}
	default:
		messageData = msg
	}

	switch evt := ctx.Event.(type) {
	case *event.PrivateMessageEvent:
		resp, err := ctx.Bot.API.SendPrivateMsg(evt.UserID, messageData)
		if err != nil {
			return 0, err
		}
		return resp.Data.MessageID, nil
	case *event.GroupMessageEvent:
		resp, err := ctx.Bot.API.SendGroupMsg(evt.GroupID, messageData)
		if err != nil {
			return 0, err
		}
		return resp.Data.MessageID, nil
	default:
		return 0, nil
	}
}

// ReplyText 回复文本消息
// 返回值：消息ID, 错误
func (ctx *Context) ReplyText(text string) (int64, error) {
	return ctx.Reply(text)
}

// WaitNextMessage 等待下一条消息
func (ctx *Context) WaitNextMessage(timeout time.Duration) (*Context, error) {
	userID := ctx.GetUserID()
	groupID := ctx.GetGroupID()

	// 创建等待会话
	sess := ctx.Session.CreateWaitSession(userID, groupID, timeout)

	// 等待响应
	data, err := sess.Wait(timeout)
	if err != nil {
		return nil, err
	}

	// 返回新的上下文
	if nextCtx, ok := data.(*Context); ok {
		return nextCtx, nil
	}

	return nil, nil
}

// SaveData 保存数据到存储
func (ctx *Context) SaveData(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ctx.Storage.Set(key, data)
}

// LoadData 从存储加载数据
func (ctx *Context) LoadData(key string, value interface{}) error {
	data, err := ctx.Storage.Get(key)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, value)
}

// Set 设置状态
func (ctx *Context) Set(key string, value interface{}) {
	ctx.State[key] = value
}

// Get 获取状态
func (ctx *Context) Get(key string) (interface{}, bool) {
	value, ok := ctx.State[key]
	return value, ok
}

// IsSuperUser 是否是超级用户
func (ctx *Context) IsSuperUser() bool {
	userID := ctx.GetUserID()
	for _, su := range ctx.Bot.Config.SuperUsers {
		if userID == su {
			return true
		}
	}
	return false
}

// ========== Event 辅助方法 ==========

// GroupEvent 获取群消息事件
// 返回类型断言后的群消息事件和是否成功
func (ctx *Context) GroupEvent() (*event.GroupMessageEvent, bool) {
	evt, ok := ctx.Event.(*event.GroupMessageEvent)
	return evt, ok
}

// PrivateEvent 获取私聊消息事件
// 返回类型断言后的私聊消息事件和是否成功
func (ctx *Context) PrivateEvent() (*event.PrivateMessageEvent, bool) {
	evt, ok := ctx.Event.(*event.PrivateMessageEvent)
	return evt, ok
}

// GroupRequestEvent 获取群请求事件
// 返回类型断言后的群请求事件和是否成功
func (ctx *Context) GroupRequestEvent() (*event.GroupRequestEvent, bool) {
	evt, ok := ctx.Event.(*event.GroupRequestEvent)
	return evt, ok
}

// FriendRequestEvent 获取好友请求事件
// 返回类型断言后的好友请求事件和是否成功
func (ctx *Context) FriendRequestEvent() (*event.FriendRequestEvent, bool) {
	evt, ok := ctx.Event.(*event.FriendRequestEvent)
	return evt, ok
}

// IsGroupMessage 判断是否是群消息
func (ctx *Context) IsGroupMessage() bool {
	_, ok := ctx.Event.(*event.GroupMessageEvent)
	return ok
}

// IsPrivateMessage 判断是否是私聊消息
func (ctx *Context) IsPrivateMessage() bool {
	_, ok := ctx.Event.(*event.PrivateMessageEvent)
	return ok
}

// MustGroupEvent 获取群消息事件（必须是群消息，否则panic）
// 仅在确定是群消息的Handler中使用
func (ctx *Context) MustGroupEvent() *event.GroupMessageEvent {
	if evt, ok := ctx.Event.(*event.GroupMessageEvent); ok {
		return evt
	}
	panic("Context.MustGroupEvent: event is not a GroupMessageEvent")
}

// MustPrivateEvent 获取私聊消息事件（必须是私聊消息，否则panic）
// 仅在确定是私聊消息的Handler中使用
func (ctx *Context) MustPrivateEvent() *event.PrivateMessageEvent {
	if evt, ok := ctx.Event.(*event.PrivateMessageEvent); ok {
		return evt
	}
	panic("Context.MustPrivateEvent: event is not a PrivateMessageEvent")
}

// ========== 流程控制方法 ==========

// Next 显式标记继续执行后续匹配器
// 用于表示当前Handler没有处理该消息，应该继续匹配其他Handler
//
// 使用场景：
//   - 可选功能（如关键词回复）没有匹配到时
//   - 条件判断后决定不处理时
//
// 示例：
//
//	engine.OnMessage(...).Handle(func(ctx *xbot.Context) {
//	    if !hasConfig {
//	        ctx.Next()  // 没有配置，继续其他匹配器
//	        return
//	    }
//	    if !matched {
//	        ctx.Next()  // 没有匹配，继续其他匹配器
//	        return
//	    }
//	    // 处理逻辑
//	})
func (ctx *Context) Next() {
	ctx.shouldContinue = true
}

// Abort 中止后续匹配器的执行
// 用于表示当前Handler已经完全处理了该消息，不需要其他Handler再处理
//
// 注意：由于使用goroutine并发执行，Abort()只能阻止尚未启动的匹配器，
// 已经启动的匹配器仍会执行完成。
//
// 使用场景：
//   - 命令已处理，不需要其他Handler
//   - 敏感词已撤回，不需要继续处理
//   - 通常配合SetBlock()使用效果更好
//
// 示例：
//
//	engine.OnMessage(...).Handle(func(ctx *xbot.Context) {
//	    if matched {
//	        ctx.Reply("已处理")
//	        ctx.Abort()  // 中止后续匹配器
//	        return
//	    }
//	})
func (ctx *Context) Abort() {
	ctx.aborted = true
	ctx.shouldContinue = false
}

// IsAborted 检查是否已中止
func (ctx *Context) IsAborted() bool {
	return ctx.aborted
}

// ========== 消息操作方法 ==========

// Delete 撤回当前消息
func (ctx *Context) Delete() error {
	messageID := ctx.GetMessageID()
	if messageID == 0 {
		return fmt.Errorf("无法获取消息ID")
	}
	return ctx.Bot.API.DeleteMsg(messageID)
}

// SendPrivateMessage 发送私聊消息
func (ctx *Context) SendPrivateMessage(userID int64, msg interface{}) (int64, error) {
	var messageData interface{}

	// 转换消息类型
	switch m := msg.(type) {
	case string:
		messageData = []message.MessageSegment{message.Text(m)}
	case message.Message:
		messageData = m
	case []message.MessageSegment:
		messageData = m
	case message.MessageSegment:
		messageData = []message.MessageSegment{m}
	default:
		messageData = msg
	}

	resp, err := ctx.Bot.API.SendPrivateMsg(userID, messageData)
	if err != nil {
		return 0, err
	}
	return resp.Data.MessageID, nil
}

// SendGroupMessage 发送群消息
func (ctx *Context) SendGroupMessage(groupID int64, msg interface{}) (int64, error) {
	var messageData interface{}

	// 转换消息类型
	switch m := msg.(type) {
	case string:
		messageData = []message.MessageSegment{message.Text(m)}
	case message.Message:
		messageData = m
	case []message.MessageSegment:
		messageData = m
	case message.MessageSegment:
		messageData = []message.MessageSegment{m}
	default:
		messageData = msg
	}

	resp, err := ctx.Bot.API.SendGroupMsg(groupID, messageData)
	if err != nil {
		return 0, err
	}
	return resp.Data.MessageID, nil
}

// ========== 群组操作方法 ==========

// SetGroupKick 踢出群成员
func (ctx *Context) SetGroupKick(groupID, userID int64, rejectAddRequest bool) error {
	return ctx.Bot.API.SetGroupKick(groupID, userID, rejectAddRequest)
}

// SetGroupBan 禁言群成员
// duration: 禁言时长（秒），0 表示解除禁言
func (ctx *Context) SetGroupBan(groupID, userID int64, duration int64) error {
	return ctx.Bot.API.SetGroupBan(groupID, userID, int32(duration))
}

// SetGroupWholeBan 全体禁言
func (ctx *Context) SetGroupWholeBan(groupID int64, enable bool) error {
	return ctx.Bot.API.SetGroupWholeBan(groupID, enable)
}

// SetGroupCard 设置群名片
func (ctx *Context) SetGroupCard(groupID, userID int64, card string) error {
	return ctx.Bot.API.SetGroupCard(groupID, userID, card)
}

// SetGroupAdmin 设置群管理员
func (ctx *Context) SetGroupAdmin(groupID, userID int64, enable bool) error {
	return ctx.Bot.API.SetGroupAdmin(groupID, userID, enable)
}

// ========== 权限判断方法 ==========

// IsAdmin 判断当前用户是否为群管理员或群主
func (ctx *Context) IsAdmin() bool {
	if evt, ok := ctx.Event.(*event.GroupMessageEvent); ok {
		return evt.Sender.Role == "admin" || evt.Sender.Role == "owner"
	}
	return false
}

// IsOwner 判断当前用户是否为群主
func (ctx *Context) IsOwner() bool {
	if evt, ok := ctx.Event.(*event.GroupMessageEvent); ok {
		return evt.Sender.Role == "owner"
	}
	return false
}
