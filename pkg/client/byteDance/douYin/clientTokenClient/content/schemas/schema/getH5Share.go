package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// MicroAppInfo 微应用信息
type MicroAppInfo struct {
	// AppId 应用ID
	AppId string `json:"app_id"`
	// AppTitle 应用标题
	AppTitle string `json:"app_title"`
	// AppUrl 应用URL
	AppUrl string `json:"app_url"`
	// Description 应用描述
	Description string `json:"description"`
}

// DouYinContentSchemasGetH5ShareReq 获取H5分享链接请求参数
type DouYinContentSchemasGetH5ShareReq struct {
	// ClientTicket 客户端票据
	ClientTicket string `json:"client_ticket"`
	// ExpireAt 过期时间戳
	ExpireAt int `json:"expire_at"`
	// HashtagList 话题标签列表
	HashtagList []string `json:"hashtag_list"`
	// MicroAppInfo 微应用信息
	MicroAppInfo MicroAppInfo `json:"micro_app_info"`
	// PoiId POI ID
	PoiId string `json:"poi_id"`
	// ShareToPublish 是否分享后发布
	ShareToPublish int `json:"share_to_publish"`
	// State 状态信息
	State string `json:"state"`
	// Title 标题
	Title string `json:"title"`
	// DouYinContentSchemasGetH5ShareRes 获取H5分享链接响应结构
	// VideoPath 视频路径
	VideoPath string `json:"video_path"`
}

// DouYinContentSchemasGetH5ShareRes 获取H5分享链接响应结构
type DouYinContentSchemasGetH5ShareRes struct {
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
