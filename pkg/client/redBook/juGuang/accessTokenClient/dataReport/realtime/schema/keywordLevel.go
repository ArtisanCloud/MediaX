package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangDataReportRealtimeKeywordLevelReq 表示获取关键词层级实时报表数据的请求结构体
type JuGuangDataReportRealtimeKeywordLevelReq struct {
	AdvertiserID       int64  `json:"advertiser_id"`
	PageNum            int    `json:"page_num,omitempty"`
	PageSize           int    `json:"page_size,omitempty"`
	StartDate          string `json:"start_date"`
	EndDate            string `json:"end_date"`
	SortColumn         string `json:"sort_column,omitempty"`
	Sort               string `json:"sort,omitempty"`
	KeywordFilterState int    `json:"keyword_filter_state,omitempty"`
	UseBidStrategy     int    `json:"use_bid_strategy,omitempty"`
	KeywordName        string `json:"keyword_name,omitempty"`
	CampaignName       string `json:"campaign_name,omitempty"`
	UnitName           string `json:"unit_name,omitempty"`
	DataCaliber        int    `json:"data_caliber,omitempty"`
	NeedHourlyData     bool   `json:"need_hourly_data,omitempty"`
}

// JuGuangDataReportRealtimeKeywordLevelRes 表示获取关键词层级实时报表数据的响应结构体
type BaseKeywordDTO struct {
	KeywordID          int64  `json:"keyword_id"`
	Keyword            string `json:"keyword"`
	UseBidStrategy     int    `json:"use_bid_strategy"`
	KeywordEnable      int    `json:"keyword_enable"`
	KeywordFilterState int    `json:"keyword_filter_state"`
	UnitID             int64  `json:"unit_id"`
	CampaignID         string `json:"campaign_id"`
}

type KeywordDTO struct {
	Data            DataReportDTO   `json:"data"`
	BaseCampaignDTO BaseCampaignDTO `json:"base_campaign_dto"`
	BaseUnitDTO     BaseUnitDTO     `json:"base_unit_dto"`
	BaseKeywordDTO  BaseKeywordDTO  `json:"base_keyword_dto"`
	SubKeywordDtos  []KeywordDTO    `json:"sub_keyword_dtos"`
	HourlyData      []DataReportDTO `json:"hourly_data"`
}

type JuGuangDataReportRealtimeKeywordLevelRes struct {
	response.RedBookAccessTokenRes
	Code        int           `json:"code"`
	Msg         string        `json:"msg"`
	Success     bool          `json:"success"`
	Page        PageRespDTO   `json:"page"`
	KeywordDtos []KeywordDTO  `json:"keyword_dtos"`
	TotalData   DataReportDTO `json:"total_data"`
}
