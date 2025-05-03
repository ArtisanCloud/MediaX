package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMToolAppletTemplateGetReq 查询小程序引导卡片模板请求参数
type DouYinIMToolAppletTemplateGetReq struct {
	CardTemplateId string `json:"card_template_id"` // 卡片素材ID
	Count          int32  `json:"count"`            // 每页数量
	Cursor         int64  `json:"cursor"`           // 分页游标
	Status         int32  `json:"status"`           // 卡片状态
}

// AppletTemplate 小程序模板
type AppletTemplate struct {
	Status         int    `json:"status"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	CardTemplateId string `json:"card_template_id"`
	CardType       int    `json:"card_type"`
	MediaId        string `json:"media_id"`
	Name           string `json:"name"`
	AppIconUrl     string `json:"app_icon_url"`
	CreateTime     int    `json:"create_time"`
	UpdateTime     int    `json:"update_time"`
}

// DouYinIMToolAppletTemplateGetReq 获取小程序模板请求参数
type DouYinIMToolAppletTemplateGetRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes

		List []AppletTemplate `json:"list"`
	} `json:"data"`
}
