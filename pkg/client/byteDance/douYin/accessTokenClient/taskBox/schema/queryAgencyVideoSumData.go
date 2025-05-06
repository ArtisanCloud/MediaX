package schema

// DouYinTaskBoxQueryAgencyVideoSumDataReq 查询视频任务汇总数据请求结构体
type DouYinTaskBoxQueryAgencyVideoSumDataReq struct {
	PageNum               int32  `json:"page_num"`                 // 页号
	PageSize              int32  `json:"page_size"`                // 页大小，不能超过10000
	VideoPublishEndTime   int64  `json:"video_publish_end_time"`   // 视频发布的结束时间，秒级时间戳
	VideoPublishStartTime int64  `json:"video_publish_start_time"` // 视频发布的起始时间，秒级时间戳
	AgentID               int64  `json:"agent_id"`                 // 团长ID
	AppID                 string `json:"app_id"`                   // 机构合作小程序appid
	DouYinID              string `json:"douyin_id"`                // 抖音号
}

type QueryAgencyVideoSumData struct {
	VideoId            int64  `json:"video_id"`
	Author             string `json:"author"`
	Comments           int    `json:"comments"`
	FeedAdShareCostTd  int    `json:"feed_ad_share_cost_td"`
	VideoLink          string `json:"video_link"`
	Shares             int    `json:"shares"`
	MicroAppTitle      string `json:"micro_app_title"`
	PublishTime        int    `json:"publish_time"`
	AgentId            string `json:"agent_id"`
	ActiveCntTd        int    `json:"active_cnt_td"`
	RefundGmvTd        int    `json:"refund_gmv_td"`
	AdShareCostTd      int    `json:"ad_share_cost_td"`
	Likes              int    `json:"likes"`
	BillingRefundGmvTd int    `json:"billing_refund_gmv_td"`
	VideoTitle         string `json:"video_title"`
	TaskName           string `json:"task_name"`
	GmvTd              int    `json:"gmv_td"`
	TalentProfitTd     int    `json:"talent_profit_td"`
	DouYinId           string `json:"douyin_id"`
	Clicks             int    `json:"clicks"`
	Date               string `json:"date"`
	BillingGmvTd       int    `json:"billing_gmv_td"`
	TaskId             int64  `json:"task_id"`
	VideoViews         int    `json:"video_views"`
}

type DouYinTaskBoxQueryAgencyVideoSumDataRes struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	LogId  string `json:"log_id"`
	Data   struct {
		Results []QueryAgencyVideoSumData `json:"results"`
	} `json:"data"`
}
