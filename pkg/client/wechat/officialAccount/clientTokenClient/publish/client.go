package publish

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	response2 "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/publish/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// OfficialAccountPublishClient 是一个用于操作微信公众号发布功能的客户端。
type OfficialAccountPublishClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 OfficialAccountPublishClient 实例。
func NewClient(c *kernel.BaseClient) *OfficialAccountPublishClient {
	return &OfficialAccountPublishClient{
		BaseClient: c,
	}
}

// ## DraftAdd 新建草稿
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Add_draft.html
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 草稿内容，包含以下字段：
//	  • Articles: 图文素材列表
//
// 返回值：
//
//	*schema.DraftAddRes 包含以下字段：
//	  • MediaID: 媒体ID
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) DraftAdd(ctx context.Context, data *schema.DraftAddReq) (*schema.DraftAddRes, error) {
	result := &schema.DraftAddRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/draft/add", nil, data, nil, result)
	return result, err
}

// ## DraftGet 获取草稿
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Get_draft.html
//
// 参数：
//
//	ctx    - 请求上下文
//	mediaID - 媒体ID
//
// 返回值：
//
//	*schema.DraftGetRes 包含以下字段：
//	  • NewsItem: 图文素材列表，每个元素包含以下字段：
//	    • Title: 标题
//	    • ThumbMediaID: 封面图片媒体ID
//	    • ShowCoverPic: 是否显示封面图片
//	    • Author: 作者
//	    • Digest: 摘要
//	    • Content: 正文内容
//	    • URL: 原文链接
//	    • ContentSourceURL: 原文链接
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) DraftGet(ctx context.Context, mediaID string) (*schema.DraftGetRes, error) {
	result := &schema.DraftGetRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/draft/get", nil, &object.HashMap{
		"media_id": mediaID,
	}, nil, result)

	return result, err
}

// ## DraftDelete 删除草稿
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Delete_draft.html
//
// 参数：
//
//	ctx    - 请求上下文
//	mediaID - 媒体ID
//
// 返回值：
//
//	*response2.OfficialAccountRes 包含以下字段：
//	  • ErrCode: 错误码
//	  • ErrMsg: 错误信息
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) DraftDelete(ctx context.Context, mediaID string) (*response2.OfficialAccountRes, error) {
	result := &response2.OfficialAccountRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/draft/delete", nil, &object.HashMap{
		"media_id": mediaID,
	}, nil, result)

	return result, err
}

// ## DraftUpdate 修改草稿
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Update_draft.html
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 草稿内容，包含以下字段：
//	  • MediaID: 媒体ID
//	  • Articles: 图文素材列表
//
// 返回值：
//
//	*response2.OfficialAccountRes 包含以下字段：
//	  • ErrCode: 错误码
//	  • ErrMsg: 错误信息
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) DraftUpdate(ctx context.Context, data *schema.DraftUpdateReq) (*response2.OfficialAccountRes, error) {
	result := &response2.OfficialAccountRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/draft/update", nil, data, nil, result)

	return result, err
}

// ## DraftCount 获取草稿总数
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Count_drafts.html
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*schema.DraftCountRes 包含以下字段：
//	  • TotalCount: 草稿总数
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) DraftCount(ctx context.Context) (*schema.DraftCountRes, error) {
	result := &schema.DraftCountRes{}

	_, err := c.BaseClient.HttpGet(ctx, "cgi-bin/draft/count", &object.StringMap{}, nil, nil, result)

	return result, err
}

// ## DraftBatchGet 获取草稿列表
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Count_drafts.html
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 查询条件，包含以下字段：
//	  • Offset: 从全部草稿的该偏移位置开始返回
//	  • Count: 返回草稿的数量
//	  • NoContent: 是否不返回草稿内容
//
// 返回值：
//
//	*schema.BatchGetRes 包含以下字段：
//	  • TotalCount: 草稿总数
//	  • ItemCount: 本次调用获取的草稿数量
//	  • Item: 草稿列表，每个元素包含以下字段：
//	    • MediaID: 媒体ID
//	    • Content: 草稿内容（如果 NoContent 为 false）
//	    • UpdateTime: 更新时间
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) DraftBatchGet(ctx context.Context, data *schema.BatchGetReq) (*schema.BatchGetRes, error) {
	result := &schema.BatchGetRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/draft/batchget", nil, data, nil, result)

	return result, err
}

