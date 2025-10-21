package message

import (
	"fmt"
	"strings"
)

// GetRawMessage 获取消息的原始内容
func (m Message) GetRawMessage() string {
	var rawMessage strings.Builder
	for _, seg := range m {
		rawMessage.WriteString(seg.Type)
		rawMessage.WriteString(fmt.Sprintf("%v", seg.Data))
	}
	return rawMessage.String()
}

// GetPlainText 获取消息中的纯文本内容
func (m Message) GetPlainText() string {
	var texts []string
	for _, seg := range m {
		if seg.Type == "text" {
			if text, ok := seg.Data["text"].(string); ok {
				texts = append(texts, text)
			}
		}
	}
	return strings.Join(texts, "")
}

// HasType 判断消息中是否包含指定类型的消息段
func (m Message) HasType(segType string) bool {
	for _, seg := range m {
		if seg.Type == segType {
			return true
		}
	}
	return false
}

// GetSegmentsByType 获取指定类型的所有消息段
func (m Message) GetSegmentsByType(segType string) []MessageSegment {
	var segments []MessageSegment
	for _, seg := range m {
		if seg.Type == segType {
			segments = append(segments, seg)
		}
	}
	return segments
}

// GetFirstSegmentByType 获取指定类型的第一个消息段
func (m Message) GetFirstSegmentByType(segType string) (MessageSegment, bool) {
	for _, seg := range m {
		if seg.Type == segType {
			return seg, true
		}
	}
	return MessageSegment{}, false
}

// Append 追加消息段
func (m *Message) Append(segments ...MessageSegment) {
	*m = append(*m, segments...)
}

// Prepend 前置消息段
func (m *Message) Prepend(segments ...MessageSegment) {
	*m = append(segments, *m...)
}

// Len 返回消息段数量
func (m Message) Len() int {
	return len(m)
}

// IsEmpty 判断消息是否为空
func (m Message) IsEmpty() bool {
	return len(m) == 0
}

// ParseMessage 解析消息（从 interface{} 或 string）
func ParseMessage(msg interface{}) Message {
	switch v := msg.(type) {
	case string:
		// 字符串消息，尝试解析 CQ 码或作为纯文本
		return ParseCQCode(v)
	case []interface{}:
		// 消息段数组
		var message Message
		for _, item := range v {
			if segMap, ok := item.(map[string]interface{}); ok {
				seg := MessageSegment{
					Type: getString(segMap, "type"),
					Data: getMap(segMap, "data"),
				}
				message = append(message, seg)
			}
		}
		return message
	case []MessageSegment:
		return Message(v)
	case Message:
		return v
	default:
		return Message{}
	}
}

// getString 从 map 中获取字符串
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getMap 从 map 中获取 map
func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
	}
	return make(map[string]interface{})
}
