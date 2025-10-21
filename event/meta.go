package event

import (
	"github.com/xiaoyi510/xbot/types"
)

// LifecycleMetaEvent 生命周期元事件
type LifecycleMetaEvent struct {
	BaseEvent
	MetaEventType types.MetaEventType    `json:"meta_event_type"` // lifecycle
	SubType       types.LifecycleSubType `json:"sub_type"`        // 子类型
}

// IsEnable 是否是 OneBot 启用
func (e *LifecycleMetaEvent) IsEnable() bool {
	return e.SubType == types.LifecycleSubTypeEnable
}

// IsDisable 是否是 OneBot 停用
func (e *LifecycleMetaEvent) IsDisable() bool {
	return e.SubType == types.LifecycleSubTypeDisable
}

// IsConnect 是否是 WebSocket 连接成功
func (e *LifecycleMetaEvent) IsConnect() bool {
	return e.SubType == types.LifecycleSubTypeConnect
}

// HeartbeatMetaEvent 心跳元事件
type HeartbeatMetaEvent struct {
	BaseEvent
	MetaEventType types.MetaEventType `json:"meta_event_type"` // heartbeat
	Status        types.Status        `json:"status"`          // 状态信息
	Interval      int64               `json:"interval"`        // 心跳间隔（毫秒）
}

// IsOnline 是否在线
func (e *HeartbeatMetaEvent) IsOnline() bool {
	return e.Status.Online
}

// IsGood 状态是否正常
func (e *HeartbeatMetaEvent) IsGood() bool {
	return e.Status.Good
}
