package message

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	cqCodeRegex = regexp.MustCompile(`\[CQ:([a-zA-Z]+)((?:,[a-zA-Z0-9_\-]+=[^,\]]*)*)\]`)
)

// ParseCQCode 解析 CQ 码为消息段数组
func ParseCQCode(cqCode string) Message {
	var message Message

	lastIndex := 0
	matches := cqCodeRegex.FindAllStringSubmatchIndex(cqCode, -1)

	for _, match := range matches {
		// 提取 CQ 码之前的纯文本
		if match[0] > lastIndex {
			text := cqCode[lastIndex:match[0]]
			if text != "" {
				message = append(message, Text(unescapeCQCode(text)))
			}
		}

		// 提取 CQ 码类型
		cqType := cqCode[match[2]:match[3]]

		// 提取 CQ 码参数
		paramsStr := ""
		if match[5] > match[4] {
			paramsStr = cqCode[match[4]:match[5]]
		}

		// 解析参数
		data := parseCQParams(paramsStr)

		message = append(message, MessageSegment{
			Type: cqType,
			Data: data,
		})

		lastIndex = match[1]
	}

	// 提取最后的纯文本
	if lastIndex < len(cqCode) {
		text := cqCode[lastIndex:]
		if text != "" {
			message = append(message, Text(unescapeCQCode(text)))
		}
	}

	// 如果没有匹配到任何 CQ 码，返回纯文本
	if len(message) == 0 && cqCode != "" {
		message = append(message, Text(cqCode))
	}

	return message
}

// parseCQParams 解析 CQ 码参数
func parseCQParams(paramsStr string) map[string]interface{} {
	data := make(map[string]interface{})

	if paramsStr == "" {
		return data
	}

	// 移除开头的逗号
	paramsStr = strings.TrimPrefix(paramsStr, ",")

	// 分割参数
	params := strings.Split(paramsStr, ",")
	for _, param := range params {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := unescapeCQCode(kv[1])

			// 尝试转换为数字
			if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
				data[key] = intVal
			} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
				data[key] = floatVal
			} else {
				data[key] = value
			}
		}
	}

	return data
}

// ToCQCode 将消息段数组转换为 CQ 码字符串
func (m Message) ToCQCode() string {
	var sb strings.Builder

	for _, seg := range m {
		if seg.Type == "text" {
			if text, ok := seg.Data["text"].(string); ok {
				sb.WriteString(escapeCQCode(text))
			}
		} else {
			sb.WriteString("[CQ:")
			sb.WriteString(seg.Type)
			for key, value := range seg.Data {
				sb.WriteString(",")
				sb.WriteString(key)
				sb.WriteString("=")
				sb.WriteString(escapeCQCode(fmt.Sprint(value)))
			}
			sb.WriteString("]")
		}
	}

	return sb.String()
}

// escapeCQCode 转义 CQ 码特殊字符
func escapeCQCode(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "[", "&#91;")
	s = strings.ReplaceAll(s, "]", "&#93;")
	s = strings.ReplaceAll(s, ",", "&#44;")
	return s
}

// unescapeCQCode 反转义 CQ 码特殊字符
func unescapeCQCode(s string) string {
	s = strings.ReplaceAll(s, "&#44;", ",")
	s = strings.ReplaceAll(s, "&#93;", "]")
	s = strings.ReplaceAll(s, "&#91;", "[")
	s = strings.ReplaceAll(s, "&amp;", "&")
	return s
}
