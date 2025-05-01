// Package schema 定义了抖音内容视频相关接口的数据请求和响应结构体。
package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentVideoPoiSearchReq 表示查询视频携带地点信息（POI）的请求参数。
type DouYinContentVideoPoiSearchReq struct {
	// DefaultHashtag 追踪分享默认的 hashtag，用于定位特定内容。
	DefaultHashtag string `json:"default_hashtag,omitempty" url:"default_hashtag,omitempty"`

	// LinkParam 分享来源的 URL 附加参数（暂未开放）。
	LinkParam string `json:"link_param,omitempty" url:"link_param,omitempty"`

	// NeedCallback 是否需要回调视频分享成功的结果。
	NeedCallback bool `json:"need_callback,omitempty" url:"need_callback,omitempty"`

	// SourceStyleID 多来源样式 ID（暂未开放）。
	SourceStyleID string `json:"source_style_id,omitempty" url:"source_style_id,omitempty"`
}

// DouYinContentVideoPoiSearchRes 表示查询 POI 信息的响应结果。
type DouYinContentVideoPoiSearchRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`

	// Data 业务数据主体。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// ShareId 视频分享 ID。
		ShareId string `json:"share_id"`
	} `json:"data"`
}
