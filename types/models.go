package types

// Sender 发送者信息
type Sender struct {
	UserID   int64  `json:"user_id"`         // 发送者 QQ 号
	Nickname string `json:"nickname"`        // 昵称
	Card     string `json:"card,omitempty"`  // 群名片/备注
	Sex      Sex    `json:"sex,omitempty"`   // 性别
	Age      int32  `json:"age,omitempty"`   // 年龄
	Area     string `json:"area,omitempty"`  // 地区
	Level    string `json:"level,omitempty"` // 等级
	Role     Role   `json:"role,omitempty"`  // 角色
	Title    string `json:"title,omitempty"` // 专属头衔
}

// Anonymous 匿名信息
type Anonymous struct {
	ID   int64  `json:"id"`   // 匿名用户 ID
	Name string `json:"name"` // 匿名用户名称
	Flag string `json:"flag"` // 匿名用户 flag，在调用禁言 API 时需要传入
}

// File 文件信息
type File struct {
	ID    string `json:"id"`    // 文件 ID
	Name  string `json:"name"`  // 文件名
	Size  int64  `json:"size"`  // 文件大小
	BusID int64  `json:"busid"` // busid
}

// Status OneBot 状态
type Status struct {
	Online bool `json:"online"` // 当前 QQ 在线
	Good   bool `json:"good"`   // 状态符合预期
}

// VersionInfo 版本信息
type VersionInfo struct {
	AppName         string `json:"app_name"`         // 应用标识
	AppVersion      string `json:"app_version"`      // 应用版本
	ProtocolVersion string `json:"protocol_version"` // OneBot 标准版本
}

// GroupInfo 群信息
type GroupInfo struct {
	GroupID        int64  `json:"group_id"`         // 群号
	GroupName      string `json:"group_name"`       // 群名称
	MemberCount    int32  `json:"member_count"`     // 成员数
	MaxMemberCount int32  `json:"max_member_count"` // 最大成员数
}

// GroupMemberInfo 群成员信息
type GroupMemberInfo struct {
	GroupID         int64  `json:"group_id"`          // 群号
	UserID          int64  `json:"user_id"`           // QQ 号
	Nickname        string `json:"nickname"`          // 昵称
	Card            string `json:"card"`              // 群名片/备注
	Sex             Sex    `json:"sex"`               // 性别
	Age             int32  `json:"age"`               // 年龄
	Area            string `json:"area"`              // 地区
	JoinTime        int64  `json:"join_time"`         // 加群时间戳
	LastSentTime    int64  `json:"last_sent_time"`    // 最后发言时间戳
	Level           string `json:"level"`             // 成员等级
	Role            Role   `json:"role"`              // 角色
	Unfriendly      bool   `json:"unfriendly"`        // 是否不良记录成员
	Title           string `json:"title"`             // 专属头衔
	TitleExpireTime int64  `json:"title_expire_time"` // 专属头衔过期时间戳
	CardChangeable  bool   `json:"card_changeable"`   // 是否允许修改群名片
}

// UserInfo 用户信息
type UserInfo struct {
	UserID   int64  `json:"user_id"`  // QQ 号
	Nickname string `json:"nickname"` // 昵称
	Sex      Sex    `json:"sex"`      // 性别
	Age      int32  `json:"age"`      // 年龄
}

// FriendInfo 好友信息
type FriendInfo struct {
	UserID   int64  `json:"user_id"`  // QQ 号
	Nickname string `json:"nickname"` // 昵称
	Remark   string `json:"remark"`   // 备注名
}

// GroupHonorInfo 群荣誉信息
type GroupHonorInfo struct {
	UserID      int64  `json:"user_id"`     // QQ 号
	Nickname    string `json:"nickname"`    // 昵称
	Avatar      string `json:"avatar"`      // 头像 URL
	Description string `json:"description"` // 荣誉描述
}

// GroupFileSystemInfo 群文件系统信息
type GroupFileSystemInfo struct {
	FileCount  int32 `json:"file_count"`  // 文件总数
	LimitCount int32 `json:"limit_count"` // 文件上限
	UsedSpace  int64 `json:"used_space"`  // 已使用空间
	TotalSpace int64 `json:"total_space"` // 空间上限
}

