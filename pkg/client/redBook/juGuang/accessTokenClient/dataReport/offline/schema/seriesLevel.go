package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// FilterClause 表示过滤条件结构体
type FilterClause struct {
	Column   string   `json:"column"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// JuGuangDataReportOfflineSeriesLevelReq 表示获取系列层级离线报表数据的请求结构体
type JuGuangDataReportOfflineSeriesLevelReq struct {
	AdvertiserID string         `json:"advertiser_id"`
	StartDate    string         `json:"start_date"`
	EndDate      string         `json:"end_date"`
	TimeUnit     string         `json:"time_unit,omitempty"`
	SortColumn   string         `json:"sort_column,omitempty"`
	Sort         string         `json:"sort,omitempty"`
	PageNum      int            `json:"page_num,omitempty"`
	PageSize     int            `json:"page_size,omitempty"`
	Filters      []FilterClause `json:"filters,omitempty"`
}

type SeriesLevelData struct {
	TotalCount      int             `json:"total_count"`
	DataList        []DataReportDTO `json:"data_list"`
	AggregationData DataReportDTO   `json:"aggregation_data"`
}

// JuGuangDataReportOfflineSeriesLevelRes 表示获取系列层级离线报表数据的响应结构体
type JuGuangDataReportOfflineSeriesLevelRes struct {
	response.RedBookAccessTokenRes
	Data SeriesLevelData `json:"data"`
}
