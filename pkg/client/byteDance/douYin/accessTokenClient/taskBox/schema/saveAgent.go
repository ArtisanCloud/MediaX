package schema

// DouYinTaskBoxSaveAgentReq 保存团长信息请求参数
type DouYinTaskBoxSaveAgentReq struct {
	// AgentID 团长ID
	// 创建团长时，agent_id不传或传0
	// 修改团长信息时，传入待修改团长的agent_id
	AgentID int64 `json:"agent_id,omitempty"`
	// AgentNickname 团长的昵称，全局唯一
	AgentNickname string `json:"agent_nickname"`
}

// SaveAgent 保存团长信息响应结果
type SaveAgent struct {
	AgentId int64 `json:"agent_id"`
}

// DouYinTaskBoxSaveAgentRes 保存团长信息响应结果
type DouYinTaskBoxSaveAgentRes struct {
	ErrNo  int       `json:"err_no"`
	ErrMsg string    `json:"err_msg"`
	LogId  string    `json:"log_id"`
	Data   SaveAgent `json:"data"`
}
