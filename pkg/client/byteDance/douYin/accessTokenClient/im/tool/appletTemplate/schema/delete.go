package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMToolRetainConsultCardReq 删除小程序引导卡片模板返回参数
type DouYinIMToolAppletTemplateDeleteRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// CardTemplateId 卡片模板ID
		CardTemplateId string `json:"card_template_id"`
	} `json:"data"`
}
