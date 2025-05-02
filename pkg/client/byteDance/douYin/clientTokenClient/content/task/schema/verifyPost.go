package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

type DouYinContentVerifyPostReq struct {
	TaskId       string `json:"task_id"`
	TargetOpenId string `json:"target_open_id"`
	VideoId      string `json:"video_id"`
}

type DouYinContentVerifyPostRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
