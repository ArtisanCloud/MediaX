package creativity

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient/promote/creativity/schema"
)

// JuGuangPromoteCreativityClient 创意客户端
type JuGuangPromoteCreativityClient struct {
	*kernel.BaseClient
}

// NewClient 新建创意客户端
func NewClient(c *kernel.BaseClient) *JuGuangPromoteCreativityClient {
	return &JuGuangPromoteCreativityClient{
		BaseClient: c,
	}
}

// ## CreateNote 创建笔记创意
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3153
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 创建笔记创意请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • unit_id: 单元ID
//	  • creativity_name: 创意名称
//	  • note_id: 笔记ID
//	  • click_urls: 点击链接列表（可选）
//	  • expo_urls: 曝光链接列表（可选）
//	  • custom_mask: 自选封面模式（可选）
//	  • custom_title: 是否自提标题（可选）
//	  • title_fills: 自提标题列表（可选）
//	  • mask_gen: 是否开启自动优化封面（可选）
//	  • title_gen: 是否开启自动优化标题（可选）
//	  • conversion_type: 组件类型
//	  • jump_url: 落地页/外链URL（可选）
//	  • landing_page_type: 落地页链接类型（可选）
//	  • bar_content: 按钮文案内容（可选）
//	  • conversion_component_types: 组件位置列表（可选）
//	  • comment: 置顶评论文案（可选）
//	  • qual_info: 资质信息（必填）
//
// 返回值：
//
//	*schema.JuGuangPromoteCreativityCreateNoteRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCreativityClient) CreateNote(ctx context.Context, data *schema.JuGuangPromoteCreativityCreateNoteReq) (*schema.JuGuangPromoteCreativityCreateNoteRes, error) {
	result := &schema.JuGuangPromoteCreativityCreateNoteRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/creativity/note/create", nil, data, nil, result)
	return result, err
}

// ## CreateLandingPage 创建落地页创意
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3154
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 创建落地页创意请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • unit_id: 单元ID
//	  • creativity_name: 创意名称
//	  • page_creativity_infos: 落地页创意信息
//	  • qual_info: 资质信息
//
// 返回值：
//
//	*schema.JuGuangPromoteCreativityCreateLandingPageRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据，包含创意ID
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCreativityClient) CreateLandingPage(ctx context.Context, data *schema.JuGuangPromoteCreativityCreateLandingPageReq) (*schema.JuGuangPromoteCreativityCreateLandingPageRes, error) {
	result := &schema.JuGuangPromoteCreativityCreateLandingPageRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/creativity/page/create", nil, data, nil, result)
	return result, err
}

// ## CreateProgrammaticPage 创建程序化创意
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3155
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 创建程序化创意请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • unit_id: 单元ID（必填）
//	  • creativity_name: 创意名称（可选）
//	  • h5_material_info: 程序化创意信息（可选）
//	  • qual_info: 资质信息（必填）
//
// 返回值：
//
//	*schema.JuGuangPromoteCreativityCreateProgrammaticPageRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据，包含创意ID
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCreativityClient) CreateProgrammaticPage(ctx context.Context, data *schema.JuGuangPromoteCreativityCreateProgrammaticPageReq) (*schema.JuGuangPromoteCreativityCreateProgrammaticPageRes, error) {
	result := &schema.JuGuangPromoteCreativityCreateProgrammaticPageRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/creativity/programmatic/page/create", nil, data, nil, result)
	return result, err
}

// ## Update 编辑创意
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3156
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 编辑创意请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • creativity_id: 创意ID（必填）
//	  • creativity_name: 创意名称（可选）
//	  • click_urls: 点击链接列表（可选）
//	  • expo_urls: 曝光链接列表（可选）
//	  • mask_perfer: 是否开启封面优选（可选）
//	  • title_mask_perfer: 是否开启标题优选（可选）
//	  • jump_url: 跳转链接（可选）
//	  • bar_content: 文案内容（可选）
//	  • conversion_component_types: 转化组件类型列表（可选）
//	  • comment: 评论区文案（可选）
//	  • poi_id: POI ID（可选）
//	  • poi_jump_type: POI组件跳转类型（可选）
//	  • monitor_company: 监测公司（可选）
//	  • monitor_params: 监测参数配置（可选）
//
// 返回值：
//
//	*schema.JuGuangPromoteCreativityUpdateRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCreativityClient) Update(ctx context.Context, data *schema.JuGuangPromoteCreativityUpdateReq) (*schema.JuGuangPromoteCreativityUpdateRes, error) {
	result := &schema.JuGuangPromoteCreativityUpdateRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/creativity/update", nil, data, nil, result)
	return result, err
}

// ## UpdateStatus 修改创意状态
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3157
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 修改创意状态请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • creativity_ids: 创意ID列表（必填，上限20个）
//	  • action_type: 操作类型（必填，1-开启，2-暂停，3-删除）
//
// 返回值：
//
//	*schema.JuGuangPromoteCreativityUpdateStatusRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据，包含所有修改状态成功的创意ID
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCreativityClient) UpdateStatus(ctx context.Context, data *schema.JuGuangPromoteCreativityUpdateStatusReq) (*schema.JuGuangPromoteCreativityUpdateStatusRes, error) {
	result := &schema.JuGuangPromoteCreativityUpdateStatusRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/creativitystatus/update", nil, data, nil, result)
	return result, err
}

// ## Search 创意查询
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3158
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 创意查询请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • campaign_id: 计划ID（可选）
//	  • unit_id: 单元ID（可选）
//	  • creativity_ids: 创意ID列表（可选，最多20个）
//	  • status: 创意状态（可选）
//	  • start_time: 开始时间（可选）
//	  • end_time: 结束时间（可选）
//	  • page_index: 页码（可选，默认1）
//	  • page_size: 每页数量（可选，默认20）
//
// 返回值：
//
//	*schema.JuGuangPromoteCreativitySearchRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据，包含创意列表和分页信息
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangPromoteCreativityClient) Search(ctx context.Context, data *schema.JuGuangPromoteCreativitySearchReq) (*schema.JuGuangPromoteCreativitySearchRes, error) {
	result := &schema.JuGuangPromoteCreativitySearchRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/creativity/search", nil, data, nil, result)
	return result, err
}
