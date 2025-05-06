package schema

import (
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// MaterialUpdateNewsReq 更新图文消息请求结构体
type MaterialUpdateNewsReq struct {
	MediaID  int64             `json:"media_id"` // 媒体ID
	Index    int64             `json:"index"`    // 要更新的文章在图文消息中的位置（从0开始）
	Articles []*object.HashMap `json:"articles"` // 图文消息列表
}
