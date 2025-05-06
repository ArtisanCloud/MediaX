package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// MaterialAddNewsReq 添加图文消息请求结构体
type MaterialAddNewsReq struct {
	Articles []*object.HashMap `json:"articles"` // 图文消息列表
}

// MaterialAddNewsRes 添加图文消息响应结构体
type MaterialAddNewsRes struct {
	response.OfficialAccountRes

	MediaID string `json:"media_id"` // 媒体ID
}
