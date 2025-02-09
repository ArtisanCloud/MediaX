package response

import "github.com/ArtisanCloud/MediaX/internal/service/wechat/core/response"

type MaterialAddNewsRes struct {
	response.OfficialAccountRes

	MediaID string `json:"media_id"`
}
