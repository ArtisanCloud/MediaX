package campaign

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient/promote/campaign/schema"
)

// ## JuGuangPromoteCampaignClient 推广计划客户端
type JuGuangPromoteCampaignClient struct {
	*kernel.BaseClient
}

// ## NewClient 新建推广计划客户端
func NewClient(c *kernel.BaseClient) *JuGuangPromoteCampaignClient {
	return &JuGuangPromoteCampaignClient{
		BaseClient: c,
	}
}

// ## Create 创建推广计划
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=2722
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • marketing_target: 营销目标（必填，4：产品种草，9：客资收集，16：应用推广）
//	  • campaign_name: 计划名称（必填，长度不超过50个字符）
//	  • placement: 广告类型（必填，1：信息流，2：搜索推广，4：通投/全站智投，7：视频内流）
//	  • promotion_target: 推广标的类型（必填，1：笔记）
//	  • enable: 计划创建后默认开启状态（可选，1-开启，0-不开启，默认开启）
//	  • time_type: 推广时间类型（必填，0:长期投放，1:自定义设置开始结束时间）
//	  • start_time: 推广开始时间（可选，格式 yyyy-MM-dd）
//	  • expire_time: 推广结束时间（可选，格式 yyyy-MM-dd）
//	  • time_period_type: 推广时段类型（必填，0: 全时段，1:自定义时间段）
//	  • bidding_strategy: 出价方式（必填，2: 手动出价，3: 自动出价-成本自动控制，4:自动出价，7: 自动出价-成本手动控制）
//	  • limit_day_budget: 预算类型（必填，0：不限预算，1：指定预算）
//	  • campaign_day_budget: 计划日预算（可选，单位分，范围 [10000~99999900)）
//	  • optimize_target: 推广目标（必填，0-点击量，1-互动量，2-展示次数，3-表单提交量等）
//
// 返回值：
//
//	*schema.JuGuangPromoteCampaignCreateRes 包含以下字段：
//	  • code: 返回码
//	  • message: 返回信息
//	  • data: 创建成功的计划ID
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCampaignClient) Create(ctx context.Context, data *schema.JuGuangPromoteCampaignCreateReq) (*schema.JuGuangPromoteCampaignCreateRes, error) {
	result := &schema.JuGuangPromoteCampaignCreateRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/campaign/create", nil, data, nil, result)
	return result, err
}

// ## Update 编辑推广计划
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3148
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • campaign_id: 计划ID（必填）
//	  • campaign_name: 计划名称（可选，长度不超过50个字符）
//	  • enable: 计划状态（可选，1-开启，0-关闭）
//	  • time_type: 推广时间类型（可选，0:长期投放，1:自定义设置开始结束时间）
//	  • start_time: 推广开始时间（可选，格式 yyyy-MM-dd）
//	  • expire_time: 推广结束时间（可选，格式 yyyy-MM-dd）
//	  • time_period_type: 推广时段类型（可选，0: 全时段，1:自定义时间段）
//	  • limit_day_budget: 预算类型（可选，0：不限预算，1：指定预算）
//	  • campaign_day_budget: 计划日预算（可选，单位分，范围 [10000~99999900)）
//
// 返回值：
//
//	*schema.JuGuangPromoteCampaignUpdateRes 包含以下字段：
//	  • code: 返回码
//	  • message: 返回信息
//	  • data: 更新成功的计划ID
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCampaignClient) Update(ctx context.Context, data *schema.JuGuangPromoteCampaignUpdateReq) (*schema.JuGuangPromoteCampaignUpdateRes, error) {
	result := &schema.JuGuangPromoteCampaignUpdateRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/campaign/update", nil, data, nil, result)
	return result, err
}

// ## UpdateStatus 修改推广计划状态
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3149
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • campaign_ids: 计划ID列表（必填，至少传一个，最多20个）
//	  • action_type: 操作类型（必填，1：开启，2：暂停，3：删除）
//
// 返回值：
//
//	*schema.JuGuangPromoteCampaignUpdateStatusRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 修改状态成功的计划ID列表
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCampaignClient) UpdateStatus(ctx context.Context, data *schema.JuGuangPromoteCampaignUpdateStatusReq) (*schema.JuGuangPromoteCampaignUpdateStatusRes, error) {
	result := &schema.JuGuangPromoteCampaignUpdateStatusRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/campaign/status/update", nil, data, nil, result)
	return result, err
}

// ## List 查询推广计划列表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3150
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • campaign_ids: 计划ID列表（可选，最多20个）
//	  • start_time: 开始时间（可选，格式 yyyy-MM-dd）
//	  • expire_time: 结束时间（可选，格式 yyyy-MM-dd）
//	  • status: 计划状态（可选，1-有效 2-暂停 3-已删除）
//	  • page: 分页参数（可选）
//	    • page_index: 页码（默认1）
//	    • page_size: 每页数量（默认20，最大100）
//
// 返回值：
//
//	*schema.JuGuangPromoteCampaignListRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 查询结果，包含以下字段：
//	    • page: 分页信息
//	      • page_index: 当前页码
//	      • total_count: 总数量
//	    • base_campaign_dtos: 计划列表，包含以下字段：
//	      • campaign_id: 计划ID
//	      • campaign_name: 计划名称
//	      • campaign_filter_state: 计划状态
//	      • campaign_create_time: 计划创建时间
//	      • marketing_target: 营销目标
//	      • placement: 广告类型
//	      • optimize_target: 优化目标
//	      • promotion_target: 投放标的
//	      • bidding_strategy: 出价策略
//	      • limit_day_budget: 预算类型
//	      • campaign_day_budget: 计划日预算
//	      • start_time: 推广开始时间
//	      • expire_time: 推广结束时间
//	      • time_period_type: 推广时间段类型
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCampaignClient) List(ctx context.Context, data *schema.JuGuangPromoteCampaignListReq) (*schema.JuGuangPromoteCampaignListRes, error) {
	result := &schema.JuGuangPromoteCampaignListRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/campaign/list", nil, data, nil, result)
	return result, err
}
