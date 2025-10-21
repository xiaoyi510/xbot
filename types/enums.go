package types

// PostType 事件类型
type PostType string

const (
	PostTypeMessage   PostType = "message"    // 消息事件
	PostTypeNotice    PostType = "notice"     // 通知事件
	PostTypeRequest   PostType = "request"    // 请求事件
	PostTypeMetaEvent PostType = "meta_event" // 元事件
)

// MessageType 消息类型
type MessageType string

const (
	MessageTypePrivate MessageType = "private" // 私聊消息
	MessageTypeGroup   MessageType = "group"   // 群消息
)

// PrivateMessageSubType 私聊消息子类型
type PrivateMessageSubType string

const (
	PrivateMessageSubTypeFriend PrivateMessageSubType = "friend" // 好友
	PrivateMessageSubTypeGroup  PrivateMessageSubType = "group"  // 群临时会话
	PrivateMessageSubTypeOther  PrivateMessageSubType = "other"  // 其他
)

// GroupMessageSubType 群消息子类型
type GroupMessageSubType string

const (
	GroupMessageSubTypeNormal    GroupMessageSubType = "normal"    // 正常消息
	GroupMessageSubTypeAnonymous GroupMessageSubType = "anonymous" // 匿名消息
	GroupMessageSubTypeNotice    GroupMessageSubType = "notice"    // 系统提示
)

// NoticeType 通知类型
type NoticeType string

const (
	NoticeTypeGroupUpload   NoticeType = "group_upload"   // 群文件上传
	NoticeTypeGroupAdmin    NoticeType = "group_admin"    // 群管理员变动
	NoticeTypeGroupDecrease NoticeType = "group_decrease" // 群成员减少
	NoticeTypeGroupIncrease NoticeType = "group_increase" // 群成员增加
	NoticeTypeGroupBan      NoticeType = "group_ban"      // 群禁言
	NoticeTypeFriendAdd     NoticeType = "friend_add"     // 好友添加
	NoticeTypeGroupRecall   NoticeType = "group_recall"   // 群消息撤回
	NoticeTypeFriendRecall  NoticeType = "friend_recall"  // 好友消息撤回
	NoticeTypeNotify        NoticeType = "notify"         // 群内提示
	NoticeTypeGroupCard     NoticeType = "group_card"     // 群名片变更
	NoticeTypeOfflineFile   NoticeType = "offline_file"   // 接收离线文件
)

// NotifySubType 提示类型子类型
type NotifySubType string

const (
	NotifySubTypePoke  NotifySubType = "poke"       // 戳一戳
	NotifySubTypeHonor NotifySubType = "honor"      // 群荣誉变更
	NotifySubTypeLucky NotifySubType = "lucky_king" // 群红包运气王
	NotifySubTypeTitle NotifySubType = "title"      // 群成员头衔变更
)

// RequestType 请求类型
type RequestType string

const (
	RequestTypeFriend RequestType = "friend" // 加好友请求
	RequestTypeGroup  RequestType = "group"  // 加群请求/邀请
)

// GroupRequestSubType 群请求子类型
type GroupRequestSubType string

const (
	GroupRequestSubTypeAdd    GroupRequestSubType = "add"    // 加群请求
	GroupRequestSubTypeInvite GroupRequestSubType = "invite" // 邀请登录号入群
)

// MetaEventType 元事件类型
type MetaEventType string

const (
	MetaEventTypeLifecycle MetaEventType = "lifecycle" // 生命周期
	MetaEventTypeHeartbeat MetaEventType = "heartbeat" // 心跳
)

// LifecycleSubType 生命周期子类型
type LifecycleSubType string

const (
	LifecycleSubTypeEnable  LifecycleSubType = "enable"  // OneBot 启用
	LifecycleSubTypeDisable LifecycleSubType = "disable" // OneBot 停用
	LifecycleSubTypeConnect LifecycleSubType = "connect" // WebSocket 连接成功
)

// Role 角色
type Role string

const (
	RoleOwner  Role = "owner"  // 群主
	RoleAdmin  Role = "admin"  // 管理员
	RoleMember Role = "member" // 普通成员
)

// Sex 性别
type Sex string

const (
	SexMale    Sex = "male"    // 男
	SexFemale  Sex = "female"  // 女
	SexUnknown Sex = "unknown" // 未知
)

// HonorType 荣誉类型
type HonorType string

const (
	HonorTypeTalkative HonorType = "talkative" // 龙王
	HonorTypePerformer HonorType = "performer" // 群聊之火
	HonorTypeEmotion   HonorType = "emotion"   // 快乐源泉
)