// ## DraftSwitch MP端开关
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Temporary_MP_Switch.html
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*response2.OfficialAccountRes 包含以下字段：
//	  • ErrCode: 错误码
//	  • ErrMsg: 错误信息
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) DraftSwitch(ctx context.Context) (*response2.OfficialAccountRes, error) {
	result := &response2.OfficialAccountRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/draft/switch", nil, &object.HashMap{}, nil, result)

	return result, err
}

// ## DraftCheckSwitch M检查P端开关
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Temporary_MP_Switch.html
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*schema.CheckSwitchRes 包含以下字段：
//	  • IsOpen: 开关状态（true/false）
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) DraftCheckSwitch(ctx context.Context) (*schema.CheckSwitchRes, error) {
	result := &schema.CheckSwitchRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/draft/switch", nil, &object.HashMap{}, &object.StringMap{
		"checkonly": "1",
	}, result)

	return result, err
}

// ## PublishSubmit 发布接口
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Publish.html
//
// 参数：
//
//	ctx    - 请求上下文
//	mediaID - 媒体ID
//
// 返回值：
//
//	*schema.PublishSubmitRes 包含以下字段：
//	  • PublishID: 发布任务ID
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) PublishSubmit(ctx context.Context, mediaID string) (*schema.PublishSubmitRes, error) {
	result := &schema.PublishSubmitRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/submit", nil, &object.HashMap{
		"media_id": mediaID,
	}, nil, result)

	return result, err
}

// ## PublishGet 发布状态轮询接口
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Get_status.html
//
// 参数：
//
//	ctx      - 请求上下文
//	publishID - 发布任务ID
//
// 返回值：
//
//	*schema.PublishGetRes 包含以下字段：
//	  • PublishID: 发布任务ID
//	  • Status: 发布状态
//	  • ArticleURL: 文章链接
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) PublishGet(ctx context.Context, publishID uint64) (*schema.PublishGetRes, error) {
	result := &schema.PublishGetRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/get", nil, &object.HashMap{
		"publish_id": publishID,
	}, nil, result)

	return result, err
}

// ## PublishDelete 删除发布
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Delete_posts.html
//
// 参数：
//
//	ctx      - 请求上下文
//	articleID - 文章ID
//	index    - 文章在图文素材中的位置
//
// 返回值：
//
//	*response2.OfficialAccountRes 包含以下字段：
//	  • ErrCode: 错误码
//	  • ErrMsg: 错误信息
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) PublishDelete(ctx context.Context, articleID string, index int) (*response2.OfficialAccountRes, error) {
	result := &response2.OfficialAccountRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/delete", nil, &object.HashMap{
		"article_id": articleID,
		"index":      index,
	}, nil, result)

	return result, err
}

// ## PublishGetArticle 通过 article_id 获取已发布文章
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Get_article_from_id.html
//
// 参数：
//
//	ctx      - 请求上下文
//	articleID - 文章ID
//
// 返回值：
//
//	*schema.PublishGetArticleRes 包含以下字段：
//	  • Article: 文章内容
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) PublishGetArticle(ctx context.Context, articleID string) (*schema.PublishGetArticleRes, error) {
	result := &schema.PublishGetArticleRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/getarticle", nil, &object.HashMap{
		"article_id": articleID,
	}, nil, result)

	return result, err
}

// ## PublishBatchGet 获取成功发布列表
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Publish/Get_publication_records.html
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 查询条件，包含以下字段：
//	  • Offset: 从全部发布记录的该偏移位置开始返回
//	  • Count: 返回发布记录的数量
//
// 返回值：
//
//	*schema.BatchGetRes 包含以下字段：
//	  • TotalCount: 发布记录总数
//	  • ItemCount: 本次调用获取的发布记录数量
//	  • Item: 发布记录列表，每个元素包含以下字段：
//	    • ArticleID: 文章ID
//	    • Title: 文章标题
//	    • URL: 文章链接
//	    • UpdateTime: 更新时间
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountPublishClient) PublishBatchGet(ctx context.Context, data *schema.BatchGetReq) (*schema.BatchGetRes, error) {
	result := &schema.BatchGetRes{}

	_, err := c.BaseClient.HttpPost(ctx, "cgi-bin/freepublish/batchget", nil, data, nil, result)

	return result, err
}
