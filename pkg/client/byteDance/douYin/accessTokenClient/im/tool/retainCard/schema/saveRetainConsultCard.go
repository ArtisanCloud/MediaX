package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMToolRetainConsultCardReq 留资卡片创建请求参数
type DouYinIMToolRetainConsultCardReq struct {
	Title      string `json:"title"`
	MediaId    string `json:"media_id"`
	Components []int  `json:"components"`
}

// DouYinIMToolRetainConsultCardRes 留资卡片创建响应参数
type DouYinIMToolRetainConsultCardRes struct {
	// CardId 留资卡片ID，创建成功后返回的留资卡片ID。
	CardId string `json:"card_id"`
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
