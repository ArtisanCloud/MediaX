package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

type JuGuangDataReportOfflineCampaignLevelReq struct {
	AdvertiserID    int64          `json:"advertiser_id"`
	StartDate       string         `json:"start_date"`
	EndDate         string         `json:"end_date"`
	TimeUnit        string         `json:"time_unit,omitempty"`
	MarketingTarget []int          `json:"marketing_target,omitempty"`
	BiddingStrategy []int          `json:"bidding_strategy,omitempty"`
	OptimizeTarget  []int          `json:"optimize_target,omitempty"`
	Placement       []int          `json:"placement,omitempty"`
	PromotionTarget []int          `json:"promotion_target,omitempty"`
	BuildType       []int          `json:"build_type,omitempty"`
	DeliveryMode    []int          `json:"delivery_mode,omitempty"`
	SplitColumns    []string       `json:"split_columns,omitempty"`
	SortColumn      string         `json:"sort_column,omitempty"`
	Sort            string         `json:"sort,omitempty"`
	PageNum         int            `json:"page_num,omitempty"`
	PageSize        int            `json:"page_size,omitempty"`
	DataCaliber     int            `json:"data_caliber,omitempty"`
	Filters         []FilterClause `json:"filters,omitempty"`
}

type JuGuangDataReportOfflineCampaignLevelRes struct {
	response.RedBookAccessTokenRes
	Data CampaignLevelData `json:"data"`
}

type CampaignLevelData struct {
	DataList        []map[string]DataReportDTO `json:"data_list"`
	AggregationData map[string]DataReportDTO   `json:"aggregation_data"`
}
