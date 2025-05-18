package schema

// JuGuangDataReportOfflineSpuLevelReq 表示SPU级别离线数据报表请求参数
type JuGuangDataReportOfflineSpuLevelReq struct {
	AdvertiserID int64  `json:"advertiser_id"` // 广告主ID
	StartDate    string `json:"start_date"`    // 开始时间，格式 yyyy-MM-dd
	EndDate      string `json:"end_date"`      // 结束时间，格式 yyyy-MM-dd
	TimeUnit     string `json:"time_unit"`     // 时间维度："DAY"：分天 "SUMMARY"：汇总
	SortColumn   string `json:"sort_column"`   // 排序字段
	Sort         string `json:"sort"`          // 升降序asc：升序desc：降序
	PageNum      int    `json:"page_num"`      // 页数，默认1
	PageSize     int    `json:"page_size"`     // 页大小，默认20,最大500
}

// SpuLevelData 表示SPU级别离线数据报表数据
type SpuLevelData struct {
	TotalCount      int             `json:"total_count"`      // 总条数
	DataList        []DateReportDTO `json:"data_list"`        // 详细数据
	AggregationData DateReportDTO   `json:"aggregation_data"` // 汇总数据
}

// JuGuangDataReportOfflineSpuLevelRes 表示SPU级别离线数据报表响应
type JuGuangDataReportOfflineSpuLevelRes struct {
	Code    int          `json:"code"`    // 返回码
	Msg     string       `json:"msg"`     // 返回信息
	Success bool         `json:"success"` // 接口是否成功
	Data    SpuLevelData `json:"data"`    // SPU数据
}

// DateReportDTO 表示数据报表DTO
type DateReportDTO struct {
	SpuID                 string `json:"spu_id"`
	SpuName               string `json:"spu_name"`
	Time                  string `json:"time"`
	Fee                   string `json:"fee"`
	Impression            string `json:"impression"`
	Click                 string `json:"click"`
	Ctr                   string `json:"ctr"`
	Acp                   string `json:"acp"`
	Cpm                   string `json:"cpm"`
	Like                  string `json:"like"`
	Comment               string `json:"comment"`
	Collect               string `json:"collect"`
	Follow                string `json:"follow"`
	Share                 string `json:"share"`
	Interaction           string `json:"interaction"`
	Cpi                   string `json:"cpi"`
	ActionButtonClick     string `json:"action_button_click"`
	ActionButtonCtr       string `json:"action_button_ctr"`
	Screenshot            string `json:"screenshot"`
	PicSave               string `json:"pic_save"`
	Leads                 string `json:"leads"`
	LeadsCpl              string `json:"leads_cpl"`
	LandingPageVisit      string `json:"landing_page_visit"`
	LeadsButtonImpression string `json:"leads_button_impression"`
	ValidLeads            string `json:"valid_leads"`
	ValidLeadsCpl         string `json:"valid_leads_cpl"`
	LeadsCvr              string `json:"leads_cvr"`
	ExternalLeads         string `json:"external_leads"`
	ExternalLeadsCpl      string `json:"external_leads_cpl"`
	GoodsVisit            string `json:"goods_visit"`
	GoodsVisitPrice       string `json:"goods_visit_price"`
	ShoppingCartAdd       string `json:"shopping_cart_add"`
	AddCartPrice          string `json:"add_cart_price"`
	IUserNum              string `json:"i_user_num"`
	TiUserNum             string `json:"ti_user_num"`
	IUserPrice            string `json:"i_user_price"`
	TiUserPrice           string `json:"ti_user_price"`
}
