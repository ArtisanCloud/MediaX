package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangPromoteUnitCreateReq 创建推广单元请求参数
type JuGuangPromoteUnitCreateReq struct {
	AdvertiserID        int64        `json:"advertiser_id"`                    // 广告主ID
	CampaignID          int64        `json:"campaign_id"`                      // 计划ID
	UnitName            string       `json:"unit_name"`                        // 单元名称
	EventBid            *int         `json:"event_bid,omitempty"`              // 出价/目标成本
	NoteIDs             []string     `json:"note_ids"`                         // 笔记ID列表
	PromotionTarget     *int         `json:"promotion_target,omitempty"`       // 推广标的
	TargetType          int          `json:"target_type"`                      // 定向类型
	TargetConfig        TargetConfig `json:"target_config"`                    // 定向配置
	KeywordTargetPeriod *int         `json:"keyword_target_period,omitempty"`  // 关键词时间周期
	KeywordTargetAction []int        `json:"keyword_target_action,omitempty"`  // 关键词行为类型
	SubstitutedUserID   *string      `json:"substituted_user_id,omitempty"`    // 代投账号
	KeywordGenType      *int         `json:"keyword_gen_type,omitempty"`       // 单元选词方式
	PageID              *string      `json:"page_id,omitempty"`                // 落地页ID
	LandingPageURL      *string      `json:"landing_page_url,omitempty"`       // 落地页URL
	UnitExternalPageURL *string      `json:"unit_external_page_url,omitempty"` // 外链URL
	UnitLandingPageDesc []string     `json:"unit_landing_page_desc,omitempty"` // 落地页表单描述
	TargetTemplateID    *int64       `json:"target_template_id,omitempty"`     // 定向包ID
}

// TargetConfig 定向配置
type TargetConfig struct {
	TargetGender           *string                 `json:"target_gender,omitempty"`             // 性别
	TargetAge              *string                 `json:"target_age,omitempty"`                // 年龄
	TargetCity             string                  `json:"target_city"`                         // 城市
	TargetDevice           *string                 `json:"target_device,omitempty"`             // 设备
	IndustryInterestTarget *IndustryInterestTarget `json:"industry_interest_target,omitempty"`  // 行业兴趣
	CrowdTarget            *CrowdTarget            `json:"crowd_target,omitempty"`              // 人群包
	InterestKeywords       []string                `json:"interest_keywords,omitempty"`         // 关键词兴趣
	Keywords               []string                `json:"keywords,omitempty"`                  // 关键词行为
	IntelligentExpansion   *int                    `json:"intelligent_expansion,omitempty"`     // 智能扩量
	SearchTargetCityIntent *string                 `json:"search_target_city_intent,omitempty"` // 搜索地域意图
	ReverseTargetCrowd     []string                `json:"reverse_target_crowd,omitempty"`      // 排除特定人群
}

// IndustryInterestTarget 行业兴趣定向
type IndustryInterestTarget struct {
	ContentInterests  []CodeNamePair `json:"content_interests,omitempty"`  // 行业阅读兴趣
	ShoppingInterests []CodeNamePair `json:"shopping_interests,omitempty"` // 行业购物兴趣
}

// CodeNamePair 代码名称对
type CodeNamePair struct {
	Code     string         `json:"code"`               // 代码
	Name     string         `json:"name"`               // 名称
	Children []CodeNamePair `json:"children,omitempty"` // 子节点
}

// CrowdTarget 人群包定向
type CrowdTarget struct {
	CrowdPkg      []CrowdPackageVO `json:"crowd_pkg,omitempty"`      // 人群包
	DmpPermission bool             `json:"dmp_permission,omitempty"` // DMP权限
}

// CrowdPackageVO 人群包VO
type CrowdPackageVO struct {
	Value string `json:"value"` // 人群包ID
	Name  string `json:"name"`  // 人群包名称
}

type UnitCreateData struct {
	UnitID int64 `json:"unit_id"` // 单元ID
}

// JuGuangPromoteUnitCreateRes 创建推广单元响应参数
type JuGuangPromoteUnitCreateRes struct {
	response.RedBookAccessTokenRes
	// UnitCreateData 创建单元数据
	Data UnitCreateData `json:"data"`
}
