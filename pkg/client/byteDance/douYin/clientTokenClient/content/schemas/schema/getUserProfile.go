package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentSchemasGetUserProfileReq 获取用户Profile请求参数
type DouYinContentSchemasGetUserProfileReq struct {
	// ExpireAt 过期时间戳
	ExpireAt int `json:"expire_at"`
	// OpenId 用户OpenId
	OpenId string `json:"open_id"`
}

// DouYinContentSchemasGetUserProfileRes 获取用户Profile响应结构
type DouYinContentSchemasGetUserProfileRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// Schema 获取用户Profile分享链接
		Schema string `json:"schema"`
	} `json:"data"`
}
