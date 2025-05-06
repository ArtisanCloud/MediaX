package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// FanResult 粉丝情况结果。
type FanResult struct {
	// Date 日期，格式为 yyyy-MM-dd。
	Date string `json:"date"`
	// NewFans 新增粉丝数，单位：万
	NewFans string `json:"new_fans"`
	// TotalFans 累计粉丝数，单位：万
	TotalFans string `json:"total_fans"`
}

// DouYinConnectionDataUserFansCountRes 获取用户粉丝情况响应。
type DouYinConnectionDataUserFansCountRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// ResultList 获取用户粉丝情况列表。
		ResultList []FanResult `json:"result_list"`
	} `json:"data"`
}
