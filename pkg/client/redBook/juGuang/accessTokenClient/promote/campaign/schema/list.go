package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangPromoteCampaignListReq 表示获取推广计划列表的请求参数
type JuGuangPromoteCampaignListReq struct {
	AdvertiserID int64    `json:"advertiser_id"`
	CampaignIDs  []int64  `json:"campaign_ids,omitempty"`
	StartTime    string   `json:"start_time,omitempty"`
	ExpireTime   string   `json:"expire_time,omitempty"`
	Status       int      `json:"status,omitempty"`
	Page         PageInfo `json:"page,omitempty"`
}

// PageInfo 分页信息
type PageInfo struct {
	PageIndex int `json:"page_index,omitempty"`
	PageSize  int `json:"page_size,omitempty"`
}

type CampaignListData struct {
	Page             PageInfo          `json:"page"`
	BaseCampaignDTOs []BaseCampaignDTO `json:"base_campaign_dtos"`
}

// JuGuangPromoteCampaignListRes 表示获取推广计划列表的响应参数
type JuGuangPromoteCampaignListRes struct {
	response.RedBookAccessTokenRes
	// Data 推广计划列表数据
	Data CampaignListData `json:"data"`
}

// BaseCampaignDTO 表示推广计划基础信息
type BaseCampaignDTO struct {
	CampaignID          int64   `json:"campaign_id"`
	CampaignName        int64   `json:"campaign_name"`
	CampaignFilterState int     `json:"campaign_filter_state"`
	CampaignCreateTime  string  `json:"campaign_create_time"`
	CampaignEnable      int     `json:"campaign_enable"`
	MarketingTarget     int     `json:"marketing_target"`
	Placement           int     `json:"placement"`
	OptimizeTarget      int     `json:"optimize_target"`
	PromotionTarget     int     `json:"promotion_target"`
	BiddingStrategy     int     `json:"bidding_strategy"`
	ConstraintType      int     `json:"constraint_type"`
	ConstraintValue     int     `json:"constraint_value"`
	LimitDayBudget      int     `json:"limit_day_budget"`
	CampaignDayBudget   int     `json:"campaign_day_budget"`
	BudgetState         int     `json:"budget_state"`
	SmartSwitch         int     `json:"smart_switch"`
	Platform            int     `json:"platform"`
	PacingMode          int     `json:"pacing_mode"`
	StartTime           string  `json:"start_time"`
	ExpireTime          string  `json:"expire_time"`
	TimePeriod          string  `json:"time_period"`
	TimePeriodType      int     `json:"time_period_type"`
	FeedFlag            int     `json:"feed_flag"`
	BuildType           int     `json:"build_type"`
	CreativityState     int     `json:"creativity_state"`
	EventAssetID        int64   `json:"event_asset_id"`
	AssetEvent          int64   `json:"asset_event"`
	AssetEventID        int64   `json:"asset_event_id"`
	PageCategory        int     `json:"page_category"`
	SearchFlag          int     `json:"search_flag"`
	SearchBidRatio      float64 `json:"search_bid_ratio"`
	DeeplinkID          int64   `json:"deeplink_id"`
	UniversalLinkID     int64   `json:"universal_link_id"`
	DetectUrlLink       string  `json:"detect_url_link"`
}
