package note

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient/note/schema"
)

// JuGuangNoteClient 表示巨量引擎直链管理 API 客户端
type JuGuangNoteClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 JuGuangNoteClient 实例
func NewClient(c *kernel.BaseClient) *JuGuangNoteClient {
	return &JuGuangNoteClient{
		BaseClient: c,
	}
}

// ## GetDirectLinkList 获取直达链接列表
//
// 接口文档参考：
// https://ad.xiaohongshu.com/openApiDoc?uba_pre=18.target_package..1712048419713&uba_ppre=18.aurora_asset_manage..1712048416043&uba_index=11&articleId=3194
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • id: 直达链接ID
//	  • page_num: 页码，从1开始
//	  • page_size: 页大小，最大100
//	  • remark_name: 备注名称，支持模糊匹配
//	  • type_str: 链接类型，1-deeplink，2-ulk
//	  • status_str: 链接状态，1-审核中，2-审核通过，3-审核拒绝
//
// 返回值：
//
//	*schema.JuGuangNoteGetDirectLinkListRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 直达链接数据
//	    - total: 直达链接总数
//	    - direct_link_list: 直达链接列表
//	      • id: 直达链接ID
//	      • url: 链接地址
//	      • type: 链接类型
//	      • status: 链接状态
//	      • remark_name: 备注名称
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) GetDirectLinkList(ctx context.Context, data *schema.JuGuangNoteGetDirectLinkListReq) (*schema.JuGuangNoteGetDirectLinkListRes, error) {
	result := &schema.JuGuangNoteGetDirectLinkListRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/direct-link/list", nil, data, nil, result)
	return result, err
}

// ## DeleteDirectLink 删除直达链接
//
// 接口文档参考：
// https://ad.xiaohongshu.com/openApiDoc?uba_pre=18.target_package..1712048419713&uba_ppre=18.aurora_asset_manage..1712048416043&uba_index=11&articleId=3193
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • id: 要删除的直达链接ID
//
// 返回值：
//
//	*schema.JuGuangNoteDeleteDirectLinkRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) DeleteDirectLink(ctx context.Context, data *schema.JuGuangNoteDeleteDirectLinkReq) (*schema.JuGuangNoteDeleteDirectLinkRes, error) {
	result := &schema.JuGuangNoteDeleteDirectLinkRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/direct-link/delete", nil, data, nil, result)
	return result, err
}

// ## CreateDirectLink 创建直达链接
//
// 接口文档参考：
// https://ad.xiaohongshu.com/openApiDoc?uba_pre=18.target_package..1712048419713&uba_ppre=18.aurora_asset_manage..1712048416043&uba_index=11&articleId=3191
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • direct_link_list: 直达链接信息列表
//	    - type: 直达链接类型，1-deeplink，2-ulk
//	    - url: url内容
//	    - remark_name: 备注名称（可选）
//
// 返回值：
//
//	*schema.JuGuangNoteCreateDirectLinkRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 创建结果
//	    - request_id: 请求ID
//	    - direct_link_list: 创建的直达链接列表
//	      • type: 直达链接类型
//	      • url: 链接地址
//	      • status: 链接状态（1-审核中，2-审核通过，3-审核拒绝）
//	      • remark_name: 备注名称
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) CreateDirectLink(ctx context.Context, data *schema.JuGuangNoteCreateDirectLinkReq) (*schema.JuGuangNoteCreateDirectLinkRes, error) {
	result := &schema.JuGuangNoteCreateDirectLinkRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/direct-link/create", nil, data, nil, result)
	return result, err
}

