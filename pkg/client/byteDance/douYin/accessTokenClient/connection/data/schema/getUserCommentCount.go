package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// CommentResult 评论情况结果。
type CommentResult struct {
	// Date 日期，格式为 yyyy-MM-dd。
	Date string `json:"date"`
	// NewComment 新增评论数，单位：万
	NewComment string `json:"new_comment"`
}

// DouYinConnectionDataUserCommentCountRes 获取用户评论情况响应。
type DouYinConnectionDataUserCommentCountRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// ResultList 获取用户评论情况列表。
		ResultList []CommentResult `json:"result_list"`
	} `json:"data"`
}
