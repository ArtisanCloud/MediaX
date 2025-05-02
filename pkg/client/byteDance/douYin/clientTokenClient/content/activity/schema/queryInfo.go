package schema

// DouYinContentActivityQueryInfoReq 查询活动信息请求参数
type DouYinContentActivityQueryInfoReq struct {
	ActivityId int64   `json:"activity_id"`  // 活动ID
	TaskIdList []int64 `json:"task_id_list"` // 任务ID列表
}

// ActivityInfo 活动信息
type ActivityInfo struct {
	ActivityName string `json:"activity_name"` // 活动名称
	StartTime    int64  `json:"start_time"`    // 活动开始时间，秒级时间戳
	EndTime      int64  `json:"end_time"`      // 活动结束时间，秒级时间戳
}

// DouYinContentActivityQueryInfoRes 查询活动信息响应结构
type DouYinContentActivityQueryInfoRes struct {
	BusinessTaskInfoList []BusinessTaskInfo `json:"business_task_info_list"` // 业务任务信息列表
	ActivityInfo         ActivityInfo       `json:"ActivityInfo"`            // 活动信息
}
