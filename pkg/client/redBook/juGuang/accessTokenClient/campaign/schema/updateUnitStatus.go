package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// ## JuGuangCampaignUpdateUnitStatusReq 修改广告单元状态请求参数
//
// 字段说明：
//   AdvertiserId - 广告主ID (必填)
//   CampaignIds  - 广告计划ID列表 (必填)
//   ActionType   - 操作类型 (必填，1: 启用, 2: 暂停)
type JuGuangCampaignUpdateUnitStatusReq struct {
	AdvertiserId int   `json:"advertiser_id"`
	CampaignIds  []int `json:"campaign_ids"`
	ActionType   int   `json:"action_type"`
}

// ## JuGuangCampaignUpdateUnitStatusRes 修改广告单元状态返回结果
//   JuGuangRes 基础响应
//   Data: 业务数据主体，包含具体的业务响应信息
//     • CampaignIds: 成功更新的广告计划ID列表

type JuGuangCampaignUpdateUnitStatusRes struct {
	response.JuGuangRes
	Data struct {
		CampaignIds []int `json:"campaign_ids"`
	} `json:"data"`
}