// ## GetQualInfo 获取资质列表
//
// 接口文档参考：
// https://ad.xiaohongshu.com/openApiDoc?uba_pre=18.target_package..1712048419713&uba_ppre=18.aurora_asset_manage..1712048416043&uba_index=11&articleId=3189
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • page_index: 页码，默认1
//	  • page_size: 每页查询量级，默认20
//
// 返回值：
//
//	*schema.JuGuangNoteGetQualInfoRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 资质信息
//	    - qual_infos: 资质列表
//	    - product_qual_infos: 行业资质信息
//	    - brand_qual_infos: 产品资质信息
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) GetQualInfo(ctx context.Context, data *schema.JuGuangNoteGetQualInfoReq) (*schema.JuGuangNoteGetQualInfoRes, error) {
	result := &schema.JuGuangNoteGetQualInfoRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/data/qual/info", nil, data, nil, result)
	return result, err
}

// ## GetPOIList 获取门店信息列表
//
// 接口文档参考：
// https://ad.xiaohongshu.com/openApiDoc?uba_pre=18.target_package..1712048419713&uba_ppre=18.aurora_asset_manage..1712048416043&uba_index=11&articleId=3169
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • page_index: 页码，默认1
//	  • page_size: 每页查询量级，默认20
//
// 返回值：
//
//	*schema.JuGuangNoteGetGetPOIListRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 门店信息
//	    - poi_info_list: 门店列表
//	      • poi_id: 门店ID
//	      • poi_name: 门店名称
//	    - page: 分页信息
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) GetPOIList(ctx context.Context, data *schema.JuGuangNoteGetGetPOIListReq) (*schema.JuGuangNoteGetGetPOIListRes, error) {
	result := &schema.JuGuangNoteGetGetPOIListRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/data/poi/list", nil, data, nil, result)
	return result, err
}

// ## ListLandingPage 获取落地页表单列表
//
// 接口文档参考：
// https://ad.xiaohongshu.com/openApiDoc?uba_pre=18.target_package..1712048419713&uba_ppre=18.aurora_asset_manage..1712048416043&uba_index=11&articleId=3161
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • page: 请求的页码，默认为1
//	  • page_size: 每页行数，默认为10
//	  • keyword: 搜索词
//	  • status: 固定传2(审核通过)
//
// 返回值：
//
//	*schema.JuGuangNoteListLandingPageRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据
//	    - total: 总数
//	    - list: 落地页信息列表
//	      • id: 页面ID
//	      • page_name: 页面名称
//	      • status: 落地页状态
//	      • audit_comment: 审核结果备注
//	      • content: 页面内容
//	      • web_domain: 在线域名
//	      • web_path: path
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) ListLandingPage(ctx context.Context, data *schema.JuGuangNoteListLandingPageReq) (*schema.JuGuangNoteListLandingPageRes, error) {
	result := &schema.JuGuangNoteListLandingPageRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/landing_page/list_landing_page", nil, data, nil, result)
	return result, err
}

// ## GetSpuList 获取SPU列表
//
// 接口文档参考：
// `https://ad.xiaohongshu.com/openApiDoc?uba_pre=18.target_package..1712048419713&uba_ppre=18.aurora_asset_manage..1712048416043&uba_index=11&articleId=2730`
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • can_bind: 是否可以绑定
//	  • keyword: 搜索关键词
//	  • page: 页码
//	  • page_size: 页大小
//
// 返回值：
//
//	*schema.JuGuangNoteGetSpuListRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据
//	    - total: 总数
//	    - spu: SPU信息列表
//	      • main_spu_id: 主SPU ID
//	      • spu_id: SPU ID
//	      • spu_name: SPU名称
//	      • spu_status: SPU状态
//	      • taxonomy_code: SPU类目code
//	      • brand_id: 品牌ID
//	      • nick_name_list: SPU昵称列表
//	      • series_list: SPU系列名称列表
//	      • pic_url_list: SPU图片链接列表
//	      • spu_audit_reason: 审核拒绝原因
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) GetSpuList(ctx context.Context, data *schema.JuGuangNoteGetSpuListReq) (*schema.JuGuangNoteGetSpuListRes, error) {
	result := &schema.JuGuangNoteGetSpuListRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/spu/list", nil, data, nil, result)
	return result, err
}
