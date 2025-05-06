package schema

// DouYinTaskBoxGenAgentLinkReq 生成代理链接请求参数
type DouYinTaskBoxGenAgentLinkReq struct {
	// AgentID 团长ID
	AgentID int64 `json:"agent_id"`
	// AgencyTalentUID 开发者透传的达人uid
	AgencyTalentUID string `json:"agency_talent_uid,omitempty"`
	// AppID 小程序appid
	AppID string `json:"app_id,omitempty"`
	// TaskCategory 任务类别枚举
	TaskCategory int `json:"task_category,omitempty"`
	// TaskID 任务ID
	TaskID int64 `json:"task_id,omitempty"`
}

// GenAgentLink 生成代理链接响应结果
type GenAgentLink struct {
	// WebLink 网页链接
	WebLink string `json:"web_link"`
	// AppLink 小程序链接
	AppLink string `json:"app_link"`
}

// DouYinTaskBoxGenAgentLinkRes 生成代理链接响应结果
type DouYinTaskBoxGenAgentLinkRes struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	LogId  string `json:"log_id"`
	// Data 业务数据主体，包含具体的业务响应信息
	Data GenAgentLink `json:"data"`
}
