package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangDataReportRealtimeCampaignLevelReq 表示获取计划层级实时报表数据的请求结构体
type JuGuangDataReportRealtimeCampaignLevelReq struct {
	AdvertiserID            int64  `json:"advertiser_id"`
	StartDate               string `json:"start_date"`
	EndDate                 string `json:"end_date"`
	SortColumn              string `json:"sort_column,omitempty"`
	Sort                    string `json:"sort,omitempty"`
	PageNum                 int    `json:"page_num,omitempty"`
	PageSize                int    `json:"page_size,omitempty"`
	MarketingTargetList     []int  `json:"marketing_target_list,omitempty"`
	CampaignFilterState     int    `json:"campaign_filter_state,omitempty"`
	CampaignCreateBeginTime string `json:"campaign_create_begin_time,omitempty"`
	CampaignCreateEndTime   string `json:"campaign_create_end_time,omitempty"`
	PlacementList           []int  `json:"placement_list,omitempty"`
	LimitDayBudgetList      []int  `json:"limit_day_budget_list,omitempty"`
	OptimizeTargetList      []int  `json:"optimize_target_list,omitempty"`
	BuildTypeList           []int  `json:"build_type_list,omitempty"`
	BiddingStrategyList     []int  `json:"bidding_strategy_list,omitempty"`
	ConstraintTypeList      []int  `json:"constraint_type_list,omitempty"`
	PromotionTargetList     []int  `json:"promotion_target_list,omitempty"`
	CombineAuditStatus      int    `json:"combine_audit_status,omitempty"`
	MigrationStatusList     []int  `json:"migration_status_list,omitempty"`
	Name                    string `json:"name,omitempty"`
	ID                      int    `json:"id,omitempty"`
	DataCaliber             int    `json:"data_caliber,omitempty"`
	NeedHourlyData          bool   `json:"need_hourly_data,omitempty"`
}

// JuGuangDataReportRealtimeCampaignLevelRes 表示获取计划层级实时报表数据的响应结构体
type CampaignDTO struct {
	Data            DataReportDTO   `json:"data"`
	BaseCampaignDTO BaseCampaignDTO `json:"base_campaign_dto"`
	HourlyData      []DataReportDTO `json:"hourly_data"`
}

type JuGuangDataReportRealtimeCampaignLevelRes struct {
	response.RedBookAccessTokenRes
	Page         PageRespDTO   `json:"page"`
	CampaignDTOs []CampaignDTO `json:"campaign_dtos"`
	TotalData    DataReportDTO `json:"total_data"`
}
