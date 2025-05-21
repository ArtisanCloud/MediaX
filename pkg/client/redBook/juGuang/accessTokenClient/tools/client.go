package tools

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient/tools/schema"
)

type JuGuangToolClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangToolClient {
	return &JuGuangToolClient{
		BaseClient: c,
	}
}

// ## GetOperationRecord 获取历史操作记录
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=4096
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 品牌方ID（必填）
//	  • start_time: 开始时间（必填，格式yyyy-mm-dd）
//	  • end_time: 结束时间（必填，格式yyyy-mm-dd）
//	  • page: 页码（可选，默认1）
//	  • page_size: 每页大小（可选，默认20，最大50）
//	  • level: 操作层级（可选，默认-1表示全部）
//	  • opt_type: 操作类型（可选，默认-1表示全部）
//
// 返回值：
//
//	*schema.JuGuangToolGetOperationRecordRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 操作记录列表
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) GetOperationRecord(ctx context.Context, data *schema.JuGuangToolGetOperationRecordReq) (*schema.JuGuangToolGetOperationRecordRes, error) {
	result := &schema.JuGuangToolGetOperationRecordRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/history/list", nil, data, nil, result)
	return result, err
}

// ## KeywordRecommend 定向推词-以词推词
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3168
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • keyword: 关键词（必填）
//	  • scene: 使用场景（可选，如搜索场景、推荐场景等）
//	  • count: 返回结果数量（可选，默认10）
//
// 返回值：
//
//	*schema.JuGuangToolKeywordRecommendRes 包含以下字段：
//	  • code: 返回码
//	  • message: 返回信息
//	  • data: 推荐关键词列表
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) KeywordRecommend(ctx context.Context, data *schema.JuGuangToolKeywordRecommendReq) (*schema.JuGuangToolKeywordRecommendRes, error) {
	result := &schema.JuGuangToolKeywordRecommendRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/keyword/common/recommend", nil, data, nil, result)
	return result, err
}

// ## GetIndustryAttribute 获取行业类目属性
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3167
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • taxonomy_id: 行业类目ID（必填）
//
// 返回值：
//
//	*schema.JuGuangToolGetIndustryAttributeRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 行业类目属性列表
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) GetIndustryAttribute(ctx context.Context, data *schema.JuGuangToolGetIndustryAttributeReq) (*schema.JuGuangToolGetIndustryAttributeRes, error) {
	result := &schema.JuGuangToolGetIndustryAttributeRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/keyword/industry/taxonomy/attribute", nil, nil, nil, result)
	return result, err
}

// ## GetKeywordIndustry 获取行业类目信息
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3166
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//
// 返回值：
//
//	*schema.JuGuangToolGetKeywordIndustryRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 行业类目信息
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) GetKeywordIndustry(ctx context.Context, data *schema.JuGuangToolGetKeywordIndustryReq) (*schema.JuGuangToolGetKeywordIndustryRes, error) {
	result := &schema.JuGuangToolGetKeywordIndustryRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/keyword/industry/taxonomy", nil, nil, nil, result)
	return result, err
}

// ## GetWordBagList 获取词包推荐列表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3165
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • name: 搜索名称（可选）
//	  • category: 词包类别（可选）
//	  • page_num: 页码（可选，默认1）
//	  • page_size: 每页大小（可选，默认10）
//	  • start_time: 开始时间（可选，格式yyyy-mm-dd）
//	  • end_time: 结束时间（可选，格式yyyy-mm-dd）
//
// 返回值：
//
//	*schema.JuGuangToolGetWordBagListRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 词包列表信息
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) GetWordBagList(ctx context.Context, data *schema.JuGuangToolGetWordBagListReq) (*schema.JuGuangToolGetWordBagListRes, error) {
	result := &schema.JuGuangToolGetWordBagListRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/keyword/word/bag/list", nil, data, nil, result)
	return result, err
}

