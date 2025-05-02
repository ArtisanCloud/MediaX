package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentSchemasGetChatReq 个人会话页跳转链接获取请求参数
type DouYinContentSchemasGetChatReq struct {
	// ExpireAt 过期时间戳
	ExpireAt int `json:"expire_at"`
	// OpenId 用户OpenId
	OpenId string `json:"open_id"`
}

// DouYinContentSchemasGetChatRes 个人会话页跳转链接获取响应结构
type DouYinContentSchemasGetChatRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// Schema 个人会话页跳转链接获取分享链接
		Schema string `json:"schema"`
	} `json:"data"`
}
