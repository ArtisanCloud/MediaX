package response

import "github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/core/response"

type MaterialAddNewsRes struct {
	response.OfficialAccountRes

	MediaID string `json:"media_id"`
}
