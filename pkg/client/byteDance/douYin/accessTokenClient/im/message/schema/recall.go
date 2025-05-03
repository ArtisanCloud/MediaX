package schema

// DouYinIMMessageRecallReq 撤回消息请求参数
type DouYinIMMessageRecallReq struct {
	ConversationId   string `json:"conversation_id"`   // 会话ID，来源于私信webhook
	ConversationType int32  `json:"conversation_type"` // 会话类型，1-单聊，2-群聊
	MsgId            string `json:"msg_id"`            // 消息ID，来源于私信webhook
}

// DouYinIMMessageRecallRes 撤回消息响应参数
type DouYinIMMessageRecallRes struct {
	ErrMsg string `json:"err_msg"` // 错误信息
	ErrNo  int    `json:"err_no"`  // 错误码，0-成功，非0-失败
	LogId  string `json:"log_id"`  // 日志ID，用于问题定位
}
