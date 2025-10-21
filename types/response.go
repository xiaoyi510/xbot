package types

// Response OneBot API 响应结构（泛型）
type Response[T any] struct {
	Status  string `json:"status"`            // 状态：ok 表示成功，failed 表示失败
	RetCode int    `json:"retcode"`           // 返回码，0 表示成功，非 0 表示失败
	Data    T      `json:"data"`              // 数据
	Message string `json:"message,omitempty"` // 错误信息
	Wording string `json:"wording,omitempty"` // 错误信息（某些实现使用这个字段）
	Echo    string `json:"echo,omitempty"`    // 回声，用于识别 API 调用
}

// IsSuccess 判断 API 调用是否成功
func (r *Response[T]) IsSuccess() bool {
	return r.Status == "ok" && r.RetCode == 0
}

// GetError 获取错误信息
func (r *Response[T]) GetError() string {
	if r.Message != "" {
		return r.Message
	}
	return r.Wording
}

// APIResponse 通用 API 响应（使用 interface{}）
type APIResponse = Response[interface{}]

// EmptyResponse 空响应
type EmptyResponse = Response[struct{}]

// MessageResponse 发送消息响应
type MessageResponse struct {
	MessageID int64 `json:"message_id"` // 消息 ID
}

// ForwardMessageResponse 发送合并转发响应
type ForwardMessageResponse struct {
	MessageID int64  `json:"message_id"` // 消息 ID
	ForwardID string `json:"forward_id"` // 转发消息 ID
}

// ImageInfoResponse 图片信息响应
type ImageInfoResponse struct {
	File string `json:"file"` // 图片文件名
}

// StatusResponse 状态响应
type StatusResponse struct {
	Online bool `json:"online"` // 当前 QQ 在线
	Good   bool `json:"good"`   // 状态符合预期
	Stat   struct {
		PacketReceived  uint64 `json:"packet_received"`   // 收到的数据包总数
		PacketSent      uint64 `json:"packet_sent"`       // 发送的数据包总数
		PacketLost      uint32 `json:"packet_lost"`       // 数据包丢失总数
		MessageReceived uint64 `json:"message_received"`  // 接受信息总数
		MessageSent     uint64 `json:"message_sent"`      // 发送信息总数
		DisconnectTimes uint32 `json:"disconnect_times"`  // TCP 连接断开次数
		LostTimes       uint32 `json:"lost_times"`        // 账号掉线次数
		LastMessageTime int64  `json:"last_message_time"` // 最后消息时间
	} `json:"stat,omitempty"` // 统计信息
}

// EssenceMessageListResponse 精华消息列表响应
type EssenceMessageListResponse = Response[[]EssenceMessage]
