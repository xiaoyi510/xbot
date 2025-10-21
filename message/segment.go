package message

// MessageSegment 消息段
type MessageSegment struct {
	Type string                 `json:"type"` // 消息段类型
	Data map[string]interface{} `json:"data"` // 消息段数据
}

// Message 消息（消息段数组）
type Message []MessageSegment

// Text 纯文本消息段
func Text(text string) MessageSegment {
	return MessageSegment{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	}
}

// Face QQ 表情消息段
func Face(id int) MessageSegment {
	return MessageSegment{
		Type: "face",
		Data: map[string]interface{}{
			"id": id,
		},
	}
}

// Image 图片消息段
// file: 图片文件名、绝对路径、URL、base64://
func Image(file string) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

// ImageWithOptions 图片消息段（带选项）
func ImageWithOptions(file string, imageType string, cache bool, proxy bool, timeout int) MessageSegment {
	data := map[string]interface{}{
		"file": file,
	}
	if imageType != "" {
		data["type"] = imageType
	}
	if !cache {
		data["cache"] = 0
	}
	if !proxy {
		data["proxy"] = 0
	}
	if timeout > 0 {
		data["timeout"] = timeout
	}
	return MessageSegment{
		Type: "image",
		Data: data,
	}
}

// Record 语音消息段
func Record(file string) MessageSegment {
	return MessageSegment{
		Type: "record",
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

// RecordWithMagic 语音消息段（带变声）
func RecordWithMagic(file string, magic bool) MessageSegment {
	data := map[string]interface{}{
		"file": file,
	}
	if magic {
		data["magic"] = 1
	}
	return MessageSegment{
		Type: "record",
		Data: data,
	}
}

// Video 短视频消息段
func Video(file string) MessageSegment {
	return MessageSegment{
		Type: "video",
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

// At @某人消息段
func At(qq int64) MessageSegment {
	return MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": qq,
		},
	}
}

// AtAll @全体成员消息段
func AtAll() MessageSegment {
	return MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": "all",
		},
	}
}

// RPS 猜拳魔法表情消息段
func RPS() MessageSegment {
	return MessageSegment{
		Type: "rps",
		Data: map[string]interface{}{},
	}
}

// Dice 掷骰子魔法表情消息段
func Dice() MessageSegment {
	return MessageSegment{
		Type: "dice",
		Data: map[string]interface{}{},
	}
}

// Shake 窗口抖动消息段
func Shake() MessageSegment {
	return MessageSegment{
		Type: "shake",
		Data: map[string]interface{}{},
	}
}

// Poke 戳一戳消息段
func Poke(pokeType int, id int) MessageSegment {
	return MessageSegment{
		Type: "poke",
		Data: map[string]interface{}{
			"type": pokeType,
			"id":   id,
		},
	}
}

// Anonymous 匿名发消息段（配合群消息使用）
func Anonymous(ignore bool) MessageSegment {
	data := map[string]interface{}{}
	if ignore {
		data["ignore"] = 1
	}
	return MessageSegment{
		Type: "anonymous",
		Data: data,
	}
}

// Share 链接分享消息段
func Share(url, title string) MessageSegment {
	return MessageSegment{
		Type: "share",
		Data: map[string]interface{}{
			"url":   url,
			"title": title,
		},
	}
}

// ShareWithImage 链接分享消息段（带图片和描述）
func ShareWithImage(url, title, content, image string) MessageSegment {
	return MessageSegment{
		Type: "share",
		Data: map[string]interface{}{
			"url":     url,
			"title":   title,
			"content": content,
			"image":   image,
		},
	}
}

// Contact 推荐好友消息段
func Contact(contactType string, id int64) MessageSegment {
	return MessageSegment{
		Type: "contact",
		Data: map[string]interface{}{
			"type": contactType, // qq 或 group
			"id":   id,
		},
	}
}

// ContactFriend 推荐好友
func ContactFriend(qq int64) MessageSegment {
	return Contact("qq", qq)
}

// ContactGroup 推荐群
func ContactGroup(groupID int64) MessageSegment {
	return Contact("group", groupID)
}

// Location 位置消息段
func Location(lat, lon float64, title, content string) MessageSegment {
	return MessageSegment{
		Type: "location",
		Data: map[string]interface{}{
			"lat":     lat,
			"lon":     lon,
			"title":   title,
			"content": content,
		},
	}
}

// Music 音乐分享消息段
// musicType: qq, 163, xm
func Music(musicType string, id int64) MessageSegment {
	return MessageSegment{
		Type: "music",
		Data: map[string]interface{}{
			"type": musicType,
			"id":   id,
		},
	}
}

// CustomMusic 音乐自定义分享消息段
func CustomMusic(url, audio, title string) MessageSegment {
	return MessageSegment{
		Type: "music",
		Data: map[string]interface{}{
			"type":  "custom",
			"url":   url,
			"audio": audio,
			"title": title,
		},
	}
}

// CustomMusicWithImage 音乐自定义分享消息段（带图片和副标题）
func CustomMusicWithImage(url, audio, title, content, image string) MessageSegment {
	return MessageSegment{
		Type: "music",
		Data: map[string]interface{}{
			"type":    "custom",
			"url":     url,
			"audio":   audio,
			"title":   title,
			"content": content,
			"image":   image,
		},
	}
}

// Reply 回复消息段
func Reply(messageID int64) MessageSegment {
	return MessageSegment{
		Type: "reply",
		Data: map[string]interface{}{
			"id": messageID,
		},
	}
}

// Forward 合并转发消息段
func Forward(id string) MessageSegment {
	return MessageSegment{
		Type: "forward",
		Data: map[string]interface{}{
			"id": id,
		},
	}
}

// Node 合并转发节点消息段
func Node(messageID int32) MessageSegment {
	return MessageSegment{
		Type: "node",
		Data: map[string]interface{}{
			"id": messageID,
		},
	}
}

// CustomNode 合并转发自定义节点消息段
func CustomNode(userID int64, nickname string, content interface{}) MessageSegment {
	return MessageSegment{
		Type: "node",
		Data: map[string]interface{}{
			"user_id":  userID,
			"nickname": nickname,
			"content":  content,
		},
	}
}

// XML XML 消息段
func XML(data string) MessageSegment {
	return MessageSegment{
		Type: "xml",
		Data: map[string]interface{}{
			"data": data,
		},
	}
}

// JSON JSON 消息段
func JSON(data string) MessageSegment {
	return MessageSegment{
		Type: "json",
		Data: map[string]interface{}{
			"data": data,
		},
	}
}

// CardImage 卡片图片消息段（某些 XML 的图片）
func CardImage(file string) MessageSegment {
	return MessageSegment{
		Type: "cardimage",
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

// TTS 文本转语音消息段
func TTS(text string) MessageSegment {
	return MessageSegment{
		Type: "tts",
		Data: map[string]interface{}{
			"text": text,
		},
	}
}
