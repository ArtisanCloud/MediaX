package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMToolAppletTemplateSetReq 设置小程序模板请求参数
type DouYinIMToolAppletTemplateSetReq struct {
	AppId          string `json:"app_id"`                  // 小程序账号名
	Name           string `json:"name,omitempty"`          // 小程序/小游戏名
	IconMediaId    string `json:"icon_media_id,omitempty"` // 小程序/小游戏头像媒体ID
	CardTemplateId string `json:"card_template_id"`        // 卡片模板ID
	CardType       int32  `json:"card_type"`               // 卡片类型
	Content        string `json:"content"`                 // 卡片内容
	MediaId        string `json:"media_id"`                // 媒体ID
	Title          string `json:"title"`                   // 卡片标题
}

// DouYinIMToolAppletTemplateSetRes 设置小程序模板响应参数
type DouYinIMToolAppletTemplateSetRes struct {
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
