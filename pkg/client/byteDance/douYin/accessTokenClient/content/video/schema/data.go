package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
)

// DouYinContentVideoDataReq 表示获取抖音视频数据的请求结构体。
// 目前该请求结构体为空，可根据实际需求添加字段。
type DouYinContentVideoDataReq struct {
}

// Data 表示视频的详细信息。
type Data struct {
	// Title 视频标题。
	Title string `json:"title"`
	// CreateTime 视频创建时间，Unix 时间戳格式。
	CreateTime int `json:"create_time"`
	// VideoStatus 视频状态。
	VideoStatus int `json:"video_status"`
	// ShareUrl 视频分享链接。
	ShareUrl string `json:"share_url"`
	// Cover 视频封面图片 URL。
	Cover string `json:"cover"`
	// IsTop 视频是否置顶。
	IsTop bool `json:"is_top"`
	// Statistics 视频统计数据。
	Statistics Statistics `json:"statistics"`
	// ItemId 视频唯一标识 ID。
	ItemId string `json:"item_id"`
	// IsReviewed 视频是否已审核。
	IsReviewed bool `json:"is_reviewed"`
	// MediaType 视频媒体类型。
	MediaType int `json:"media_type"`
}

// DouYinContentVideoDataRes 表示获取抖音视频数据的响应结构体。
type DouYinContentVideoDataRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// List 视频列表，包含多个视频信息。
		List []Data `json:"list"`
	}
}
