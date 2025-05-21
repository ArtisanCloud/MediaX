package unit

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient/promote/unit/schema"
)

type JuGuangPromoteUnitClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangPromoteUnitClient {
	return &JuGuangPromoteUnitClient{
		BaseClient: c,
	}
}

// ## Create 创建推广单元
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=2723
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • campaign_id: 计划ID（必填）
//	  • unit_name: 单元名称（必填）
//	  • event_bid: 出价/目标成本（可选）
//	  • note_ids: 标的-笔记id（必填）
//	  • promotion_target: 推广标的（可选）
//	  • target_type: 定向类型（必填）
//	  • target_config: 定向配置（必填）
//	  • keyword_target_period: 关键词时间周期（可选）
//	  • keyword_target_action: 关键词行为类型（可选）
//	  • business_tree_name: 推广业务信息（可选）
//	  • spu_note_info: spu&笔记标的信息（可选）
//	  • keyword_with_bid: 关键词（可选）
//	  • substituted_user_id: 代投账号（可选）
//	  • keyword_gen_type: 单元选词方式（可选）
//	  • page_id: 落地页ID（可选）
//	  • landing_page_url: 落地页Url（可选）
//	  • unit_external_page_url: 外链Url（可选）
//	  • unit_landing_page_desc: 落地页表单描述（可选）
//	  • target_template_id: 定向包ID（可选）
//
// 返回值：
//
//	*schema.JuGuangPromoteUnitCreateRes 包含以下字段：
//	  • code: 响应状态码
//	  • message: 响应消息
//	  • data: 响应数据
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteUnitClient) Create(ctx context.Context, data *schema.JuGuangPromoteUnitCreateReq) (*schema.JuGuangPromoteUnitCreateRes, error) {
	result := &schema.JuGuangPromoteUnitCreateRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/unit/create", nil, data, nil, result)
	return result, err
}

// ## Update 编辑推广单元
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3151
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • unit_id: 单元ID（必填）
//	  • unit_name: 单元名称（可选）
//	  • event_bid: 出价/目标成本（可选）
//	  • note_ids: 标的-笔记id（可选）
//	  • target_type: 定向类型（可选）
//	  • target_config: 定向配置（可选）
//	  • keyword_target_period: 关键词时间周期（可选）
//	  • keyword_target_action: 关键词行为类型（可选）
//	  • business_tree_name: 推广业务信息（可选）
//	  • spu_note_info: spu&笔记标的信息（可选）
//	  • keyword_with_bid: 关键词（可选）
//	  • substituted_user_id: 代投账号（可选）
//	  • page_id: 落地页ID（可选）
//	  • landing_page_url: 落地页Url（可选）
//	  • unit_external_page_url: 外链Url（可选）
//	  • unit_landing_page_desc: 落地页表单描述（可选）
//	  • target_template_id: 定向包ID（可选）
//
// 返回值：
//
//	*schema.JuGuangPromoteUnitUpdateRes 包含以下字段：
//	  • code: 响应状态码
//	  • message: 响应消息
//	  • data: 响应数据
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteUnitClient) Update(ctx context.Context, data *schema.JuGuangPromoteUnitUpdateReq) (*schema.JuGuangPromoteUnitUpdateRes, error) {
	result := &schema.JuGuangPromoteUnitUpdateRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/unit/update", nil, data, nil, result)
	return result, err
}

// ## UpdateStatus 修改推广单元状态
//
// 接口文档参考：
// `https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3152`
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • unit_ids: 单元ID列表（必填，最多20个）
//	  • status: 状态枚举（必填）
//	    - 1: 开启
//	    - 2: 暂停
//	    - 3: 删除
//
// 返回值：
//
//	*schema.JuGuangPromoteUnitUpdateStatusRes 包含以下字段：
//	  • code: 响应状态码
//	  • message: 响应消息
//	  • success: 接口是否成功
//	  • data: 响应数据
//	    - unit_ids: 单元ID列表
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteUnitClient) UpdateStatus(ctx context.Context, data *schema.JuGuangPromoteUnitUpdateStatusReq) (*schema.JuGuangPromoteUnitUpdateStatusRes, error) {
	result := &schema.JuGuangPromoteUnitUpdateStatusRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/unit/update/status", nil, data, nil, result)
	return result, err
}

// ## List 获取推广单元列表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3044
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • campaign_id: 计划ID（可选）
//	  • unit_ids: 单元ID列表（可选，最多10个）
//	  • status: 投放状态（可选，1：投放中 2：暂停）
//	  • unit_name: 单元名称（可选）
//	  • start_date: 创建单元开始时间（可选）
//	  • end_date: 创建单元结束时间（可选）
//	  • page: 请求的页码数（可选，默认1）
//	  • page_size: 请求的每页行数（可选，默认20）
//
// 返回值：
//
//	*schema.JuGuangPromoteUnitListRes 包含以下字段：
//	  • code: 响应状态码
//	  • msg: 响应消息
//	  • success: 接口是否成功
//	  • data: 响应数据
//	    - total_count: 总数
//	    - unit_infos: 单元信息列表
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteUnitClient) List(ctx context.Context, data *schema.JuGuangPromoteUnitListReq) (*schema.JuGuangPromoteUnitListRes, error) {
	result := &schema.JuGuangPromoteUnitListRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/unit/list", nil, data, nil, result)
	return result, err
}
