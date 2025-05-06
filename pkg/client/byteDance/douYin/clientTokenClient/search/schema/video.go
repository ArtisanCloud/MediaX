package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/content/video/schema"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
)

// DouYinSearchVideoReq 抖音视频搜索请求参数
type DouYinSearchVideoReq struct {
	Keyword  string `json:"keyword"`   // 搜索关键词
	Count    int    `json:"count"`     // 返回数量
	Cursor   int    `json:"cursor"`    // 游标，用于分页
	DeviceId string `json:"device_id"` // 设备ID
}

// DouYinSearchVideoRes 抖音视频搜索响应结构
type DouYinSearchVideoRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`

	// Data 业务数据主体。
	Data struct {
		Data struct {
			// Data通用返回信息
			response.DouYinRes
			// VideoList 视频列表 []schema.Video
			VideoList []schema.Video `json:"video_list"`
			// SearchId 搜索ID
			SearchId string `json:"search_id"`
		} `json:"data"`
	} `json:"data"`
}
