package response

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type MaterialBatchGetMaterialRes struct {
	response.OfficialAccountRes

	TotalCount int               `json:"total_count"`
	ItemCount  int               `json:"item_count"`
	Item       []*object.HashMap `json:"item"`
}
