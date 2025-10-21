package event

import (
	"encoding/json"
	"xbot/message"
	"xbot/types"
)

// Event 基础事件接口
type Event interface {
	GetTime() int64
	GetSelfID() int64
	GetPostType() string
}

// BaseEvent 基础事件
type BaseEvent struct {
	Time     int64          `json:"time"`      // 事件发生的时间戳
	SelfID   int64          `json:"self_id"`   // 收到事件的机器人 QQ 号
	PostType types.PostType `json:"post_type"` // 事件类型
}

// GetTime 获取事件时间
func (e *BaseEvent) GetTime() int64 {
	return e.Time
}

// GetSelfID 获取机器人 QQ 号
func (e *BaseEvent) GetSelfID() int64 {
	return e.SelfID
}

// GetPostType 获取事件类型
func (e *BaseEvent) GetPostType() string {
	return string(e.PostType)
}

// ParseEvent 解析事件
func ParseEvent(data []byte) (Event, error) {
	// 首先解析基础事件以确定类型
	var base BaseEvent
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	// 根据事件类型解析具体事件
	switch base.PostType {
	case types.PostTypeMessage:
		return parseMessageEvent(data)
	case types.PostTypeNotice:
		return parseNoticeEvent(data)
	case types.PostTypeRequest:
		return parseRequestEvent(data)
	case types.PostTypeMetaEvent:
		return parseMetaEvent(data)
	default:
		// 未知事件类型，返回基础事件
		return &base, nil
	}
}

// parseMessageEvent 解析消息事件
func parseMessageEvent(data []byte) (Event, error) {
	var temp struct {
		MessageType types.MessageType `json:"message_type"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	switch temp.MessageType {
	case types.MessageTypePrivate:
		var evt PrivateMessageEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			return nil, err
		}
		evt.ParsedMessage = message.ParseMessage(evt.Message)
		return &evt, nil
	case types.MessageTypeGroup:
		var evt GroupMessageEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			return nil, err
		}
		evt.ParsedMessage = message.ParseMessage(evt.Message)
		return &evt, nil
	default:
		return nil, json.Unmarshal(data, &BaseEvent{})
	}
}

// parseNoticeEvent 解析通知事件
func parseNoticeEvent(data []byte) (Event, error) {
	var temp struct {
		NoticeType types.NoticeType `json:"notice_type"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	switch temp.NoticeType {
	case types.NoticeTypeGroupUpload:
		var evt GroupUploadNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.NoticeTypeGroupAdmin:
		var evt GroupAdminNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.NoticeTypeGroupDecrease:
		var evt GroupDecreaseNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.NoticeTypeGroupIncrease:
		var evt GroupIncreaseNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.NoticeTypeGroupBan:
		var evt GroupBanNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.NoticeTypeFriendAdd:
		var evt FriendAddNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.NoticeTypeGroupRecall:
		var evt GroupRecallNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.NoticeTypeFriendRecall:
		var evt FriendRecallNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.NoticeTypeNotify:
		var evt NotifyNoticeEvent
		return &evt, json.Unmarshal(data, &evt)
	default:
		var evt BaseEvent
		return &evt, json.Unmarshal(data, &evt)
	}
}

// parseRequestEvent 解析请求事件
func parseRequestEvent(data []byte) (Event, error) {
	var temp struct {
		RequestType types.RequestType `json:"request_type"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	switch temp.RequestType {
	case types.RequestTypeFriend:
		var evt FriendRequestEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.RequestTypeGroup:
		var evt GroupRequestEvent
		return &evt, json.Unmarshal(data, &evt)
	default:
		var evt BaseEvent
		return &evt, json.Unmarshal(data, &evt)
	}
}

// parseMetaEvent 解析元事件
func parseMetaEvent(data []byte) (Event, error) {
	var temp struct {
		MetaEventType types.MetaEventType `json:"meta_event_type"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	switch temp.MetaEventType {
	case types.MetaEventTypeLifecycle:
		var evt LifecycleMetaEvent
		return &evt, json.Unmarshal(data, &evt)
	case types.MetaEventTypeHeartbeat:
		var evt HeartbeatMetaEvent
		return &evt, json.Unmarshal(data, &evt)
	default:
		var evt BaseEvent
		return &evt, json.Unmarshal(data, &evt)
	}
}
