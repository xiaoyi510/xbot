package api

import (
	"encoding/json"
	"fmt"
	"github.com/xiaoyi510/xbot/driver"
	"github.com/xiaoyi510/xbot/types"
)

// Client API 客户端
type Client struct {
	driver driver.Driver
}

// NewClient 创建 API 客户端
func NewClient(d driver.Driver) *Client {
	return &Client{
		driver: d,
	}
}

// CallAPI 调用 API
func (c *Client) CallAPI(action string, params map[string]interface{}) (*types.APIResponse, error) {
	return c.driver.CallAPI(action, params)
}

// SendPrivateMsg 发送私聊消息
func (c *Client) SendPrivateMsg(userID int64, message interface{}) (*types.Response[types.MessageResponse], error) {
	params := map[string]interface{}{
		"user_id": userID,
		"message": message,
	}

	resp, err := c.CallAPI(ActionSendPrivateMsg, params)
	if err != nil {
		return nil, err
	}

	var result types.Response[types.MessageResponse]
	result.Status = resp.Status
	result.RetCode = resp.RetCode
	result.Message = resp.Message
	result.Wording = resp.Wording

	if resp.Data != nil {
		if data, ok := resp.Data.(map[string]interface{}); ok {
			if msgID, ok := data["message_id"].(float64); ok {
				result.Data.MessageID = int64(msgID)
			}
		}
	}

	return &result, nil
}

// SendGroupMsg 发送群消息
func (c *Client) SendGroupMsg(groupID int64, message interface{}) (*types.Response[types.MessageResponse], error) {
	params := map[string]interface{}{
		"group_id": groupID,
		"message":  message,
	}

	resp, err := c.CallAPI(ActionSendGroupMsg, params)
	if err != nil {
		return nil, err
	}

	var result types.Response[types.MessageResponse]
	result.Status = resp.Status
	result.RetCode = resp.RetCode
	result.Message = resp.Message
	result.Wording = resp.Wording

	if resp.Data != nil {
		if data, ok := resp.Data.(map[string]interface{}); ok {
			if msgID, ok := data["message_id"].(float64); ok {
				result.Data.MessageID = int64(msgID)
			}
		}
	}

	return &result, nil
}

// DeleteMsg 撤回消息
func (c *Client) DeleteMsg(messageID int64) error {
	params := map[string]interface{}{
		"message_id": messageID,
	}

	_, err := c.CallAPI(ActionDeleteMsg, params)
	return err
}

// GetMsg 获取消息
func (c *Client) GetMsg(messageID int64) (*types.Response[types.MessageData], error) {
	params := map[string]interface{}{
		"message_id": messageID,
	}

	resp, err := c.CallAPI(ActionGetMsg, params)
	if err != nil {
		return nil, err
	}

	var result types.Response[types.MessageData]
	result.Status = resp.Status
	result.RetCode = resp.RetCode

	// TODO: 解析 Data

	return &result, nil
}

// SetGroupKick 群组踢人
func (c *Client) SetGroupKick(groupID, userID int64, rejectAddRequest bool) error {
	params := map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
	}

	if rejectAddRequest {
		params["reject_add_request"] = true
	}

	_, err := c.CallAPI(ActionSetGroupKick, params)
	return err
}

// SetGroupBan 群组禁言
func (c *Client) SetGroupBan(groupID, userID int64, duration int32) error {
	params := map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"duration": duration,
	}

	_, err := c.CallAPI(ActionSetGroupBan, params)
	return err
}

// SetGroupWholeBan 群组全员禁言
func (c *Client) SetGroupWholeBan(groupID int64, enable bool) error {
	params := map[string]interface{}{
		"group_id": groupID,
		"enable":   enable,
	}

	_, err := c.CallAPI(ActionSetGroupWholeBan, params)
	return err
}

// SetGroupAdmin 设置群管理员
func (c *Client) SetGroupAdmin(groupID, userID int64, enable bool) error {
	params := map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"enable":   enable,
	}

	_, err := c.CallAPI(ActionSetGroupAdmin, params)
	return err
}

// SetGroupCard 设置群名片
func (c *Client) SetGroupCard(groupID, userID int64, card string) error {
	params := map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"card":     card,
	}

	_, err := c.CallAPI(ActionSetGroupCard, params)
	return err
}

// GetLoginInfo 获取登录号信息
func (c *Client) GetLoginInfo() (*types.Response[types.LoginInfo], error) {
	resp, err := c.CallAPI(ActionGetLoginInfo, nil)
	if err != nil {
		return nil, err
	}

	var result types.Response[types.LoginInfo]
	result.Status = resp.Status
	result.RetCode = resp.RetCode

	// TODO: 解析 Data

	return &result, nil
}

// GetGroupList 获取群列表
func (c *Client) GetGroupList() (*types.Response[[]types.GroupInfo], error) {
	resp, err := c.CallAPI(ActionGetGroupList, nil)
	if err != nil {
		return nil, err
	}

	var result types.Response[[]types.GroupInfo]
	result.Status = resp.Status
	result.RetCode = resp.RetCode

	// TODO: 解析 Data

	return &result, nil
}

// SetGroupAddRequest 处理加群请求
func (c *Client) SetGroupAddRequest(flag, subType string, approve bool, reason string) error {
	params := map[string]interface{}{
		"flag":     flag,
		"sub_type": subType,
		"approve":  approve,
		"reason":   reason,
	}

	_, err := c.CallAPI(ActionSetGroupAddRequest, params)
	return err
}

// SendGroupNotice 发送群公告
func (c *Client) SendGroupNotice(groupID int64, content string) error {
	params := map[string]interface{}{
		"group_id": groupID,
		"content":  content,
	}

	_, err := c.CallAPI(ActionSendGroupNotice, params)
	return err
}

