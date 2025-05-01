package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentVideoGetIFrameByVideoRes 抖音视频获取iframe响应结构体
type DouYinContentVideoGetIFrameByVideoRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// IframeCode iframe嵌入代码
		IframeCode string `json:"iframe_code"`
		// VideoWidth 视频宽度
		VideoWidth int `json:"video_width"`
		// VideoHeight 视频高度
		VideoHeight int `json:"video_height"`
		// VideoTitle 视频标题
		VideoTitle string `json:"video_title"`
	} `json:"data"`
}
