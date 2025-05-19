package schema

// JuGuangDataReportRealtimeUnitLevelReq 表示获取单元层级实时报表数据的请求结构体
type JuGuangDataReportRealtimeUnitLevelReq struct {
	AdvertiserID        int64  `json:"advertiser_id"`
	StartDate           string `json:"start_date"`
	EndDate             string `json:"end_date"`
	PageNum             int    `json:"page_num,omitempty"`
	PageSize            int    `json:"page_size,omitempty"`
	SortColumn          string `json:"sort_column,omitempty"`
	Sort                string `json:"sort,omitempty"`
	MarketingTargetList []int  `json:"marketing_target_list,omitempty"`
	UnitFilterState     int    `json:"unit_filter_state,omitempty"`
	UnitCreateBeginTime string `json:"unit_create_begin_time,omitempty"`
	UnitCreateEndTime   string `json:"unit_create_end_time,omitempty"`
	PlacementList       []int  `json:"placement_list,omitempty"`
	BiddingStrategyList []int  `json:"bidding_strategy_list,omitempty"`
	PromotionTargetList []int  `json:"promotion_target_list,omitempty"`
	CombineAuditStatus  int    `json:"combine_audit_status,omitempty"`
	Name                string `json:"name,omitempty"`
	ID                  int    `json:"id,omitempty"`
	DataCaliber         int    `json:"data_caliber,omitempty"`
	NeedHourlyData      bool   `json:"need_hourly_data,omitempty"`
}

// JuGuangDataReportRealtimeUnitLevelRes 表示获取单元层级实时报表数据的响应结构体
type JuGuangDataReportRealtimeUnitLevelRes struct {
	Code      int           `json:"code"`
	Msg       string        `json:"msg"`
	Success   bool          `json:"success"`
	Page      PageRespDTO   `json:"page"`
	UnitDTOs  []UnitDTO     `json:"unit_dtos"`
	TotalData DataReportDTO `json:"total_data"`
}

type UnitDTO struct {
	Data            DataReportDTO   `json:"data"`
	BaseUnitDTO     BaseUnitDTO     `json:"base_unit_dto"`
	BaseCampaignDTO BaseCampaignDTO `json:"base_campaign_dto"`
	HourlyData      []DataReportDTO `json:"hourly_data,omitempty"`
}
