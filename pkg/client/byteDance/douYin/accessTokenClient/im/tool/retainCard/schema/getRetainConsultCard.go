package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// Card 留资卡片信息
type Card struct {
	CardId     string `json:"card_id"`
	Title      string `json:"title"`
	MediaId    string `json:"media_id"`
	Components []int  `json:"components"`
	Status     int    `json:"status"`
}

// DouYinIMToolRetainConsultCardReq 留资卡片请求参数
type DouYinIMToolGetRetainConsultCardRes struct {
	// Cards 留资卡片列表，每个卡片包含 card_id、title、media_id、components 和 status 等字段。
	Cards []Card `json:"cards"`
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
