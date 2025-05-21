package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// TimePeriod 定义推广时间段
type TimePeriod struct {
	Mon  string `json:"mon"`
	Tues string `json:"tues"`
	Wed  string `json:"wed"`
	Thur string `json:"thur"`
	Fri  string `json:"fri"`
	Sat  string `json:"sat"`
	Sun  string `json:"sun"`
}

// JuGuangPromoteCampaignUpdateReq 定义更新推广计划的请求结构
type JuGuangPromoteCampaignUpdateReq struct {
	AdvertiserID          int64       `json:"advertiser_id"`
	CampaignID            int64       `json:"campaign_id"`
	CampaignName          *string     `json:"campaign_name,omitempty"`
	LimitDayBudget        *int        `json:"limit_day_budget,omitempty"`
	CampaignDayBudget     *int        `json:"campaign_day_budget,omitempty"`
	SmartSwitch           *int        `json:"smart_switch,omitempty"`
	TimeType              *int        `json:"time_type,omitempty"`
	StartTime             *string     `json:"start_time,omitempty"`
	ExpireTime            *string     `json:"expire_time,omitempty"`
	TimePeriodType        *int        `json:"time_period_type,omitempty"`
	TimePeriod            *TimePeriod `json:"time_period,omitempty"`
	PacingMode            *int        `json:"pacing_mode,omitempty"`
	BiddingStrategy       *int        `json:"bidding_strategy,omitempty"`
	SearchFlag            *int        `json:"search_flag,omitempty"`
	TargetExtensionSwitch *int        `json:"target_extension_switch,omitempty"`
	SearchBidRatio        *float64    `json:"search_bid_ratio,omitempty"`
	DeeplinkID            *int64      `json:"deeplink_id,omitempty"`
	UniversalLinkID       *int64      `json:"universal_link_id,omitempty"`
	DetectUrlLink         *string     `json:"detect_url_link,omitempty"`
}

type CampaignUpdateData struct {
	CampaignID int64 `json:"campaign_id"`
}

// JuGuangPromoteCampaignUpdateRes 定义更新推广计划的响应结构
type JuGuangPromoteCampaignUpdateRes struct {
	response.RedBookAccessTokenRes
	// Data 更新计划数据
	Data CampaignUpdateData `json:"data"`
}
