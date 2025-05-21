package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangPromoteUnitUpdateStatusReq 更新推广单元状态请求参数
type JuGuangPromoteUnitUpdateStatusReq struct {
	AdvertiserID int64   `json:"advertiser_id"` // 广告主ID
	UnitIDs      []int64 `json:"unit_ids"`      // 单元ID列表
	Status       int     `json:"status"`        // 单元状态
}

// UnitUpdateStatusData 更新单元状态数据
type UnitUpdateStatusData struct {
	UnitIDs []int64 `json:"unit_ids"` // 单元ID列表
}

// JuGuangPromoteUnitUpdateStatusRes 更新推广单元状态响应参数
type JuGuangPromoteUnitUpdateStatusRes struct {
	response.RedBookAccessTokenRes
	Data UnitUpdateStatusData `json:"data"` // 更新单元状态数据
}
