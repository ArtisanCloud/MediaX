package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// Media 视频列表，包含视频ID、URL列表等信息。
type Media struct {
	// MediaId 视频ID
	MediaId string `json:"media_id"`
	// UrlList URL列表，包含视频的不同分辨率URL。
	UrlList []string `json:"url_list"`
	// Status 视频状态
	Status string `json:"status"`
}

// DouYinToolMediaUploadRes 媒体上传响应
type DouYinToolMediaUploadRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// Media 视频列表，包含视频ID、URL列表等信息。
		Media Media `json:"media"`
	} `json:"data"`
}
