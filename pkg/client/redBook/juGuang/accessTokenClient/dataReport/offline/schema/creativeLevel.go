package schema

type JuGuangDataReportOfflineCreativeLevelReq struct {
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

type JuGuangDataReportOfflineCreativeLevelRes struct {
	Code    int               `json:"code"`
	Msg     string            `json:"msg"`
	Success bool              `json:"success"`
	Data    CreativeLevelData `json:"data"`
}

type CreativeLevelData struct {
	DataList        []CreativeLevelItem `json:"data_list"`
	AggregationData AggregationData     `json:"aggregation_data"`
}

type CreativeLevelItem struct {
	CreativityID string  `json:"creativity_id"`
	Fee          float64 `json:"fee"`
	Impression   int64   `json:"impression"`
	Click        int64   `json:"click"`
	Interaction  int64   `json:"interaction"`
}
