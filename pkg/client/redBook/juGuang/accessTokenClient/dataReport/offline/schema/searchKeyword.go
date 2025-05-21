package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangDataReportOfflineSearchKeywordReq 表示获取搜索关键词层级离线报表数据的请求结构体
type JuGuangDataReportOfflineSearchKeywordReq struct {
	AdvertiserID    int64  `json:"advertiser_id"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	TimeUnit        string `json:"time_unit,omitempty"`
	MarketingTarget []int  `json:"marketing_target,omitempty"`
	BiddingStrategy []int  `json:"bidding_strategy,omitempty"`
	OptimizeTarget  []int  `json:"optimize_target,omitempty"`
	Placement       []int  `json:"placement,omitempty"`
	PromotionTarget []int  `json:"promotion_target,omitempty"`
	Programmatic    []int  `json:"programmatic,omitempty"`
	BuildType       []int  `json:"build_type,omitempty"`
	SortColumn      string `json:"sort_column,omitempty"`
	Sort            string `json:"sort,omitempty"`
	PageNum         int    `json:"page_num,omitempty"`
	PageSize        int    `json:"page_size,omitempty"`
	DataCaliber     int    `json:"data_caliber,omitempty"`
}

type PageRespDTO struct {
	PageIndex  int `json:"page_index"`
	TotalCount int `json:"total_count"`
}

type SearchKeywordData struct {
	Page            PageRespDTO     `json:"page"`
	DataList        []DataReportDTO `json:"data_list"`
	AggregationData DataReportDTO   `json:"aggregation_data"`
}

// JuGuangDataReportOfflineSearchKeywordRes 表示获取搜索关键词层级离线报表数据的响应结构体
type JuGuangDataReportOfflineSearchKeywordRes struct {
	response.RedBookAccessTokenRes
	Data SearchKeywordData `json:"data"`
}
