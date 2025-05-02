package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentSchemasGetItemInfoReq 视频详情页跳转链接获取请求参数
type DouYinContentSchemasGetItemInfoReq struct {
	ExpireAt int64  `json:"expire_at"` // 短链过期时间，必填
	ItemId   string `json:"item_id"`   // 视频id
	VideoId  int64  `json:"video_id"`  // 视频id
}

// DouYinContentSchemasGetItemInfoRes 视频详情页跳转链接获取响应结构
type DouYinContentSchemasGetItemInfoRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// Schema 视频上传的h5分享链接
		Schema string `json:"schema"`
	} `json:"data"`
}
