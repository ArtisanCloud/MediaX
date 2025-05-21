package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangToolGetKeywordIndustryReq 表示获取关键词行业分类的请求参数
//
// 字段说明：
//
//	AdvertiserID - 广告主ID，必填项
type JuGuangToolGetKeywordIndustryReq struct {
	AdvertiserID int64 `json:"advertiser_id"`
}

// DictDto 表示行业分类数据
//
// 字段说明：
//
//	TaxonomyID - 行业ID
//	FullPathName - 全路径名
//	TaxonomyName - 行业名称
//	TaxonomyLevel - 层级数
//	Children - 子层级
type DictDto struct {
	TaxonomyID    string    `json:"taxonomy_id"`
	FullPathName  string    `json:"full_path_name"`
	TaxonomyName  string    `json:"taxonomy_name"`
	TaxonomyLevel int       `json:"taxonomy_level"`
	Children      []DictDto `json:"children"`
}

// GetKeywordIndustryData 表示获取关键词行业分类的返回数据
//
// 字段说明：
//
//	AllIndustryTaxonomys - 所有可选的行业类目
//	AdsIndustryTaxonomyDictDto - 行业细节
type GetKeywordIndustryData struct {
	AllIndustryTaxonomys       string    `json:"all_industry_taxonomys"`
	AdsIndustryTaxonomyDictDto []DictDto `json:"ads_industry_taxonomy_dict_dto"`
}

// JuGuangToolGetKeywordIndustryRes 表示获取关键词行业分类的返回结果
//
// 字段说明：
//
//	JuGuangRes - 聚光平台基础响应结构
//	Data - 业务数据主体，包含行业分类数据
type JuGuangToolGetKeywordIndustryRes struct {
	response.JuGuangRes
	Data GetKeywordIndustryData `json:"data"`
}
