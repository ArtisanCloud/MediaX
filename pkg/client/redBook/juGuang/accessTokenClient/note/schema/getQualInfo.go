package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// ## JuGuangNoteGetQualInfoReq 获取资质信息请求参数
//
// 字段说明：
//
//	AdvertiserId - 广告主ID (必填)
//	Page - 分页信息 (可选)
//	  • PageIndex - 页码 (默认1)
//	  • PageSize - 每页查询量级 (默认20)
type JuGuangNoteGetQualInfoReq struct {
	AdvertiserId int `json:"advertiser_id"`
	Page         struct {
		PageIndex *int `json:"page_index,omitempty"`
		PageSize  *int `json:"page_size,omitempty"`
	} `json:"page"`
}

// ## JuGuangNoteGetQualInfoRes 获取资质信息返回结果
//
// 字段说明：
//
//	Code - 返回码
//	Msg - 返回信息
//	Success - 接口是否成功
//	Data - 业务数据主体
//	  • QualInfos - 资质信息列表
//	  • ProductQualInfos - 行业资质信息列表
//	  • BrandQualInfos - 产品资质信息列表
type JuGuangNoteGetQualInfoRes struct {
	response.JuGuangRes
	Data struct {
		QualInfos        []QualInfo        `json:"qual_infos"`
		ProductQualInfos []ProductQualInfo `json:"product_qual_infos"`
		BrandQualInfos   []BrandQualInfo   `json:"brand_qual_infos"`
	} `json:"data"`
}

// QualInfo 资质信息
//
// 字段说明：
//
//	ApplyId - 资质ID
//	TradeTypeFirstCode - 一级行业编码
//	TradeTypeFirstName - 一级行业名称
//	TradeTypeSecondCode - 二级行业编码
//	TradeTypeSecondName - 二级行业名称
type QualInfo struct {
	ApplyId             string `json:"apply_id"`
	TradeTypeFirstCode  string `json:"trade_type_first_code"`
	TradeTypeFirstName  string `json:"trade_type_first_name"`
	TradeTypeSecondCode string `json:"trade_type_second_code"`
	TradeTypeSecondName string `json:"trade_type_second_name"`
}

// ProductQualInfo 行业资质信息
//
// 字段说明：
//
//	ProductQualId - 产品资质ID
//	ProductName - 产品资质名称
type ProductQualInfo struct {
	ProductQualId int    `json:"product_qual_id"`
	ProductName   string `json:"product_name"`
}

// BrandQualInfo 产品资质信息
//
// 字段说明：
//
//	BrandQualId - 品牌资质ID
//	QualTypeName - 资质名称
//	UserRemark - 备注
type BrandQualInfo struct {
	BrandQualId  int    `json:"brand_qual_id"`
	QualTypeName string `json:"qual_type_name"`
	UserRemark   string `json:"user_remark"`
}
