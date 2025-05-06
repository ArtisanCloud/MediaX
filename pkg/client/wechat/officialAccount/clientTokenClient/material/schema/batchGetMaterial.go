package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// MaterialBatchGetMaterialReq 批量获取素材请求结构体
type MaterialBatchGetMaterialReq struct {
	Type   string `json:"type"`   // 素材类型
	Offset int64  `json:"offset"` // 偏移量
	Count  int64  `json:"count"`  // 获取数量
}

// MaterialBatchGetMaterialRes 批量获取素材响应结构体
type MaterialBatchGetMaterialRes struct {
	response.OfficialAccountRes

	TotalCount int               `json:"total_count"` // 素材总数
	ItemCount  int               `json:"item_count"`  // 本次获取的素材数量
	Item       []*object.HashMap `json:"item"`        // 素材列表
}
