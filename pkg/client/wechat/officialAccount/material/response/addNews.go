package response

import "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"

type MaterialAddNewsRes struct {
	response.OfficialAccountRes

	MediaID string `json:"media_id"`
}
