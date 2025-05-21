package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// ## JuGuangNoteQueryTemplateReq 查询定向包请求参数
// 字段说明：
//
//	advertiser_id - 广告主id（必填）
//	target_template_id - 定向包id（可选）
//	marketing_target - 营销诉求（可选）
//	placement - 广告类型（可选）
//	page - 页码（可选，默认1）
//	page_size - 每页行数（可选，默认20，最大100）
type JuGuangNoteQueryTemplateReq struct {
	AdvertiserID     int64  `json:"advertiser_id"`                // 广告主id
	TargetTemplateID *int64 `json:"target_template_id,omitempty"` // 定向包id
	MarketingTarget  *int   `json:"marketing_target,omitempty"`   // 营销诉求
	Placement        *int   `json:"placement,omitempty"`          // 广告类型
	Page             *int   `json:"page,omitempty"`               // 页码
	PageSize         *int   `json:"page_size,omitempty"`          // 每页行数
}

// ## JuGuangNoteQueryTemplateRes 查询定向包响应结构体
// 字段说明：
//
//	code - 响应码
//	msg - 响应消息
//	success - 是否成功
//	data - 业务数据
//	request_id - 请求id
type JuGuangNoteQueryTemplateRes struct {
	response.RedBookAccessTokenRes
	Data *JuGuangNoteQueryTemplateData `json:"data"` // 业务数据
}

// ## JuGuangNoteQueryTemplateData 定向包数据结构体
type JuGuangNoteQueryTemplateData struct {
	ID              int64                `json:"id"`               // 定向包ID
	Name            string               `json:"name"`             // 定向包名称
	Desc            string               `json:"desc"`             // 定向包描述
	TargetType      int                  `json:"target_type"`      // 定向类型
	MarketingTarget int                  `json:"marketing_target"` // 营销诉求
	Placement       int                  `json:"placement"`        // 广告类型
	DeliveryMode    int                  `json:"delivery_mode"`    // 投放模式
	UnitList        []int64              `json:"unit_list"`        // 已关联单元ID列表
	TargetConfig    *JuGuangTargetConfig `json:"target_config"`    // 定向配置
}

// ## JuGuangTargetConfig 定向配置结构体
type JuGuangTargetConfig struct {
	TargetGender           string                  `json:"target_gender"`             // 定向性别
	TargetAge              string                  `json:"target_age"`                // 定向年龄
	TargetDevice           string                  `json:"target_device"`             // 定向设备
	TargetDevicePrice      string                  `json:"target_device_price"`       // 手机价格
	TargetCity             string                  `json:"target_city"`               // 城市定向
	TargetAreaCode         string                  `json:"target_area_code"`          // 城市编码
	IndustryInterestTarget *IndustryInterestTarget `json:"industry_interest_target"`  // 行业兴趣
	CrowdTarget            *CrowdTarget            `json:"crowd_target"`              // 人群包
	Keywords               []string                `json:"keywords"`                  // 关键词定向
	InterestKeywords       []string                `json:"interest_keywords"`         // 关键词兴趣定向
	IntelligentExpansion   int                     `json:"intelligent_expansion"`     // 智能扩量
	SearchTargetCityIntent string                  `json:"search_target_city_intent"` // 搜索意图城市定向
	ReverseTargetCrowd     []string                `json:"reverse_target_crowd"`      // 排除特定人群
}

// ## IndustryInterestTarget 行业兴趣结构体
type IndustryInterestTarget struct {
	ContentInterests  []CodeNamePair `json:"content_interests"`  // 内容兴趣
	ShoppingInterests []CodeNamePair `json:"shopping_interests"` // 购物兴趣
}

// ## CodeNamePair 代码-名称对结构体
type CodeNamePair struct {
	Code     string         `json:"code"`     // code
	Name     string         `json:"name"`     // 名称
	Children []CodeNamePair `json:"children"` // 子节点
}

// ## CrowdTarget 人群包结构体
type CrowdTarget struct {
	CrowdPkg []CrowdPackageVO `json:"crowd_pkg"` // 人群包列表
}

// ## CrowdPackageVO 人群包VO结构体
type CrowdPackageVO struct {
	Value string `json:"value"` // 人群包id
	Name  string `json:"name"`  // 人群包名称
	Type  string `json:"type"`  // 人群包类型
}
