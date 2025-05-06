package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// LikeResult 点赞情况结果。
type LikeResult struct {
	// Date 日期，格式为 yyyy-MM-dd。
	Date string `json:"date"`
	// NewLike 新增点赞数，单位：万
	NewLike string `json:"new_like"`
}

// DouYinConnectionDataUserLikeNumberRes 获取用户点赞情况响应。
type DouYinConnectionDataUserLikeNumberRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// ResultList 获取用户点赞情况列表。
		ResultList []LikeResult `json:"result_list"`
	} `json:"data"`
}
