package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentCreatePostReq 创建任务请求参数
type DouYinContentCreatePostReq struct {
	// EndTime 任务结束时间，秒级时间戳，必填
	EndTime int64 `json:"end_time"`

	// StartTime 任务开始时间，秒级时间戳，必填
	StartTime int64 `json:"start_time"`

	// TaskCondition 任务条件，必填
	TaskCondition TaskCondition `json:"task_condition"`

	// TaskName 任务名称，长度不超过50个字符，必填
	TaskName string `json:"task_name"`
}

// TaskCondition 定义任务条件结构
type TaskCondition struct {
	// Condition 条件类型，例如："collection"
	Condition string `json:"condition"`

	// MinValue 最小值
	MinValue int64 `json:"min_value"`

	// MaxValue 最大值
	MaxValue int64 `json:"max_value"`
}

type DouYinContentCreatePostRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体
	Data struct {
		c.
		TaskId     string `json:"task_id"`
		TaskStatus int    `json:"task_status"`
	} `json:"data"`
}
