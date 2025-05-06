package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// Result 视频情况结果。
type VideoResult struct {
	// Date       日期，格式为 yyyy-MM-dd。
	Date string `json:"date"`
	// NewComment 新增评论数，单位：万
	NewIssue int `json:"new_issue"`
	// NewLike    新增点赞数，单位：万
	NewPlay int `json:"new_play"`
	// TotalIssue 累计播放数，单位：万
	TotalIssue int `json:"total_issue"`
}

// DouYinConnectionDataUserVideoStatusRes 获取用户视频情况响应。
type DouYinConnectionDataUserVideoStatusRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// ResultList 获取用户视频情况列表。
		ResultList []VideoResult `json:"result_list"`
	} `json:"data"`
}
