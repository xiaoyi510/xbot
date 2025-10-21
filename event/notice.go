package event

import (
	"xbot/types"
)

// GroupUploadNoticeEvent 群文件上传事件
type GroupUploadNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType `json:"notice_type"` // group_upload
	GroupID    int64            `json:"group_id"`    // 群号
	UserID     int64            `json:"user_id"`     // 发送者 QQ 号
	File       types.File       `json:"file"`        // 文件信息
}

// GroupAdminNoticeEvent 群管理员变动事件
type GroupAdminNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType `json:"notice_type"` // group_admin
	SubType    string           `json:"sub_type"`    // set 设置管理员, unset 取消管理员
	GroupID    int64            `json:"group_id"`    // 群号
	UserID     int64            `json:"user_id"`     // 管理员 QQ 号
}

// IsSet 是否是设置管理员
func (e *GroupAdminNoticeEvent) IsSet() bool {
	return e.SubType == "set"
}

// IsUnset 是否是取消管理员
func (e *GroupAdminNoticeEvent) IsUnset() bool {
	return e.SubType == "unset"
}

// GroupDecreaseNoticeEvent 群成员减少事件
type GroupDecreaseNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType `json:"notice_type"` // group_decrease
	SubType    string           `json:"sub_type"`    // leave 主动退群, kick 成员被踢, kick_me 登录号被踢
	GroupID    int64            `json:"group_id"`    // 群号
	OperatorID int64            `json:"operator_id"` // 操作者 QQ 号
	UserID     int64            `json:"user_id"`     // 离开者 QQ 号
}

// IsLeave 是否是主动退群
func (e *GroupDecreaseNoticeEvent) IsLeave() bool {
	return e.SubType == "leave"
}

// IsKick 是否是成员被踢
func (e *GroupDecreaseNoticeEvent) IsKick() bool {
	return e.SubType == "kick"
}

// IsKickMe 是否是登录号被踢
func (e *GroupDecreaseNoticeEvent) IsKickMe() bool {
	return e.SubType == "kick_me"
}

// GroupIncreaseNoticeEvent 群成员增加事件
type GroupIncreaseNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType `json:"notice_type"` // group_increase
	SubType    string           `json:"sub_type"`    // approve 管理员同意入群, invite 管理员邀请入群
	GroupID    int64            `json:"group_id"`    // 群号
	OperatorID int64            `json:"operator_id"` // 操作者 QQ 号
	UserID     int64            `json:"user_id"`     // 加入者 QQ 号
}

// IsApprove 是否是管理员同意入群
func (e *GroupIncreaseNoticeEvent) IsApprove() bool {
	return e.SubType == "approve"
}

// IsInvite 是否是管理员邀请入群
func (e *GroupIncreaseNoticeEvent) IsInvite() bool {
	return e.SubType == "invite"
}

// GroupBanNoticeEvent 群禁言事件
type GroupBanNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType `json:"notice_type"` // group_ban
	SubType    string           `json:"sub_type"`    // ban 禁言, lift_ban 解除禁言
	GroupID    int64            `json:"group_id"`    // 群号
	OperatorID int64            `json:"operator_id"` // 操作者 QQ 号
	UserID     int64            `json:"user_id"`     // 被禁言 QQ 号
	Duration   int64            `json:"duration"`    // 禁言时长（秒）
}

// IsBan 是否是禁言
func (e *GroupBanNoticeEvent) IsBan() bool {
	return e.SubType == "ban"
}

// IsLiftBan 是否是解除禁言
func (e *GroupBanNoticeEvent) IsLiftBan() bool {
	return e.SubType == "lift_ban"
}

// FriendAddNoticeEvent 好友添加事件
type FriendAddNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType `json:"notice_type"` // friend_add
	UserID     int64            `json:"user_id"`     // 新添加好友 QQ 号
}

// GroupRecallNoticeEvent 群消息撤回事件
type GroupRecallNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType `json:"notice_type"` // group_recall
	GroupID    int64            `json:"group_id"`    // 群号
	UserID     int64            `json:"user_id"`     // 消息发送者 QQ 号
	OperatorID int64            `json:"operator_id"` // 操作者 QQ 号
	MessageID  int64            `json:"message_id"`  // 被撤回的消息 ID
}

// FriendRecallNoticeEvent 好友消息撤回事件
type FriendRecallNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType `json:"notice_type"` // friend_recall
	UserID     int64            `json:"user_id"`     // 好友 QQ 号
	MessageID  int64            `json:"message_id"`  // 被撤回的消息 ID
}

// NotifyNoticeEvent 群内提示事件
type NotifyNoticeEvent struct {
	BaseEvent
	NoticeType types.NoticeType    `json:"notice_type"`          // notify
	SubType    types.NotifySubType `json:"sub_type"`             // 提示类型
	GroupID    int64               `json:"group_id"`             // 群号
	UserID     int64               `json:"user_id"`              // 发送者 QQ 号
	TargetID   int64               `json:"target_id,omitempty"`  // 被戳者 QQ 号
	HonorType  types.HonorType     `json:"honor_type,omitempty"` // 荣誉类型
}

// IsPoke 是否是戳一戳
func (e *NotifyNoticeEvent) IsPoke() bool {
	return e.SubType == types.NotifySubTypePoke
}

// IsHonor 是否是荣誉变更
func (e *NotifyNoticeEvent) IsHonor() bool {
	return e.SubType == types.NotifySubTypeHonor
}

// IsLuckyKing 是否是红包运气王
func (e *NotifyNoticeEvent) IsLuckyKing() bool {
	return e.SubType == types.NotifySubTypeLucky
}

// IsTitle 是否是头衔变更
func (e *NotifyNoticeEvent) IsTitle() bool {
	return e.SubType == types.NotifySubTypeTitle
}