// GroupFileInfo 群文件信息
type GroupFileInfo struct {
	GroupID       int64  `json:"group_id"`       // 群号
	FileID        string `json:"file_id"`        // 文件 ID
	FileName      string `json:"file_name"`      // 文件名
	BusID         int64  `json:"busid"`          // 文件类型
	FileSize      int64  `json:"file_size"`      // 文件大小
	UploadTime    int64  `json:"upload_time"`    // 上传时间
	DeadTime      int64  `json:"dead_time"`      // 过期时间
	ModifyTime    int64  `json:"modify_time"`    // 最后修改时间
	DownloadTimes int32  `json:"download_times"` // 下载次数
	Uploader      int64  `json:"uploader"`       // 上传者 ID
	UploaderName  string `json:"uploader_name"`  // 上传者名字
}

// GroupFolderInfo 群文件夹信息
type GroupFolderInfo struct {
	GroupID        int64  `json:"group_id"`         // 群号
	FolderID       string `json:"folder_id"`        // 文件夹 ID
	FolderName     string `json:"folder_name"`      // 文件夹名
	CreateTime     int64  `json:"create_time"`      // 创建时间
	Creator        int64  `json:"creator"`          // 创建者
	CreatorName    string `json:"creator_name"`     // 创建者名字
	TotalFileCount int32  `json:"total_file_count"` // 子文件数量
}

// OfflineFileInfo 离线文件信息
type OfflineFileInfo struct {
	Name string `json:"name"` // 文件名
	Size int64  `json:"size"` // 文件大小
	URL  string `json:"url"`  // 下载链接
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	AppID      int64  `json:"app_id"`      // 客户端 ID
	DeviceName string `json:"device_name"` // 设备名称
	DeviceKind string `json:"device_kind"` // 设备类型
}

// MessageData 消息数据（用于获取消息API）
type MessageData struct {
	MessageID  int64       `json:"message_id"`  // 消息 ID
	RealID     int32       `json:"real_id"`     // 真实 ID
	Sender     Sender      `json:"sender"`      // 发送者
	Time       int64       `json:"time"`        // 发送时间
	Message    interface{} `json:"message"`     // 消息内容
	RawMessage string      `json:"raw_message"` // 原始消息
}

// ForwardMessage 合并转发消息
type ForwardMessage struct {
	Messages []ForwardNode `json:"messages"` // 消息节点
}

// ForwardNode 合并转发节点
type ForwardNode struct {
	Type string                 `json:"type"` // 节点类型 node
	Data map[string]interface{} `json:"data"` // 节点数据
}

// ImageInfo 图片信息
type ImageInfo struct {
	Size     int32  `json:"size"`     // 图片大小
	FileName string `json:"filename"` // 图片文件名
	URL      string `json:"url"`      // 图片链接
}

// RecordInfo 语音信息
type RecordInfo struct {
	File      string `json:"file"`                 // 文件名
	URL       string `json:"url"`                  // 下载链接
	OutFormat string `json:"out_format,omitempty"` // 转换后的格式
}

// LoginInfo 登录账号信息
type LoginInfo struct {
	UserID   int64  `json:"user_id"`  // QQ 号
	Nickname string `json:"nickname"` // QQ 昵称
}

// Credentials 认证信息
type Credentials struct {
	Cookies   string `json:"cookies"`    // Cookies
	CSRFToken int32  `json:"csrf_token"` // CSRF Token
}

// EssenceMessage 精华消息信息
type EssenceMessage struct {
	SenderID     int64       `json:"sender_id"`     // 发送者 QQ 号
	SenderNick   string      `json:"sender_nick"`   // 发送者昵称
	SenderTime   int64       `json:"sender_time"`   // 发送时间戳
	OperatorID   int64       `json:"operator_id"`   // 操作者 QQ 号
	OperatorNick string      `json:"operator_nick"` // 操作者昵称
	OperatorTime int64       `json:"operator_time"` // 操作时间戳
	MessageID    int64       `json:"message_id"`    // 消息 ID
	MessageSeq   int64       `json:"msg_seq"`       // 消息序列号
	MessageRand  int64       `json:"msg_random"`    // 消息随机数
	Content      interface{} `json:"content"`       // 消息内容（可能是消息段数组或其他格式）
}
