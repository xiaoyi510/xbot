package event

import (
	"xbot/types"
)

// FriendRequestEvent 加好友请求事件
type FriendRequestEvent struct {
	BaseEvent
	RequestType types.RequestType `json:"request_type"` // friend
	UserID      int64             `json:"user_id"`      // 发送请求的 QQ 号
	Comment     string            `json:"comment"`      // 验证信息
	Flag        string            `json:"flag"`         // 请求 flag，在调用处理请求的 API 时需要传入
}

// GroupRequestEvent 加群请求事件
type GroupRequestEvent struct {
	BaseEvent
	RequestType types.RequestType         `json:"request_type"` // group
	SubType     types.GroupRequestSubType `json:"sub_type"`     // add 加群请求, invite 邀请登录号入群
	GroupID     int64                     `json:"group_id"`     // 群号
	UserID      int64                     `json:"user_id"`      // 发送请求的 QQ 号
	Comment     string                    `json:"comment"`      // 验证信息
	Flag        string                    `json:"flag"`         // 请求 flag
}

// IsAdd 是否是加群请求
func (e *GroupRequestEvent) IsAdd() bool {
	return e.SubType == types.GroupRequestSubTypeAdd
}

// IsInvite 是否是邀请入群
func (e *GroupRequestEvent) IsInvite() bool {
	return e.SubType == types.GroupRequestSubTypeInvite
}
