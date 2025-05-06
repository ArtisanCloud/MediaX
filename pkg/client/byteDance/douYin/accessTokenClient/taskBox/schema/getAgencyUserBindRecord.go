package schema

// DouYinTaskBoxGetAgencyUserBindRecordReq 获取机构用户绑定记录请求参数
type DouYinTaskBoxGetAgencyUserBindRecordReq struct {
	// PageNo 页号
	PageNo int32 `json:"page_no"`
	// PageSize 页大小
	PageSize int32 `json:"page_size"`
	// AgencyTalentUID 机构侧达人uid
	AgencyTalentUID string `json:"agency_talent_uid,omitempty"`
	// AgentID 团长ID
	AgentID int64 `json:"agent_id,omitempty"`
	// DouyinID 达人抖音号
	DouyinID string `json:"douyin_id,omitempty"`
}

// BindResult 机构用户绑定结果
type BindResult struct {
	AgentId    int    `json:"agent_id"`
	DouyinId   string `json:"douyin_id"`
	BindTime   int    `json:"bind_time"`
	UnbindTime int    `json:"unbind_time"`
}

// GetAgencyUserBindRecord 获取机构用户绑定记录响应结果
type GetAgencyUserBindRecord struct {
	Results    []BindResult `json:"results"`
	TotalCount int          `json:"total_count"`
}

// DouYinTaskBoxGetAgencyUserBindRecordRes 获取机构用户绑定记录响应结果
type DouYinTaskBoxGetAgencyUserBindRecordRes struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	LogId  string `json:"log_id"`
	// Data 业务数据主体，包含具体的业务响应信息
	Data GetAgencyUserBindRecord `json:"data"`
}
