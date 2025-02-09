package response

import "github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/core/response"

type MaterialAddMaterialRes struct {
	response.OfficialAccountRes

	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}
