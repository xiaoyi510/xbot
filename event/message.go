package event

import (
	"xbot/message"
	"xbot/types"
)

// PrivateMessageEvent 私聊消息事件
type PrivateMessageEvent struct {
	BaseEvent
	MessageType types.MessageType           `json:"message_type"` // message
	SubType     types.PrivateMessageSubType `json:"sub_type"`     // 消息子类型
	MessageID   int64                       `json:"message_id"`   // 消息 ID
	UserID      int64                       `json:"user_id"`      // 发送者 QQ 号
	Message     interface{}                 `json:"message"`      // 消息内容
	RawMessage  string                      `json:"raw_message"`  // 原始消息内容
	Font        int32                       `json:"font"`         // 字体
	Sender      types.Sender                `json:"sender"`       // 发送者信息

	// 解析后的消息
	ParsedMessage message.Message `json:"-"`
}

// GetUserID 获取用户 ID
func (e *PrivateMessageEvent) GetUserID() int64 {
	return e.UserID
}

// GetMessageID 获取消息 ID
func (e *PrivateMessageEvent) GetMessageID() int64 {
	return e.MessageID
}

// GetMessage 获取解析后的消息
func (e *PrivateMessageEvent) GetMessage() message.Message {
	return e.ParsedMessage
}

// GetPlainText 获取纯文本消息
func (e *PrivateMessageEvent) GetPlainText() string {
	return e.ParsedMessage.GetPlainText()
}

// IsFromFriend 是否来自好友
func (e *PrivateMessageEvent) IsFromFriend() bool {
	return e.SubType == types.PrivateMessageSubTypeFriend
}

// IsFromGroup 是否来自群临时会话
func (e *PrivateMessageEvent) IsFromGroup() bool {
	return e.SubType == types.PrivateMessageSubTypeGroup
}

// GroupMessageEvent 群消息事件
type GroupMessageEvent struct {
	BaseEvent
	MessageType types.MessageType         `json:"message_type"` // message
	SubType     types.GroupMessageSubType `json:"sub_type"`     // 消息子类型
	MessageID   int64                     `json:"message_id"`   // 消息 ID
	GroupID     int64                     `json:"group_id"`     // 群号
	UserID      int64                     `json:"user_id"`      // 发送者 QQ 号
	Anonymous   *types.Anonymous          `json:"anonymous"`    // 匿名信息
	Message     interface{}               `json:"message"`      // 消息内容
	RawMessage  string                    `json:"raw_message"`  // 原始消息内容
	Font        int32                     `json:"font"`         // 字体
	Sender      types.Sender              `json:"sender"`       // 发送者信息

	// 解析后的消息
	ParsedMessage message.Message `json:"-"`
}

// GetUserID 获取用户 ID
func (e *GroupMessageEvent) GetUserID() int64 {
	return e.UserID
}

// GetGroupID 获取群号
func (e *GroupMessageEvent) GetGroupID() int64 {
	return e.GroupID
}

// GetMessageID 获取消息 ID
func (e *GroupMessageEvent) GetMessageID() int64 {
	return e.MessageID
}

// GetMessage 获取解析后的消息
func (e *GroupMessageEvent) GetMessage() message.Message {
	return e.ParsedMessage
}

// GetPlainText 获取纯文本消息
func (e *GroupMessageEvent) GetPlainText() string {
	return e.ParsedMessage.GetPlainText()
}

// IsAnonymous 是否是匿名消息
func (e *GroupMessageEvent) IsAnonymous() bool {
	return e.Anonymous != nil && e.SubType == types.GroupMessageSubTypeAnonymous
}

// IsNormal 是否是正常消息
func (e *GroupMessageEvent) IsNormal() bool {
	return e.SubType == types.GroupMessageSubTypeNormal
}

// IsNotice 是否是系统提示
func (e *GroupMessageEvent) IsNotice() bool {
	return e.SubType == types.GroupMessageSubTypeNotice
}

// GetRole 获取发送者角色
func (e *GroupMessageEvent) GetRole() types.Role {
	return e.Sender.Role
}

// IsOwner 是否是群主
func (e *GroupMessageEvent) IsOwner() bool {
	return e.Sender.Role == types.RoleOwner
}

// IsAdmin 是否是管理员
func (e *GroupMessageEvent) IsAdmin() bool {
	return e.Sender.Role == types.RoleAdmin
}

// IsMember 是否是普通成员
func (e *GroupMessageEvent) IsMember() bool {
	return e.Sender.Role == types.RoleMember
}
