package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliVideoDeleteReq 表示删除B站视频的请求参数
type BiliBiliVideoDeleteReq struct {
	ResourceID string `json:"resource_id"` // 稿件唯一ID
}

// BiliBiliVideoDeleteRes 表示删除B站视频的响应参数
type BiliBiliVideoDeleteRes struct {
	response.BiliBiliRes
}
