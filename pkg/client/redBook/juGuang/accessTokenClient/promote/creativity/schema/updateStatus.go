package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangPromoteCreativityUpdateStatusReq 更新创意状态请求参数
type JuGuangPromoteCreativityUpdateStatusReq struct {
	AdvertiserID  int64   `json:"advertiser_id"`  // 广告主ID
	CreativityIDs []int64 `json:"creativity_ids"` // 创意ID列表
	ActionType    int     `json:"action_type"`    // 操作类型
}

// CreativityUpdateStatusData 更新创意状态数据
type CreativityUpdateStatusData struct {
	CreativityIDs []int64 `json:"creativity_ids"` // 创意ID列表
}

// JuGuangPromoteCreativityUpdateStatusRes 更新创意状态响应参数
type JuGuangPromoteCreativityUpdateStatusRes struct {
	response.RedBookAccessTokenRes
	Data CreativityUpdateStatusData `json:"data"` // 更新创意状态数据
}
