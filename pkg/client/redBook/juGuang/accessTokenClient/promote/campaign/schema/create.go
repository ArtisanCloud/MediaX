package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangPromoteCampaignCreateReq 创建推广计划请求参数
type JuGuangPromoteCampaignCreateReq struct {
	AdvertiserID          int64  `json:"advertiser_id"`                     // 广告主ID
	MarketingTarget       int    `json:"marketing_target"`                  // 营销目标
	CampaignName          string `json:"campaign_name"`                     // 计划名称
	Placement             int    `json:"placement"`                         // 广告类型
	PromotionTarget       int    `json:"promotion_target"`                  // 推广标的类型
	Enable                *int   `json:"enable,omitempty"`                  // 计划状态
	TimeType              int    `json:"time_type"`                         // 时间类型
	StartTime             string `json:"start_time,omitempty"`              // 开始时间
	ExpireTime            string `json:"expire_time,omitempty"`             // 结束时间
	TimePeriodType        int    `json:"time_period_type"`                  // 时段类型
	TimePeriod            string `json:"time_period,omitempty"`             // 自定义时段
	BiddingStrategy       int    `json:"bidding_strategy"`                  // 出价方式
	LimitDayBudget        int    `json:"limit_day_budget"`                  // 预算类型
	CampaignDayBudget     *int   `json:"campaign_day_budget,omitempty"`     // 计划日预算
	OptimizeTarget        int    `json:"optimize_target"`                   // 推广目标
	ConstraintType        *int   `json:"constraint_type,omitempty"`         // 成本控制类型
	SmartSwitch           *int   `json:"smart_switch,omitempty"`            // 节假日预算上浮
	PacingMode            *int   `json:"pacing_mode,omitempty"`             // 投放速率
	BuildType             *int   `json:"build_type,omitempty"`              // 搭建类型
	DeliveryMode          *int   `json:"delivery_mode,omitempty"`           // 投放模式
	EventAssetID          *int64 `json:"event_asset_id,omitempty"`          // 资产ID
	AssetEvent            *int64 `json:"asset_event,omitempty"`             // 资产事件类型
	AssetEventID          *int64 `json:"asset_event_id,omitempty"`          // 资产事件ID
	PageCategory          *int   `json:"page_category,omitempty"`           // 落地页类型
	SearchFlag            *int   `json:"search_flag,omitempty"`             // 搜索快投
	TargetExtensionSwitch *int   `json:"target_extension_switch,omitempty"` // 搜索快投定向拓展
	SearchBidRatio        *int   `json:"search_bid_ratio,omitempty"`        // 搜索快投-出价系数
	DeeplinkID            *int64 `json:"deeplink_id,omitempty"`             // deeplink链接ID
	UniversalLinkID       *int64 `json:"universal_link_id,omitempty"`       // ulk的链接ID
	DetectUrlLink         string `json:"detect_url_link,omitempty"`         // 监测链接
}

type CampaignCreateData struct {
	CampaignID int64 `json:"campaign_id"` // 计划ID
}

// JuGuangPromoteCampaignCreateRes 创建推广计划响应参数
type JuGuangPromoteCampaignCreateRes struct {
	response.RedBookAccessTokenRes
	// Data 创建计划数据
	Data CampaignCreateData `json:"data"`
}
