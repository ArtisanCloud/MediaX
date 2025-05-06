package schema

// DouYinTaskBoxChangeUserBindAgentReq 切换用户绑定团长请求参数
type DouYinTaskBoxChangeUserBindAgentReq struct {
	// DouyinID 达人的抖音号
	DouyinID string `json:"douyin_id"`
	// NewAgentID 达人换绑的团长ID
	NewAgentID int64 `json:"new_agent_id"`
	// OldAgentID 达人当前绑定的团长ID
	OldAgentID int64 `json:"old_agent_id"`
}

type DouYinTaskBoxChangeUserBindAgentRes struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	LogId  string `json:"log_id"`
}
