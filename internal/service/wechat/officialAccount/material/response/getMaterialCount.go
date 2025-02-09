package response

import (
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/core/response"
)

type MaterialGetMaterialCountRes struct {
	response.OfficialAccountRes

	VoiceCount int    `json:"voice_count"`
	VideoCount int    `json:"video_count"`
	ImageCount int    `json:"image_count"`
	NewsCount  int    `json:"news_count"`
	MediaID    string `json:"media_id"`
}
