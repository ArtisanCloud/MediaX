package schema

type JuGuangDataReportRealtimeTargetLevelReq struct {
	AdvertiserID        int64  `json:"advertiser_id"`
	PageNum             int    `json:"page_num,omitempty"`
	PageSize            int    `json:"page_size,omitempty"`
	StartDate           string `json:"start_date"`
	EndDate             string `json:"end_date"`
	SortColumn          string `json:"sort_column,omitempty"`
	Sort                string `json:"sort,omitempty"`
	Name                string `json:"name,omitempty"`
	MarketingTargetList []int  `json:"marketing_target_list,omitempty"`
	NeedHourlyData      bool   `json:"need_hourly_data,omitempty"`
}

type JuGuangDataReportRealtimeTargetLevelRes struct {
	Code       int           `json:"code"`
	Msg        string        `json:"msg"`
	Success    bool          `json:"success"`
	Page       PageRespDTO   `json:"page"`
	TargetDTOs []TargetDTO   `json:"target_dtos"`
	TotalData  DataReportDTO `json:"total_data"`
}

type PageRespDTO struct {
	PageIndex  int `json:"page_index"`
	TotalCount int `json:"total_count"`
}

type TargetDTO struct {
	Data            DataReportDTO   `json:"data"`
	BaseCampaignDTO BaseCampaignDTO `json:"base_campaign_dto"`
	BaseUnitDTO     BaseUnitDTO     `json:"base_unit_dto"`
	BaseTargetDTO   BaseTargetDTO   `json:"base_target_dto"`
	HourlyData      []DataReportDTO `json:"hourly_data,omitempty"`
}

type BaseCampaignDTO struct {
	CampaignID              int64  `json:"campaign_id"`
	CampaignName            string `json:"campaign_name"`
	CampaignFilterState     int    `json:"campaign_filter_state"`
	CampaignCreateTime      string `json:"campaign_create_time"`
	CampaignEnable          int    `json:"campaign_enable"`
	MarketingTarget         int    `json:"marketing_target"`
	Placement               int    `json:"placement"`
	OptimizeTarget          int    `json:"optimize_target"`
	PromotionTarget         int    `json:"promotion_target"`
	BiddingStrategy         int    `json:"bidding_strategy"`
	ConstraintType          int    `json:"constraint_type"`
	ConstraintValue         int    `json:"constraint_value"`
	LimitDayBudget          int    `json:"limit_day_budget"`
	OriginCampaignDayBudget int    `json:"origin_campaign_day_budget"`
	BudgetState             int    `json:"budget_state"`
	SmartSwitch             int    `json:"smart_switch"`
	PacingMode              int    `json:"pacing_mode"`
	StartTime               string `json:"start_time"`
	ExpireTime              string `json:"expire_time"`
	TimePeriod              string `json:"time_period"`
	TimePeriodType          int    `json:"time_period_type"`
	BuildType               int    `json:"build_type"`
	FeedFlag                int    `json:"feed_flag"`
	SearchFlag              int    `json:"search_flag"`
	MigrationStatus         int    `json:"migration_status"`
}

type BaseUnitDTO struct {
	UnitID          int64  `json:"unit_id"`
	UnitName        string `json:"unit_name"`
	UnitFilterState int    `json:"unit_filter_state"`
	UnitCreateTime  string `json:"unit_create_time"`
	UnitEnable      int    `json:"unit_enable"`
	CampaignID      int64  `json:"campaign_id"`
	EventBid        int    `json:"event_bid"`
}

type BaseTargetDTO struct {
	TargetID     int64  `json:"target_id"`
	TargetName   string `json:"target_name"`
	TargetStatus int    `json:"target_status"`
	UnitID       int64  `json:"unit_id"`
	CampaignID   int64  `json:"campaign_id"`
}

