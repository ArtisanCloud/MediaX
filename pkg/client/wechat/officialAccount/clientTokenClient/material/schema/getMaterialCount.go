package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"

// MaterialGetMaterialCountRes 获取素材总数响应结构体
type MaterialGetMaterialCountRes struct {
	response.OfficialAccountRes

	VoiceCount int    `json:"voice_count"` // 语音素材总数
	VideoCount int    `json:"video_count"` // 视频素材总数
	ImageCount int    `json:"image_count"` // 图片素材总数
	NewsCount  int    `json:"news_count"`  // 图文素材总数
	MediaID    string `json:"media_id"`    // 媒体ID
}
