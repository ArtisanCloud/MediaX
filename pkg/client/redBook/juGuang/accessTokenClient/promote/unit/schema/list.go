package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangPromoteUnitListReq 表示获取小红书聚光平台推广单元列表的请求参数
type JuGuangPromoteUnitListReq struct {
	AdvertiserID int64   `json:"advertiser_id"`         // 广告主ID
	CampaignID   *int64  `json:"campaign_id,omitempty"` // 计划ID
	UnitIDs      []int64 `json:"unit_ids,omitempty"`    // 单元ID列表
	Status       *int    `json:"status,omitempty"`      // 投放状态
	UnitName     *string `json:"unit_name,omitempty"`   // 单元名称
	StartDate    *string `json:"start_date,omitempty"`  // 创建开始时间
	EndDate      *string `json:"end_date,omitempty"`    // 创建结束时间
	Page         *int    `json:"page,omitempty"`        // 页码
	PageSize     *int    `json:"page_size,omitempty"`   // 每页数量
}

// UnitInfo 表示推广单元信息
type UnitInfo struct {
	ID                  int64    `json:"id"`                     // 单元ID
	CampaignID          int64    `json:"campaign_id"`            // 计划ID
	Name                string   `json:"name"`                   // 单元名称
	Enable              int      `json:"enable"`                 // 投放状态
	UnitFilterState     int      `json:"unit_filter_state"`      // 单元过滤状态
	EventBid            int      `json:"event_bid"`              // 出价
	TargetType          int      `json:"target_type"`            // 定向类型
	ItemIDs             []string `json:"item_ids"`               // 商品ID
	NoteIDs             []string `json:"note_ids"`               // 笔记ID
	LiveUserID          string   `json:"live_user_id"`           // 直播用户ID
	PageID              string   `json:"page_id"`                // 落地页ID
	LandingPageURL      string   `json:"landing_page_url"`       // 落地页URL
	UnitExternalPageURL string   `json:"unit_external_page_url"` // 外链URL
	LandingPageType     int      `json:"landing_page_type"`      // 落地页类型
	TargetPosition      int      `json:"target_position"`        // 抢占位置
	TargetGoal          int      `json:"target_goal"`            // 抢占目标
	WordTagName         string   `json:"word_tag_name"`          // 词包名称
	ProportionGoal      float64  `json:"proportion_goal"`        // 占比目标
	BusinessTreeName    string   `json:"business_tree_name"`     // 推广业务信息
	UnitLandingPageDesc []string `json:"unit_landing_page_desc"` // 落地页描述
	KeywordTargetPeriod int      `json:"keyword_target_period"`  // 关键词定向周期
	KeywordTargetAction []int    `json:"keyword_target_action"`  // 关键词定向行为
	SubstitutedUserID   string   `json:"substituted_user_id"`    // 代投账号ID
	CreateTime          string   `json:"create_time"`            // 创建时间
	UpdateTime          string   `json:"update_time"`            // 更新时间
}

type UnitListData struct {
	TotalCount int        `json:"total_count"` // 总数量
	UnitInfos  []UnitInfo `json:"unit_infos"`  // 单元信息列表
}

// JuGuangPromoteUnitListRes 表示获取小红书聚光平台推广单元列表的响应
type JuGuangPromoteUnitListRes struct {
	response.RedBookAccessTokenRes
	// Data 单元列表数据
	Data UnitListData `json:"data"`
}