// SetEssenceMsg 设置精华消息
func (c *Client) SetEssenceMsg(messageID int64) error {
	params := map[string]interface{}{
		"message_id": messageID,
	}

	_, err := c.CallAPI(ActionSetEssenceMsg, params)
	return err
}

// DeleteEssenceMsg 移出精华消息
func (c *Client) DeleteEssenceMsg(messageID int64) error {
	params := map[string]interface{}{
		"message_id": messageID,
	}

	_, err := c.CallAPI(ActionDeleteEssenceMsg, params)
	return err
}

// GetEssenceMsgList 获取精华消息列表
func (c *Client) GetEssenceMsgList(groupID int64) ([]types.EssenceMessage, error) {
	params := map[string]interface{}{
		"group_id": groupID,
	}

	resp, err := c.CallAPI(ActionGetEssenceMsgList, params)
	if err != nil {
		return nil, err
	}

	// 解析为精华消息列表
	if resp.Data == nil {
		return []types.EssenceMessage{}, nil
	}

	// 尝试解析为数组
	dataArray, ok := resp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid essence message list data format")
	}

	essenceMessages := make([]types.EssenceMessage, 0, len(dataArray))
	for _, item := range dataArray {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		essence := types.EssenceMessage{
			SenderID:     getInt64(itemMap, "sender_id"),
			SenderNick:   getString(itemMap, "sender_nick"),
			SenderTime:   getInt64(itemMap, "sender_time"),
			OperatorID:   getInt64(itemMap, "operator_id"),
			OperatorNick: getString(itemMap, "operator_nick"),
			OperatorTime: getInt64(itemMap, "operator_time"),
			MessageID:    getInt64(itemMap, "message_id"),
			MessageSeq:   getInt64(itemMap, "msg_seq"),
			MessageRand:  getInt64(itemMap, "msg_random"),
			Content:      itemMap["content"], // 保持原始类型，由使用者解析
		}
		essenceMessages = append(essenceMessages, essence)
	}

	return essenceMessages, nil
}

// ParseEssenceContent 解析精华消息内容为消息段数组
// 这是一个辅助函数，用于将 Content 转换为可用的消息段
func ParseEssenceContent(content interface{}) ([]map[string]interface{}, error) {
	if content == nil {
		return []map[string]interface{}{}, nil
	}

	// 尝试解析为数组
	contentArray, ok := content.([]interface{})
	if !ok {
		return nil, fmt.Errorf("content is not an array")
	}

	segments := make([]map[string]interface{}, 0, len(contentArray))
	for _, item := range contentArray {
		if segMap, ok := item.(map[string]interface{}); ok {
			segments = append(segments, segMap)
		}
	}

	return segments, nil
}

// GetEssenceContentText 从精华消息内容中提取纯文本
// 这是一个便捷函数，用于快速获取消息的文本内容
func GetEssenceContentText(content interface{}) string {
	segments, err := ParseEssenceContent(content)
	if err != nil {
		return ""
	}

	var text string
	for _, seg := range segments {
		if segType, ok := seg["type"].(string); ok && segType == "text" {
			if data, ok := seg["data"].(map[string]interface{}); ok {
				if textVal, ok := data["text"].(string); ok {
					text += textVal
				}
			}
		}
	}

	return text
}

// 辅助函数：从 map 中获取 int64
func getInt64(m map[string]interface{}, key string) int64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		case json.Number:
			if i, err := v.Int64(); err == nil {
				return i
			}
		}
	}
	return 0
}

// 辅助函数：从 map 中获取 string
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// SendGroupForwardMsg 发送群合并转发消息
// messages: 消息节点数组，每个节点使用 message.CustomNode() 或 message.Node() 创建
func (c *Client) SendGroupForwardMsg(groupID int64, messages []interface{}) (*types.Response[types.ForwardMessageResponse], error) {
	params := map[string]interface{}{
		"group_id": groupID,
		"messages": messages,
	}

	resp, err := c.CallAPI(ActionSendGroupForwardMsg, params)
	if err != nil {
		return nil, err
	}

	var result types.Response[types.ForwardMessageResponse]
	result.Status = resp.Status
	result.RetCode = resp.RetCode
	result.Message = resp.Message
	result.Wording = resp.Wording

	if resp.Data != nil {
		if data, ok := resp.Data.(map[string]interface{}); ok {
			if msgID, ok := data["message_id"].(float64); ok {
				result.Data.MessageID = int64(msgID)
			}
			if fwdID, ok := data["forward_id"].(string); ok {
				result.Data.ForwardID = fwdID
			}
		}
	}

	return &result, nil
}

// SendPrivateForwardMsg 发送私聊合并转发消息
// messages: 消息节点数组，每个节点使用 message.CustomNode() 或 message.Node() 创建
func (c *Client) SendPrivateForwardMsg(userID int64, messages []interface{}) (*types.Response[types.ForwardMessageResponse], error) {
	params := map[string]interface{}{
		"user_id":  userID,
		"messages": messages,
	}

	resp, err := c.CallAPI(ActionSendPrivateForwardMsg, params)
	if err != nil {
		return nil, err
	}

	var result types.Response[types.ForwardMessageResponse]
	result.Status = resp.Status
	result.RetCode = resp.RetCode
	result.Message = resp.Message
	result.Wording = resp.Wording

	if resp.Data != nil {
		if data, ok := resp.Data.(map[string]interface{}); ok {
			if msgID, ok := data["message_id"].(float64); ok {
				result.Data.MessageID = int64(msgID)
			}
			if fwdID, ok := data["forward_id"].(string); ok {
				result.Data.ForwardID = fwdID
			}
		}
	}

	return &result, nil
}
