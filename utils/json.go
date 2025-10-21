package utils

import (
	"encoding/json"
)

// ToJSON 将对象转换为 JSON 字符串
func ToJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 从 JSON 字符串解析对象
func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// ToJSONBytes 将对象转换为 JSON 字节数组
func ToJSONBytes(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// FromJSONBytes 从 JSON 字节数组解析对象
func FromJSONBytes(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ToJSONIndent 将对象转换为格式化的 JSON 字符串
func ToJSONIndent(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MustToJSON 将对象转换为 JSON 字符串，失败则 panic
func MustToJSON(v interface{}) string {
	str, err := ToJSON(v)
	if err != nil {
		panic(err)
	}
	return str
}

// MustFromJSON 从 JSON 字符串解析对象，失败则 panic
func MustFromJSON(data string, v interface{}) {
	if err := FromJSON(data, v); err != nil {
		panic(err)
	}
}
