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

// ## QueryTemplate 获取定向包列表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=4309
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • target_template_id: 定向包ID（可选，根据定向包ID查询时其他参数无效）
//	  • marketing_target: 营销诉求（可选）
//	  • placement: 广告类型（可选）
//	  • page: 页码，默认1（可选）
//	  • page_size: 每页行数，默认20，最大100（可选）
//
// 返回值：
//
//	  *schema.JuGuangNoteQueryTemplateRes 包含以下字段：
//	    • code: 响应码
//	    • msg: 响应信息
//	    • success: 是否成功
//	    • data: 定向包数据列表
//	      - id: 定向包ID
//	      - name: 定向包名称
//	      - desc: 定向包描述
//	      - target_type: 定向类型
//	      - marketing_target: 营销诉求
//	      - placement: 广告类型
//	      - delivery_mode: 投放模式
//	      - unit_list: 已关联单元ID列表
//	      - target_config: 定向配置
//	    • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) QueryTemplate(ctx context.Context, data *schema.JuGuangNoteQueryTemplateReq) (*schema.JuGuangNoteQueryTemplateRes, error) {
	result := &schema.JuGuangNoteQueryTemplateRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/target/template/query", nil, data, nil, result)
	return result, err
}

// ## QueryNoteID 通过临时笔记ID获取正式笔记ID
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3787
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • temporary_id: 临时笔记ID（必填）
//
// 返回值：
//
//	*schema.JuGuangNoteQueryNoteIDRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据
//	      - temporary_id: 临时笔记ID
//	      - note_id: 正式笔记ID
//	      - create_time: 创建时间
//	  • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) QueryNoteID(ctx context.Context, data *schema.JuGuangNoteQueryNoteIDReq) (*schema.JuGuangNoteQueryNoteIDRes, error) {
	result := &schema.JuGuangNoteQueryNoteIDRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/noteid/query", nil, data, nil, result)
	return result, err
}

// ## DeleteNote 删除笔记
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3785
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • note_id: 笔记ID（必填）
//
// 返回值：
//
//	*schema.JuGuangNoteDeleteNoteRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) DeleteNote(ctx context.Context, data *schema.JuGuangNoteDeleteNoteReq) (*schema.JuGuangNoteDeleteNoteRes, error) {
	result := &schema.JuGuangNoteDeleteNoteRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/note/delete", nil, data, nil, result)
	return result, err
}

// ## PostNote 发布图文笔记
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3734
//
// <span style="color:red">接口使用门槛：品牌客户需满足日耗大于10万，代理商客户需满足日耗大于80万，如您满足要求请联系对接的渠道经理申请相应的审批流程。</span>
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • cover_datas: 图片信息列表（必填，每项包含 image_url, imag_width, imag_height）
//	  • common_data: 基础数据（非必填，包含 title, desc）
//	  • tags: 话题信息（非必填，最多45条）
//
// 返回值：
//
//	*schema.JuGuangNotePostNoteRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据
//	      - temporary_id: 笔记临时ID（用于后续查询正式笔记ID）
//	  • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) PostNote(ctx context.Context, data *schema.JuGuangNotePostNoteReq) (*schema.JuGuangNotePostNoteRes, error) {
	result := &schema.JuGuangNotePostNoteRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/image-note/post", nil, data, nil, result)
	return result, err
}

// ## GetNegativeWordList 查询否定词列表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3696
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • unit_id: 单元ID（必填）
//	  • page_num: 页码，默认1（可选）
//	  • page_size: 每页行数，默认20（可选）
//
// 返回值：
//
//	*schema.JuGuangNoteGetNegativeWordListRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • request_id: 请求ID
//	  • data: 返回数据
//	      - page: 分页信息（page_index: 当前页码, total_count: 总数量）
//	      - negative_keywords: 否定词详情列表
//	          • campaign_id: 计划ID
//	          • unit_id: 单元ID
//	          • negative_keyword_id: 否定词ID
//	          • keyword: 否定词内容
//	          • phrase_match_type: 匹配方式（0-精确匹配，1-短语匹配）
//
// error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) GetNegativeWordList(ctx context.Context, data *schema.JuGuangNoteGetNegativeWordListReq) (*schema.JuGuangNoteGetNegativeWordListRes, error) {
	result := &schema.JuGuangNoteGetNegativeWordListRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/negative/keyword/list", nil, data, nil, result)
	return result, err
}

// ## BatchDeleteNegativeWord 批量删除否定词
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3695
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • unit_id: 单元ID（必填）
//	  • negative_keyword_ids: 否定词ID列表（必填）
//
// 返回值：
//
//	*schema.JuGuangNoteBatchDeleteNegativeWordRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • request_id: 请求ID
//
// error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) BatchDeleteNegativeWord(ctx context.Context, data *schema.JuGuangNoteBatchDeleteNegativeWordReq) (*schema.JuGuangNoteBatchDeleteNegativeWordRes, error) {
	result := &schema.JuGuangNoteBatchDeleteNegativeWordRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/negative/keyword/batch/delete", nil, data, nil, result)
	return result, err
}

// ## BatchAddNegativeWord 批量添加否定词
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3676
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • unit_id: 单元ID（必填）
//	  • keywords: 否定词列表（必填，每项包含 keyword: 否定词内容, phrase_match_type: 匹配方式，0-精确匹配，1-短语匹配）
//
// 返回值：
//
//	*schema.JuGuangNoteBatchAddNegativeWordRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • request_id: 请求ID
//
// error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) BatchAddNegativeWord(ctx context.Context, data *schema.JuGuangNoteBatchAddNegativeWordReq) (*schema.JuGuangNoteBatchAddNegativeWordRes, error) {
	result := &schema.JuGuangNoteBatchAddNegativeWordRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/negative/keyword/batch/add", nil, data, nil, result)
	return result, err
}

