package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

type JuGuangPromoteCampaignUpdateStatusReq struct {
	AdvertiserID int64   `json:"advertiser_id"`
	CampaignIDs  []int64 `json:"campaign_ids"`
	ActionType   int     `json:"action_type"`
}

type CampaignUpdateStatusData struct {
	CampaignIDs []int64 `json:"campaign_ids"`
}

type JuGuangPromoteCampaignUpdateStatusRes struct {
	response.RedBookAccessTokenRes
	// Data 更新计划状态数据
	Data CampaignUpdateStatusData `json:"data"`
}
