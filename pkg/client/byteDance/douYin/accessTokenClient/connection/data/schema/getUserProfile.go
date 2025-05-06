package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// UserProfile 用户画像结果。
type UserProfile struct {
	// Date 日期，格式为 yyyy-MM-dd。
	Date string `json:"date"`
	// ProfileUv 画像UV，即用户画像的访问量。
	ProfileUv string `json:"profile_uv"`
}

// DouYinConnectionDataUserProfile 抖音用户画像响应。
type DouYinConnectionDataUserProfile struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// ResultList 用户画像列表。
		ResultList []UserProfile `json:"result_list"`
	} `json:"data"`
}
