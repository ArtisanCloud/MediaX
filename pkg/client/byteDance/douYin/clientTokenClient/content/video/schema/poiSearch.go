package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

type DouYinContentVideoPoiSearchReq struct {
	DefaultHashtag string `json:"default_hashtag,omitempty" url:"default_hashtag,omitempty"` // 追踪分享默认hashtag
	LinkParam      string `json:"link_param,omitempty" url:"link_param,omitempty"`           // 分享来源url附加参数（暂未开放）
	NeedCallback   bool   `json:"need_callback,omitempty" url:"need_callback,omitempty"`     // 如果需要知道视频分享成功的结果，need_callback设置为true
	SourceStyleID  string `json:"source_style_id,omitempty" url:"source_style_id,omitempty"` // 多来源样式id（暂未开放）
}

type DouYinContentVideoPoiSearchRes struct {
	Extra response.DouYinRes `json:"extra,omitempty"`
	Data  struct {
		ShareId     string `json:"share_id"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	} `json:"data"`
}
