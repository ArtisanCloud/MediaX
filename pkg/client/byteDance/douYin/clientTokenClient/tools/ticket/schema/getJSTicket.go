package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinToolTicketGetJSTicketRes 获取js ticket
type DouYinToolTicketGetJSTicketRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// ExpiresIn 凭证有效时间，单位：秒
		ExpiresIn int `json:"expires_in"`
		// Ticket 凭证
		Ticket string `json:"ticket"`
	} `json:"data"`
}
