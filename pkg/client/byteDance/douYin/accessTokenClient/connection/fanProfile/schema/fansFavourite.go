package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// Favourite 粉丝喜好数据
type Favourite struct {
	Rank     int    `json:"rank"`      // 排名
	Keyword  string `json:"keyword"`   // 关键词
	HotValue int    `json:"hot_value"` // 热度值
}

// DouYinConnectionFanProfileFansFavouriteRes 获取粉丝喜好数据接口响应结构
type DouYinConnectionFanProfileFansFavouriteRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// 粉丝喜好数据列表
		List []Favourite `json:"list"`
	} `json:"data"`
}
