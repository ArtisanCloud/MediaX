package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangDataReportOfflineNoteLevelReq 表示获取笔记层级离线报表数据的请求结构体
type JuGuangDataReportOfflineNoteLevelReq struct {
	AdvertiserID    int64          `json:"advertiser_id"`
	StartDate       string         `json:"start_date"`
	EndDate         string         `json:"end_date"`
	TimeUnit        string         `json:"time_unit,omitempty"`
	MarketingTarget []int          `json:"marketing_target,omitempty"`
	BiddingStrategy []int          `json:"bidding_strategy,omitempty"`
	OptimizeTarget  []int          `json:"optimize_target,omitempty"`
	Placement       []int          `json:"placement,omitempty"`
	PromotionTarget []int          `json:"promotion_target,omitempty"`
	Programmatic    []int          `json:"programmatic,omitempty"`
	DeliveryMode    []int          `json:"delivery_mode,omitempty"`
	SplitColumns    []string       `json:"split_columns,omitempty"`
	SortColumn      string         `json:"sort_column,omitempty"`
	Sort            string         `json:"sort,omitempty"`
	PageNum         int            `json:"page_num,omitempty"`
	PageSize        int            `json:"page_size,omitempty"`
	DataCaliber     int            `json:"data_caliber,omitempty"`
	Filters         []FilterClause `json:"filters,omitempty"`
}

type NoteLevelData struct {
	TotalCount      int             `json:"total_count"`
	DataList        []DataReportDTO `json:"data_list"`
	AggregationData DataReportDTO   `json:"aggregation_data"`
}

// JuGuangDataReportOfflineNoteLevelRes 表示获取笔记层级离线报表数据的响应结构体
type JuGuangDataReportOfflineNoteLevelRes struct {
	response.RedBookAccessTokenRes
	Data NoteLevelData `json:"data"`
}
