package schema

// DouYinIMGroupRecallReq 撤回群消息请求参数
type DouYinIMGroupRecallReq struct {
	MsgId            string `json:"msg_id"`
	ConversationId   string `json:"conversation_id"`
	ConversationType int    `json:"conversation_type"`
}

// DouYinIMGroupRecallRes 撤回群消息响应参数
type DouYinIMGroupRecallRes struct {
	ErrMsg string `json:"err_msg"` // 错误信息
	ErrNo  int    `json:"err_no"`  // 错误码，0-成功，非0-失败
	LogId  string `json:"log_id"`  // 日志ID，用于问题定位
}
