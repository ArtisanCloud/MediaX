package schema

// AgencyQueryVideoSum 视频汇总数据
type AgencyQueryVideoSum struct {
	// BillingGmvTd 当日GMV（Gross Merchandise Volume，商品交易总额）
	BillingGmvTd int64 `json:"billing_gmv_td"`
	// BillingRefundGmvTd 当日退款GMV
	BillingRefundGmvTd int64 `json:"biling_refund_gmv_td"`
	// VideoId 视频ID
	VideoId int64 `json:"video_id"`
	// TaskId 任务ID
	TaskId int64 `json:"task_id"`
	// AgentId 代理商ID
	AgentId string `json:"agent_id"`
	// PublishTime 发布时间（时间戳）
	PublishTime int `json:"publish_time"`
}

// DouYinTaskBoxAgencyQueryVideoSumRes 查询视频汇总数据响应
type DouYinTaskBoxAgencyQueryVideoSumRes struct {
	// LogId 日志ID，用于问题排查
	LogId string `json:"log_id"`
	// ErrNo 错误码
	ErrNo int `json:"err_no"`
	// ErrMsg 错误信息
	ErrMsg string `json:"err_msg"`
	// Data 业务数据主体
	Data struct {
		// Data 视频汇总数据列表
		Data []AgencyQueryVideoSum `json:"data"`
	} `json:"data"`
}