type DataReportDTO struct {
	Fee                         string `json:"fee"`
	Impression                  string `json:"impression"`
	Click                       string `json:"click"`
	CTR                         string `json:"ctr"`
	ACP                         string `json:"acp"`
	CPM                         string `json:"cpm"`
	Like                        string `json:"like"`
	Comment                     string `json:"comment"`
	Collect                     string `json:"collect"`
	Follow                      string `json:"follow"`
	Share                       string `json:"share"`
	Interaction                 string `json:"interaction"`
	CPI                         string `json:"cpi"`
	ActionButtonClick           string `json:"action_button_click"`
	ActionButtonCtr             string `json:"action_button_ctr"`
	Screenshot                  string `json:"screenshot"`
	PicSave                     string `json:"pic_save"`
	ReservePV                   string `json:"reserve_pv"`
	ClkLiveEntryPV              string `json:"clk_live_entry_pv"`
	ClkLiveEntryPVCost          string `json:"clk_live_entry_pv_cost"`
	ClkLiveAvgViewTime          string `json:"clk_live_avg_view_time"`
	ClkLiveAllFollow            string `json:"clk_live_all_follow"`
	ClkLive5sEntryPV            string `json:"clk_live_5s_entry_pv"`
	ClkLive5sEntryUVCost        string `json:"clk_live_5s_entry_uv_cost"`
	ClkLiveComment              string `json:"clk_live_comment"`
	SearchCmtClick              string `json:"search_cmt_click"`
	SearchCmtClickCVR           string `json:"search_cmt_click_cvr"`
	SearchCmtAfterRead          string `json:"search_cmt_after_read"`
	SearchCmtAfterReadAvg       string `json:"search_cmt_after_read_avg"`
	SellerVisit                 string `json:"seller_visit"`
	SellerVisitPrice            string `json:"seller_visit_price"`
	ShoppingCartAdd             string `json:"shopping_cart_add"`
	AddCartPrice                string `json:"add_cart_price"`
	PresaleOrderNum7d           string `json:"presale_order_num_7d"`
	PresaleOrderGmv7d           string `json:"presale_order_gmv_7d"`
	GoodsOrder                  string `json:"goods_order"`
	GoodsOrderPrice             string `json:"goods_order_price"`
	Rgmv                        string `json:"rgmv"`
	Roi                         string `json:"roi"`
	SuccessGoodsOrder           string `json:"success_goods_order"`
	ClickOrderCvr               string `json:"click_order_cvr"`
	PurchaseOrderPrice7d        string `json:"purchase_order_price_7d"`
	PurchaseOrderGmv7d          string `json:"purchase_order_gmv_7d"`
	PurchaseOrderRoi7d          string `json:"purchase_order_roi_7d"`
	ClkLiveRoomOrderNum         string `json:"clk_live_room_order_num"`
	LiveAverageOrderCost        string `json:"live_average_order_cost"`
	ClkLiveRoomRgmv             string `json:"clk_live_room_rgmv"`
	ClkLiveRoomRoi              string `json:"clk_live_room_roi"`
	Leads                       string `json:"leads"`
	LeadsCpl                    string `json:"leads_cpl"`
	LandingPageVisit            string `json:"landing_page_visit"`
	LeadsButtonImpression       string `json:"leads_button_impression"`
	ValidLeads                  string `json:"valid_leads"`
	ValidLeadsCpl               string `json:"valid_leads_cpl"`
	LeadsCvr                    string `json:"leads_cvr"`
	PhoneCallCnt                string `json:"phone_call_cnt"`
	PhoneCallSuccCnt            string `json:"phone_call_succ_cnt"`
	WechatCopyCnt               string `json:"wechat_copy_cnt"`
	WechatCopySuccCnt           string `json:"wechat_copy_succ_cnt"`
	IdentityCertiCnt            string `json:"identity_certi_cnt"`
	CommodityBuyCnt             string `json:"commodity_buy_cnt"`
	MessageUser                 string `json:"message_user"`
	Message                     string `json:"message"`
	MessageConsult              string `json:"message_consult"`
	InitiativeMessage           string `json:"initiative_message"`
	MessageConsultCpl           string `json:"message_consult_cpl"`
	InitiativeMessageCpl        string `json:"initiative_message_cpl"`
	MsgLeadsNum                 string `json:"msg_leads_num"`
	MsgLeadsCost                string `json:"msg_leads_cost"`
	ExternalGoodsVisit7         string `json:"external_goods_visit_7"`
	ExternalGoodsVisitPrice7    string `json:"external_goods_visit_price_7"`
	ExternalGoodsVisitRate7     string `json:"external_goods_visit_rate_7"`
	ExternalGoodsOrder7         string `json:"external_goods_order_7"`
	ExternalRgmv7               string `json:"external_rgmv_7"`
	ExternalGoodsOrderPrice7    string `json:"external_goods_order_price_7"`
	ExternalGoodsOrderRate7     string `json:"external_goods_order_rate_7"`
	ExternalRoi7                string `json:"external_roi_7"`
	ExternalGoodsOrder15        string `json:"external_goods_order_15"`
	ExternalRgmv15              string `json:"external_rgmv_15"`
	ExternalGoodsOrderPrice15   string `json:"external_goods_order_price_15"`
	ExternalGoodsOrderRate15    string `json:"external_goods_order_rate_15"`
	ExternalRoi15               string `json:"external_roi_15"`
	ExternalGoodsOrder30        string `json:"external_goods_order_30"`
	ExternalRgmv30              string `json:"external_rgmv_30"`
	ExternalGoodsOrderPrice30   string `json:"external_goods_order_price_30"`
	ExternalGoodsOrderRate30    string `json:"external_goods_order_rate_30"`
	ExternalRoi30               string `json:"external_roi_30"`
	WordAvgLocation             string `json:"word_avg_location"`
	WordImpressionRankFirst     string `json:"word_impression_rank_first"`
	WordImpressionRateFirst     string `json:"word_impression_rate_first"`
	WordImpressionRankThird     string `json:"word_impression_rank_third"`
	WordImpressionRateThird     string `json:"word_impression_rate_third"`
	WordClickRankFirst          string `json:"word_click_rank_first"`
	WordClickRateFirst          string `json:"word_click_rate_first"`
	WordClickRateThird          string `json:"word_click_rate_third"`
	WordClickRankThird          string `json:"word_click_rank_third"`
	InvokeAppOpenCnt            string `json:"invoke_app_open_cnt"`
	InvokeAppOpenCost           string `json:"invoke_app_open_cost"`
	InvokeAppEnterStoreCnt      string `json:"invoke_app_enter_store_cnt"`
	InvokeAppEnterStoreCost     string `json:"invoke_app_enter_store_cost"`
	InvokeAppEngagementCnt      string `json:"invoke_app_engagement_cnt"`
	InvokeAppEngagementCost     string `json:"invoke_app_engagement_cost"`
	InvokeAppPaymentCnt         string `json:"invoke_app_payment_cnt"`
	InvokeAppPaymentCost        string `json:"invoke_app_payment_cost"`
	SearchInvokeButtonClickCnt  string `json:"search_invoke_button_click_cnt"`
	SearchInvokeButtonClickCost string `json:"search_invoke_button_click_cost"`
	JdActiveUserNum             string `json:"jd_active_user_num"`
	JdActiveUserNumCvr          string `json:"jd_active_user_num_cvr"`
	JdActiveUserNumCpl          string `json:"jd_active_user_num_cpl"`
}
