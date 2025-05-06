package schema

// DouYinTaskBoxQueryAgencyVideoDailyDataReq 查询视频每日数据请求结构体
type DouYinTaskBoxQueryAgencyVideoDailyDataReq struct {
	BillingDate           string `json:"billing_date"`             // 计费日期
	PageNum               int32  `json:"page_num"`                 // 页号
	PageSize              int32  `json:"page_size"`                // 页大小，最大不能超过10000
	VideoPublishEndTime   int64  `json:"video_publish_end_time"`   // 最大的视频发布时间戳
	VideoPublishStartTime int64  `json:"video_publish_start_time"` // 视频发布的起始时间戳
	AgentID               int64  `json:"agent_id"`                 // 抖音侧团长ID
	AppID                 string `json:"app_id"`                   // 撮合中介合作的小程序appid
	DouYinID              string `json:"douyin_id"`                // 达人抖音号
}

type QueryAgencyVideoDailyData struct {
	BillingGmv1D       int    `json:"billing_gmv_1d"`
	TalentProfit1D     int    `json:"talent_profit_1d"`
	ActiveCnt1D        int    `json:"active_cnt_1d"`
	Author             string `json:"author"`
	Gmv1D              int    `json:"gmv_1d"`
	RefundGmv1D        int    `json:"refund_gmv_1d"`
	VideoTitle         string `json:"video_title"`
	VideoId            int64  `json:"video_id"`
	DouYinId           string `json:"douyin_id"`
	VideoLink          string `json:"video_link"`
	AdShareCost1D      int    `json:"ad_share_cost_1d"`
	BillingRefundGmv1D int    `json:"billing_refund_gmv_1d"`
	TaskName           string `json:"task_name"`
	Date               string `json:"date"`
	PublishTime        int    `json:"publish_time"`
	FeedAdShareCost1D  int    `json:"feed_ad_share_cost_1d"`
	TaskId             int64  `json:"task_id"`
	MicroAppTitle      string `json:"micro_app_title"`
	AgentId            string `json:"agent_id"`
}

type DouYinTaskBoxQueryAgencyVideoDailyDataRes struct {
	ErrNo  int64  `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	LogId  string `json:"log_id"`
	Data   struct {
		Results []QueryAgencyVideoDailyData `json:"results"`
	} `json:"data"`
}
