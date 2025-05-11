package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// ## JuGuangCampaignUpdateCreativeStatusReq 修改广告创意状态请求参数
type JuGuangCampaignUpdateCreativeStatusReq struct {
	// 广告主ID
	AdvertiserId int `json:"advertiser_id"`
	// 创意ID列表	至少传一个限制单次变更创意数量，最多传20
	CreativityIds []int `json:"creativity_ids"`
	// 1：开启2：暂停3：删除
	ActionType int `json:"action_type"`
}

// ## JuGuangCampaignUpdateCreativeStatusRes 修改广告创意状态返回结果
type JuGuangCampaignUpdateCreativeStatusRes struct {
	// JuGuangRes 基础响应
	response.JuGuangRes
	// Data 业务数据主体，包含具体的业务响应信息
	Data struct {
		// CampaignIds 成功更新的广告计划ID列表
		CreativityIds []int `json:"creativity_ids"`
	} `json:"data"`
}
