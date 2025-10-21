package message

import "fmt"

// Builder 消息构建器
type Builder struct {
	message Message
}

// NewBuilder 创建消息构建器
func NewBuilder() *Builder {
	return &Builder{
		message: make(Message, 0),
	}
}

// Text 添加纯文本
func (b *Builder) Text(text string) *Builder {
	b.message = append(b.message, Text(text))
	return b
}

// Textf 添加格式化文本
func (b *Builder) Textf(format string, args ...interface{}) *Builder {
	return b.Text(fmt.Sprintf(format, args...))
}

// Face 添加 QQ 表情
func (b *Builder) Face(id int) *Builder {
	b.message = append(b.message, Face(id))
	return b
}

// Image 添加图片
func (b *Builder) Image(file string) *Builder {
	b.message = append(b.message, Image(file))
	return b
}

// Record 添加语音
func (b *Builder) Record(file string) *Builder {
	b.message = append(b.message, Record(file))
	return b
}

// Video 添加视频
func (b *Builder) Video(file string) *Builder {
	b.message = append(b.message, Video(file))
	return b
}

// At 添加 @某人
func (b *Builder) At(qq int64) *Builder {
	b.message = append(b.message, At(qq))
	return b
}

// AtAll 添加 @全体成员
func (b *Builder) AtAll() *Builder {
	b.message = append(b.message, AtAll())
	return b
}

// RPS 添加猜拳
func (b *Builder) RPS() *Builder {
	b.message = append(b.message, RPS())
	return b
}

// Dice 添加骰子
func (b *Builder) Dice() *Builder {
	b.message = append(b.message, Dice())
	return b
}

// Shake 添加抖动
func (b *Builder) Shake() *Builder {
	b.message = append(b.message, Shake())
	return b
}

// Poke 添加戳一戳
func (b *Builder) Poke(pokeType, id int) *Builder {
	b.message = append(b.message, Poke(pokeType, id))
	return b
}

// Share 添加链接分享
func (b *Builder) Share(url, title string) *Builder {
	b.message = append(b.message, Share(url, title))
	return b
}

// ShareWithImage 添加链接分享（带图片）
func (b *Builder) ShareWithImage(url, title, content, image string) *Builder {
	b.message = append(b.message, ShareWithImage(url, title, content, image))
	return b
}

// ContactFriend 添加推荐好友
func (b *Builder) ContactFriend(qq int64) *Builder {
	b.message = append(b.message, ContactFriend(qq))
	return b
}

// ContactGroup 添加推荐群
func (b *Builder) ContactGroup(groupID int64) *Builder {
	b.message = append(b.message, ContactGroup(groupID))
	return b
}

// Location 添加位置
func (b *Builder) Location(lat, lon float64, title, content string) *Builder {
	b.message = append(b.message, Location(lat, lon, title, content))
	return b
}

// Music 添加音乐分享
func (b *Builder) Music(musicType string, id int64) *Builder {
	b.message = append(b.message, Music(musicType, id))
	return b
}

// CustomMusic 添加自定义音乐
func (b *Builder) CustomMusic(url, audio, title string) *Builder {
	b.message = append(b.message, CustomMusic(url, audio, title))
	return b
}

// Reply 添加回复
func (b *Builder) Reply(messageID int64) *Builder {
	b.message = append(b.message, Reply(messageID))
	return b
}

// XML 添加 XML 消息
func (b *Builder) XML(data string) *Builder {
	b.message = append(b.message, XML(data))
	return b
}

// JSON 添加 JSON 消息
func (b *Builder) JSON(data string) *Builder {
	b.message = append(b.message, JSON(data))
	return b
}

// Segment 添加自定义消息段
func (b *Builder) Segment(seg MessageSegment) *Builder {
	b.message = append(b.message, seg)
	return b
}

// Segments 添加多个消息段
func (b *Builder) Segments(segs ...MessageSegment) *Builder {
	b.message = append(b.message, segs...)
	return b
}

// Build 构建消息
func (b *Builder) Build() Message {
	return b.message
}

// BuildCQCode 构建为 CQ 码字符串
func (b *Builder) BuildCQCode() string {
	return b.message.ToCQCode()
}

// Clear 清空消息
func (b *Builder) Clear() *Builder {
	b.message = make(Message, 0)
	return b
}

// Len 返回消息段数量
func (b *Builder) Len() int {
	return len(b.message)
}
