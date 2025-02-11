package response

import "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"

type MaterialAddMaterialRes struct {
	response.OfficialAccountRes

	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}
