package schema

// DouYinTaskBoxAgencyQueryAppTaskReq 查询小程序任务请求结构体
type DouYinTaskBoxAgencyQueryAppTaskReq struct {
	AppId           string `json:"appid"`             // 查询的小程序 appid
	CreateEndTime   int64  `json:"create_end_time"`   // 任务创建终止时间，秒级时间戳
	CreateStartTime int64  `json:"create_start_time"` // 任务创建起始时间，秒级时间戳
	TaskCategory    int32  `json:"task_category"`     // 任务类别枚举
}

// DouYinTaskBoxAgencyQueryAppTaskRes 查询小程序任务响应结构体
type DouYinTaskBoxAgencyQueryAppTaskRes struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	LogId  string `json:"log_id"`
	Data   struct {
		TaskIds []int64 `json:"task_ids"`
	} `json:"data"`
}
