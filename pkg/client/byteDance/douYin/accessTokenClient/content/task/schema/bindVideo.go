package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentBindVideoReq 绑定视频请求参数
type DouYinContentBindVideoReq struct {
	// TaskId 任务ID，创建任务之后获取的任务ID，必填
	TaskId string `json:"task_id"`

	// VideoId 视频ID，必填
	VideoId string `json:"video_id"`
}

type DouYinContentBindVideoRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
