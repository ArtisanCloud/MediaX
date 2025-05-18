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