// ## PostVideoNote 发布视频笔记
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3591
//
// <span style="color:red">接口使用门槛：品牌客户需满足日耗大于10万，代理商客户需满足日耗大于80万，如您满足要求请联系对接的渠道经理申请相应的审批流程。</span>
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • video_data: 视频信息（必填，包含 video_url: 视频地址，video_width: 视频宽度，video_height: 视频高度）
//	  • cover_data: 视频封面信息（必填，包含 image_url: 封面图片地址，imag_width: 封面宽度，imag_height: 封面高度）
//	  • common_data: 基础数据（非必填，包含 title: 标题, desc: 描述）
//	  • tags: 话题信息（非必填，最多45条）
//
// 返回值：
//
//	*schema.JuGuangNotePostVideoNoteRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据
//	      - temporary_id: 笔记临时ID（用于后续查询正式笔记ID）
//	  • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) PostVideoNote(ctx context.Context, data *schema.JuGuangNotePostVideoNoteReq) (*schema.JuGuangNotePostVideoNoteRes, error) {
	result := &schema.JuGuangNotePostVideoNoteRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/video-note/post", nil, data, nil, result)
	return result, err
}

// ## GetDirectLinkList 获取直达链接列表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3194
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
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3193
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

// ## EditDirectLink 编辑直达链接
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3192
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • direct_link: 直达链接信息
//	    - id: 直达链接ID
//	    - url: URL内容
//	    - remark_name: 备注名称（可选）
//
// 返回值：
//
//	*schema.JuGuangNoteEditDirectLinkRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) EditDirectLink(ctx context.Context, data *schema.JuGuangNoteEditDirectLinkReq) (*schema.JuGuangNoteEditDirectLinkRes, error) {
	result := &schema.JuGuangNoteEditDirectLinkRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/direct-link/update", nil, data, nil, result)
	return result, err
}

// ## CreateDirectLink 创建直达链接
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3191
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

// ## GetAssetInfo 获取资产事件信息
//
// 接口文档参考：
// `https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3190`
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID
//	  • page_index: 页码，从1开始
//	  • page_size: 每页查询量级
//
// 返回值：
//
//	*schema.JuGuangNoteGetAssetInfoRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 资产事件信息
//	    - event_assert_dtos: 资产列表
//	      • event_asset_id: 资产ID
//	      • event_asset_name: 资产名称
//	      • status: 资产状态
//	    - page: 分页信息
//	  • request_id: 请求ID
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) GetAssetInfo(ctx context.Context, data *schema.JuGuangNoteGetAssetInfoReq) (*schema.JuGuangNoteGetAssetInfoRes, error) {
	result := &schema.JuGuangNoteGetAssetInfoRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/data/event/asset/info", nil, data, nil, result)
	return result, err
}

// ## GetQualInfo 获取资质列表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3189
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
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3169
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
//	*schema.JuGuangNoteGetPOIListRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 门店信息
//	    - poi_info_list: 门店列表
//	      • poi_id: 门店ID
//	      • poi_name: 门店名称
//	    - page: 分页信息
//	error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) GetPOIList(ctx context.Context, data *schema.JuGuangNoteGetPOIListReq) (*schema.JuGuangNoteGetPOIListRes, error) {
	result := &schema.JuGuangNoteGetPOIListRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/data/poi/list", nil, data, nil, result)
	return result, err
}

// ## GetNoteList 获取笔记列表
//
// 接口文档参考：
// `https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3162`
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • note_type: 笔记类型（必填，1：我的笔记，2：合作笔记，4：主理人笔记，6：员工笔记，11：授权笔记）
//	  • keyword: 搜索关键词（可选，支持笔记ID或笔记名称）
//	  • page: 请求的页码（可选，默认1）
//	  • page_size: 每页行数（可选，默认20，最大100）
//	  • base_only: 是否只拉取笔记基本信息（必填，默认false）
//
// 返回值：
//
//	*schema.JuGuangNoteGetNoteListRes 包含以下字段：
//	  • code: 返回码
//	  • msg: 返回信息
//	  • success: 接口是否成功
//	  • data: 返回数据
//	      - total: 总数
//	      - notes: 笔记信息列表
//	          • note_id: 笔记ID
//	          • image: 图片
//	          • desc: 笔记内容
//	          • create_time: 创建时间
//	          • author: 笔记作者
//	          • author_image: 作者头像
//	          • status: 笔记状态
//	          • note_content_type: 笔记类型（1：图文笔记，2：视频笔记）
//	          • title: 标题
//	          • image_list: 图片列表（仅返回笔记封面图片）
//	          • note_url: 笔记链接（有效期两个月）
//
// error 调用过程中遇到的错误（如有）
func (client *JuGuangNoteClient) GetNoteList(ctx context.Context, data *schema.JuGuangNoteGetNoteListReq) (*schema.JuGuangNoteGetNoteListRes, error) {
	result := &schema.JuGuangNoteGetNoteListRes{}
	_, err := client.BaseClient.HttpPost(ctx, "/api/open/jg/note/list", nil, data, nil, result)
	return result, err
}

// ## ListLandingPage 获取落地页表单列表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3161
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
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=2730
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
