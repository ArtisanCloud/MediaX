package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

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
	DataList        []DataReportDTO `json:"data_list"`        // 详细数据
	AggregationData DataReportDTO   `json:"aggregation_data"` // 汇总数据
}

// JuGuangDataReportOfflineSpuLevelRes 表示SPU级别离线数据报表响应
type JuGuangDataReportOfflineSpuLevelRes struct {
	response.RedBookAccessTokenRes
	Data SpuLevelData `json:"data"` // SPU数据
}

// DataReportDTO 表示数据报表DTO
type DataReportDTO struct {
	SpuID                       string `json:"spu_id"`
	AddWechatCount              string `json:"add_wechat_count"`
	AddWechatCost               string `json:"add_wechat_cost"`
	AddWechatSucCount           string `json:"add_wechat_suc_count"`
	AddWechatSucCost            string `json:"add_wechat_suc_cost"`
	WechatTalkCount             string `json:"wechat_talk_count"`
	WechatTalkCost              string `json:"wechat_talk_cost"`
	ShopPoiClickNum             string `json:"shop_poi_click_num"`
	ShopPoiPagePv               string `json:"shop_poi_page_pv"`
	ShopPoiPageVisitPrice       string `json:"shop_poi_page_visit_price"`
	ShopPoiPageNavigateClick    string `json:"shop_poi_page_navigate_click"`
	SpuName                     string `json:"spu_name"`
	Time                        string `json:"time"`
	Fee                         string `json:"fee"`
	Impression                  string `json:"impression"`
	Click                       string `json:"click"`
	Ctr                         string `json:"ctr"`
	Acp                         string `json:"acp"`
	Cpm                         string `json:"cpm"`
	Like                        string `json:"like"`
	Comment                     string `json:"comment"`
	Collect                     string `json:"collect"`
	Follow                      string `json:"follow"`
	Share                       string `json:"share"`
	Interaction                 string `json:"interaction"`
	Cpi                         string `json:"cpi"`
	ActionButtonClick           string `json:"action_button_click"`
	ActionButtonCtr             string `json:"action_button_ctr"`
	Screenshot                  string `json:"screenshot"`
	PicSave                     string `json:"pic_save"`
	ExternalLeads               string `json:"external_leads"`
	ExternalLeadsCpl            string `json:"external_leads_cpl"`
	GoodsVisit                  string `json:"goods_visit"`
	GoodsVisitPrice             string `json:"goods_visit_price"`
	IUserNum                    string `json:"i_user_num"`
	TiUserNum                   string `json:"ti_user_num"`
	IUserPrice                  string `json:"i_user_price"`
	TiUserPrice                 string `json:"ti_user_price"`
	NoteID                      string `json:"note_id"`
	NoteTitle                   string `json:"note_title"`
	NoteImage                   string `json:"note_image"`
	NoteJumpURL                 string `json:"note_jump_url"`
	Placement                   string `json:"placement"`
	OptimizeTarget              string `json:"optimize_target"`
	PromotionTarget             string `json:"promotion_target"`
	BiddingStrategy             string `json:"bidding_strategy"`
	BuildType                   string `json:"build_type"`
	MarketingTarget             string `json:"marketing_target"`
	PageID                      string `json:"page_id"`
	ItemID                      string `json:"item_id"`
	LiveRedID                   string `json:"live_red_id"`
	CountryName                 string `json:"country_name"`
	Province                    string `json:"province"`
	City                        string `json:"city"`
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
	MessageFstReplyTimeAvg      string `json:"message_fst_reply_time_avg"`
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
	WordImpressionRankAll       string `json:"word_impression_rank_all"`
	WordImpressionRateAll       string `json:"word_impression_rate_all"`
	WordClickRankAll            string `json:"word_click_rank_all"`
	WordClickRateAll            string `json:"word_click_rate_all"`
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
	AppDownloadButtonClickCnt   string `json:"app_download_button_click_cnt"`
	AppDownloadButtonClickCtr   string `json:"app_download_button_click_ctr"`
	AppDownloadButtonClickCost  string `json:"app_download_button_click_cost"`
	AppActivateCnt              string `json:"app_activate_cnt"`
	AppActivateCost             string `json:"app_activate_cost"`
	AppActivateCtr              string `json:"app_activate_ctr"`
	AppRegisterCnt              string `json:"app_register_cnt"`
	AppRegisterCost             string `json:"app_register_cost"`
	AppRegisterCtr              string `json:"app_register_ctr"`
	FirstAppPayCnt              string `json:"first_app_pay_cnt"`
	FirstAppPayCost             string `json:"first_app_pay_cost"`
	FirstAppPayCtr              string `json:"first_app_pay_ctr"`
	CurrentAppPayCnt            string `json:"current_app_pay_cnt"`
	CurrentAppPayCost           string `json:"current_app_pay_cost"`
	AppKeyActionCnt             string `json:"app_key_action_cnt"`
	AppKeyActionCost            string `json:"app_key_action_cost"`
	AppKeyActionCtr             string `json:"app_key_action_ctr"`
	AppPayCnt7d                 string `json:"app_pay_cnt_7d"`
	AppPayCost7d                string `json:"app_pay_cost_7d"`
	AppPayAmount                string `json:"app_pay_amount"`
	AppPayRoi                   string `json:"app_pay_roi"`
	AppActivateAmount1d         string `json:"app_activate_amount_1d"`
	AppActivateAmount3d         string `json:"app_activate_amount_3d"`
	AppActivateAmount7d         string `json:"app_activate_amount_7d"`
	AppActivateAmount1dRoi      string `json:"app_activate_amount_1d_roi"`
	AppActivateAmount3dRoi      string `json:"app_activate_amount_3d_roi"`
	AppActivateAmount7dRoi      string `json:"app_activate_amount_7d_roi"`
	Retention1dCnt              string `json:"retention_1d_cnt"`
	Retention3dCnt              string `json:"retention_3d_cnt"`
	Retention7dCnt              string `json:"retention_7d_cnt"`
}
