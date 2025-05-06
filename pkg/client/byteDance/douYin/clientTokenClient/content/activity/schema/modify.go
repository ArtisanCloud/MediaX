package schema

// DouYinContentActivityModifyReq 更新活动请求参数
type DouYinContentActivityModifyReq struct {
	StartTime                  int64              `json:"start_time"`
	ActivityId                 int64              `json:"activity_id"`
	ModifyBusinessTaskInfoList []BusinessTaskInfo `json:"modify_business_task_info_list"`
	NotBcCheck                 bool               `json:"not_bc_check"`
	ActivityName               string             `json:"activity_name"`
	EndTime                    int64              `json:"end_time"`
	AddBusinessTaskInfoList    []BusinessTaskInfo `json:"add_business_task_info_list"`
	DelBusinessTaskIdList      []int64            `json:"del_business_task_id_list"`
}

// DouYinContentActivityModifyRes 更新活动响应结构
type DouYinContentActivityModifyRes struct {
	ModifyBusinessTaskIdList []int64 `json:"modify_business_task_id_list"` // 修改的业务任务ID列表
	ActivityId               int64   `json:"activity_id"`                  // 活动ID
	AddBusinessTaskIdList    []int64 `json:"add_business_task_id_list"`    // 新增的业务任务ID列表
	DelBusinessTaskIdList    []int64 `json:"del_business_task_id_list"`    // 删除的业务任务ID列表
}
