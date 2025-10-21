package api

// OneBot API 常量定义
const (
	// 消息API
	ActionSendPrivateMsg = "send_private_msg"
	ActionSendGroupMsg   = "send_group_msg"
	ActionDeleteMsg      = "delete_msg"
	ActionGetMsg         = "get_msg"

	// 群管理API
	ActionSetGroupKick       = "set_group_kick"
	ActionSetGroupBan        = "set_group_ban"
	ActionSetGroupWholeBan   = "set_group_whole_ban"
	ActionSetGroupAdmin      = "set_group_admin"
	ActionSetGroupCard       = "set_group_card"
	ActionSetGroupAddRequest = "set_group_add_request"

	// 信息获取API
	ActionGetLoginInfo = "get_login_info"
	ActionGetGroupList = "get_group_list"

	// 群公告和精华消息API
	ActionSendGroupNotice   = "_send_group_notice"
	ActionSetEssenceMsg     = "set_essence_msg"
	ActionDeleteEssenceMsg  = "delete_essence_msg"
	ActionGetEssenceMsgList = "get_essence_msg_list"

	// 合并转发API
	ActionSendGroupForwardMsg   = "send_group_forward_msg"
	ActionSendPrivateForwardMsg = "send_private_forward_msg"
)
