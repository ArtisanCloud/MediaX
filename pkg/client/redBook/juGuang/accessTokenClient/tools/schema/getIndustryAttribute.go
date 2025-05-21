package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

type JuGuangToolGetIndustryAttributeReq struct {
	// 行业ID，用于指定要查询的行业属性。
	IndustryID string `json:"industry_id"`
}

// AttributeDto 表示行业属性数据
//
// 字段说明：
//
//	TaxonomyID - 分类ID，表示行业分类的唯一标识
//	TaxonomyAttributeName - 分类属性名称，表示行业分类的具体属性名称
//	TaxonomyLevel - 分类层级，表示行业分类的层级深度
type AttributeDto struct {
	TaxonomyID            string `json:"taxonomy_id"`
	TaxonomyAttributeName string `json:"taxonomy_attribute_name"`
	TaxonomyLevel         int    `json:"taxonomy_level"`
}

// GetIndustryAttributeData 表示获取行业属性数据的返回结构
//
// 字段说明：
//
//	TaxonomyAttributeDtos - 行业属性数据列表，包含多个AttributeDto对象
type GetIndustryAttributeData struct {
	TaxonomyAttributeDtos []AttributeDto `json:"taxonomy_attribute_dtos"`
}

// JuGuangToolGetIndustryAttributeRes 表示获取行业属性接口的返回结果
//
// 字段说明：
//
//	JuGuangRes - 聚光平台基础响应结构
//	Data - 业务数据主体，包含行业属性数据
type JuGuangToolGetIndustryAttributeRes struct {
	response.JuGuangRes
	Data GetIndustryAttributeData `json:"data"`
}