// ## CrowdEstimate 人群预估
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3164
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • target_type: 定向类型（必填，1-通投,2-智能定向,3-高级定向）
//	  • placement: 广告类型（必填，1-信息流,2-搜索推广,4-全站智投,7-视频内流）
//	  • optimize_target: 推广目标（必填，如点击量、互动量等）
//	  • marketing_target: 营销目标（必填，如商品销量、笔记种草等）
//	  • target_config: 定向配置（必填，包含年龄、性别、地域等）
//
// 返回值：
//
//	*schema.JuGuangToolCrowdEstimateRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 预估结果，包含：
//	    - crowd_scope: 受众范围（1-偏窄,2-合适,3-偏广）
//	    - crowd_num: 受众人数
//	    - raw_crowd_num: 受众人数数值
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) CrowdEstimate(ctx context.Context, data *schema.JuGuangToolCrowdEstimateReq) (*schema.JuGuangToolCrowdEstimateRes, error) {
	result := &schema.JuGuangToolCrowdEstimateRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/crowd/estimate", nil, data, nil, result)
	return result, err
}

// ## GetKeywordMatch 获取关键词匹配词库信息
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3163
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • keywords: 关键词列表（必填，最多150个）
//
// 返回值：
//
//	*schema.JuGuangToolGetKeywordMatchRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 匹配信息，包含：
//	    - match_distinct_count: 匹配个数
//	    - match_infos: 匹配信息列表
//	      • keyword: 关键词
//	      • in_thesaurus: 是否匹配词库
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) GetKeywordMatch(ctx context.Context, data *schema.JuGuangToolGetKeywordMatchReq) (*schema.JuGuangToolGetKeywordMatchRes, error) {
	result := &schema.JuGuangToolGetKeywordMatchRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/keyword/match", nil, data, nil, result)
	return result, err
}

// ## GetKeywordRecommend 获取推荐关键词信息
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3160
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • keyword: 关键词（可选）
//	  • note_ids: 笔记ID列表（可选）
//
// 返回值：
//
//	*schema.JuGuangToolKeywordRecommendRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 推荐关键词列表，包含：
//	    - target_word: 推荐关键词
//	    - recommend_reason: 推荐理由
//	    - cover_num: 覆盖人数
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) GetKeywordRecommend(ctx context.Context, data *schema.JuGuangToolKeywordRecommendReq) (*schema.JuGuangToolKeywordRecommendRes, error) {
	result := &schema.JuGuangToolKeywordRecommendRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/keyword/recommend", nil, data, nil, result)
	return result, err
}

// ## GetTargetInfo 获取定向信息
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3159
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • marketing_target: 营销目标（必填，如商品销量、产品种草等）
//
// 返回值：
//
//	*schema.JuGuangToolGetTargetInfoRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 定向信息，包含：
//	    - industry_interest_target: 行业兴趣
//	    - content_interests: 内容兴趣
//	    - shopping_interests: 购物兴趣
//	    - crowd_target: 人群包
//	    - gender_targets: 性别
//	    - age_targets: 年龄
//	    - area_targets: 地域
//	    - device_targets: 设备
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) GetTargetInfo(ctx context.Context, data *schema.JuGuangToolGetTargetInfoReq) (*schema.JuGuangToolGetTargetInfoRes, error) {
	result := &schema.JuGuangToolGetTargetInfoRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/target/get_available_target_info", nil, data, nil, result)
	return result, err
}

// ## CheckDuplicatedName 计划单元名称重复性校验
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3113
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • type: 查询类型（必填，1：计划，2：单元）
//	  • name: 计划名称/单元名称集合（必填，单次查询上限100）
//
// 返回值：
//
//	*schema.JuGuangToolCheckDuplicatedNameRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • data: 查询结果，包含：
//	    - type: 查询类型
//	    - check_result: 查询结果，key为查询名称，value为是否重复
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangToolClient) CheckDuplicatedName(ctx context.Context, data *schema.JuGuangToolCheckDuplicatedNameReq) (*schema.JuGuangToolCheckDuplicatedNameRes, error) {
	result := &schema.JuGuangToolCheckDuplicatedNameRes{}
	_, err := c.BaseClient.HttpPost(ctx, " /api/open/jg/data/check/name/dup", nil, data, nil, result)
	return result, err
}
