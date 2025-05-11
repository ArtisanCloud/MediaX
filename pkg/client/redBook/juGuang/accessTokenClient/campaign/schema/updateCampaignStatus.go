package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// ## JuGuangCampaignUpdateCampaignStatusReq 更新广告计划状态请求
type JuGuangCampaignUpdateCampaignStatusReq struct {
	// 广告主ID
	AdvertiserId int `json:"advertiser_id"`
	// 计划ID列表	至少传一个限制单次变更计划数量，最多传20
	CampaignIds []int `json:"campaign_ids"`
	// 1：开启2：暂停3：删除
	ActionType int `json:"action_type"`
}

// ## JuGuangCampaignUpdateCampaignStatusRes 更新广告计划状态返回结果
type JuGuangCampaignUpdateCampaignStatusRes struct {
	// JuGuangRes 基础响应
	response.JuGuangRes
	// Data 业务数据主体，包含具体的业务响应信息
	Data struct {
		// CampaignIds 成功更新的广告计划ID列表
		CampaignIds []int `json:"campaign_ids"`
	} `json:"data"`
}
