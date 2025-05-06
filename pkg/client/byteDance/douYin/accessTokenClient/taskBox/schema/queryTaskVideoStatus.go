package schema

// DouYinTaskBoxQueryTaskVideoStatusReq 查询任务视频状态请求结构体
type DouYinTaskBoxQueryTaskVideoStatusReq struct {
	PageNum               int32  `json:"page_num"`                 // 页号
	PageSize              int32  `json:"page_size"`                // 页大小，最大不能超过100
	VideoPublishEndTime   int64  `json:"video_publish_end_time"`   // 最大的视频发布时间戳
	VideoPublishStartTime int64  `json:"video_publish_start_time"` // 视频发布的起始时间戳
	AgentID               int64  `json:"agent_id"`                 // 抖音侧团长ID
	AppID                 string `json:"app_id"`                   // 撮合中介合作的小程序appid
	DouYinID              string `json:"douyin_id"`                // 达人抖音号
}

type QueryTaskVideoStatus struct {
	VideoId     int `json:"VideoId"`
	VideoStatus int `json:"VideoStatus"`
}

type DouYinTaskBoxQueryTaskVideoStatusRes struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	LogId  string `json:"log_id"`
	Data   struct {
		Results []QueryTaskVideoStatus `json:"results"`
	} `json:"data"`
}
