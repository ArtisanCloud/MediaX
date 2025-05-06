package schema

// DouYinTaskBoxAgencyQueryTaskInfoReq 查询任务信息请求结构体
type DouYinTaskBoxAgencyQueryTaskInfoReq struct {
	Appid              string `json:"appid"`                // 小程序appid
	PageNo             int32  `json:"page_no"`              // 分页编号，从1开始
	PageSize           int32  `json:"page_size"`            // 分页大小
	QueryParamsContent string `json:"query_params_content"` // 查询参数内容
	QueryParamsType    int32  `json:"query_params_type"`    // 查询参数类型
	TaskCategory       int32  `json:"task_category"`        // 任务类别
}

// 查询参数类型枚举
const (
	TASK_ID        = 1 // 任务ID
	TASK_NAME      = 2 // 任务名称
	TASK_PAGE_PATH = 3 // 任务页面地址
)

// 任务类别枚举
const (
	VIDEO_TASK = 1 // 短视频任务
	LIVE_TASK  = 2 // 直播任务
)

type OrientedTalentRel struct {
	CooperationState int    `json:"cooperation_state"`
	DouYinId         string `json:"douyin_id"`
	CancelOperator   int    `json:"cancel_operator,omitempty"`
}

type Task struct {
	AnchorTitle           string              `json:"anchor_title"`
	Appid                 string              `json:"appid"`
	PaymentAllocateRatio  int                 `json:"payment_allocate_ratio"`
	PlatformAddressApp    string              `json:"platform_address_app"`
	PlatformAddressWeb    string              `json:"platform_address_web"`
	ReferMaCaptures       []string            `json:"refer_ma_captures"`
	ReferVideoCaptures    []string            `json:"refer_video_captures"`
	RejectReason          string              `json:"reject_reason"`
	StartPage             string              `json:"start_page"`
	Status                int                 `json:"status"`
	TaskDesc              string              `json:"task_desc"`
	TaskEndTime           int                 `json:"task_end_time"`
	TaskIcon              string              `json:"task_icon"`
	TaskId                int64               `json:"task_id"`
	TaskName              string              `json:"task_name"`
	TaskRefundPeriod      int                 `json:"task_refund_period"`
	TaskSettleType        int                 `json:"task_settle_type"`
	TaskStartTime         int                 `json:"task_start_time"`
	TaskTags              []string            `json:"task_tags"`
	TaskType              int                 `json:"task_type"`
	OrientedTalentRelList []OrientedTalentRel `json:"oriented_talent_rel_list"`
}

type AgencyQueryTaskInfo struct {
	Total     int    `json:"total"`
	PageCount int    `json:"page_count"`
	AppId     string `json:"app_id"`
	Tasks     []Task `json:"tasks"`
}

type DouYinTaskBoxAgencyQueryTaskInfoRes struct {
	ErrNo  int                 `json:"err_no"`
	ErrMsg string              `json:"err_msg"`
	LogId  string              `json:"log_id"`
	Data   AgencyQueryTaskInfo `json:"data"`
}
