package video

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/article/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// BiliBiliArticleClient 表示 B 站文章相关的客户端
type BiliBiliArticleClient struct {
	*kernel.BaseClient
}

// NewClient 创建 B 站文章相关的客户端
func NewClient(c *kernel.BaseClient) *BiliBiliArticleClient {
	return &BiliBiliArticleClient{
		BaseClient: c,
	}
}

// ## Add 提交文章
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/b14b77b6-8889-8c8b-2e83-17c5a4c550fb#h1-u6587u7AE0u63D0u4EA4
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • title: 文章标题，建议30字以内，最多不超过40字
//	  • category: 文章分类，通过categories接口获取，填写子分类id
//	  • template_id: 选择稿件模板
//	  • summary: 文章简介，字数不少于200字
//	  • content: 正文内容，字数200-40000，或者至少添加三张图片
//	  • banner_url: 头图url，与top_video_bvid可都不填写，或者只能二选一
//	  • original: 是否原创：0-非原创 1-原创 默认0
//	  • image_urls: 封面默认使用正文文字图片，多个图片用英文逗号分割
//	  • tags: 用户自定义标签，多标签情况下用英文逗号分割
//	  • list_id: 文集id
//	  • up_closed_reply: 是否up主关闭评论区：0-否 1-是 默认0
//	  • top_video_bvid: 头部视频bvid，与banner_url可都不填写，或者只能二选一
//
// 返回值：
//
//	*schema.BiliBiliArticleAddRes 包含以下字段：
//	  • code: 响应码
//	  • message: 响应消息
//	  • request_id: 请求ID
//	  • ttl: TTL
//	  • data: 响应数据，包含新建的文章id
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliArticleClient) Add(ctx context.Context, data *schema.BiliBiliArticleAddReq) (*schema.BiliBiliArticleAddRes, error) {
	result := &schema.BiliBiliArticleAddRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/article/add", nil, data, nil, result)
	return result, err
}

// ## Edit 编辑文章
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/2b5284bd-9a40-247b-8da6-0ef7cd00afd3#h1-u6587u7AE0u7F16u8F91
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 文章ID
//	  • title: 文章标题，建议30字以内，最多不超过40字
//	  • category: 文章分类，通过categories接口获取，填写子分类id
//	  • template_id: 选择稿件模板
//	  • summary: 文章简介，字数不少于200字
//	  • content: 正文内容，字数200-40000，或者至少添加三张图片
//	  • banner_url: 头图url，与top_video_bvid可都不填写，或者只能二选一
//	  • original: 是否原创：0-非原创 1-原创 默认0
//	  • image_urls: 封面默认使用正文文字图片，多个图片用英文逗号分割
//	  • tags: 用户自定义标签，多标签情况下用英文逗号分割
//	  • list_id: 文集id
//	  • up_closed_reply: 是否up主关闭评论区：0-否 1-是 默认0
//	  • top_video_bvid: 头部视频bvid，与banner_url可都不填写，或者只能二选一
//
// 返回值：
//
//	*schema.BiliBiliArticleEditRes 包含以下字段：
//	  • code: 响应码
//	  • message: 响应消息
//	  • request_id: 请求ID
//	  • ttl: TTL
//	  • data: 响应数据，包含编辑后的文章id
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliArticleClient) Edit(ctx context.Context, data *schema.BiliBiliArticleEditReq) (*schema.BiliBiliArticleEditRes, error) {
	result := &schema.BiliBiliArticleEditRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/article/edit", nil, data, nil, result)
	return result, err
}

// ## Delete 删除文章
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/b63f8918-2add-0fbb-0718-d0537329ed1c#h1-u6587u7AE0u5220u9664
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • ID: 文章ID（必填，支持多个ID，用英文逗号分隔）
//
// 返回值：
//
//	*schema.BiliBiliArticleDeleteRes 包含以下字段：
//	  • Code: 返回码
//	  • Message: 返回信息
//	  • RequestID: 请求ID
//	  • TTL: 请求有效期
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliArticleClient) Delete(ctx context.Context, data *schema.BiliBiliArticleDeleteReq) (*schema.BiliBiliArticleDeleteRes, error) {
	result := &schema.BiliBiliArticleDeleteRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/article/delete", nil, data, nil, result)
	return result, err
}

// ## Detail 查询文章详情
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/243457f6-9d44-9348-06ef-366cdbb44eb4#h1-u67E5u8BE2u6587u7AE0u8BE6u60C5
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • ID: 文章ID（必填）
//
// 返回值：
//
//	*schema.BiliBiliArticleDetailRes 包含以下字段：
//	  • Code: 返回码
//	  • Message: 返回信息
//	  • TTL: 请求有效期
//	  • Data: 文章详情数据
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliArticleClient) Detail(ctx context.Context, data *schema.BiliBiliArticleDetailReq) (*schema.BiliBiliArticleDetailRes, error) {
	result := &schema.BiliBiliArticleDetailRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/article/detail", nil, data, nil, result)
	return result, err
}

// ## List 查询文章列表
// 接口文档参考：
// https://open.bilibili.com/doc/4/c8057666-2b92-fc72-3607-f4199933dc13#h1-u67E5u8BE2u6587u7AE0u5217u8868
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • PN: 页码（可选，默认1）
//	  • PS: 每页大小（可选，默认10）
//	  • Sort: 排序类型（可选，默认1）
//	  • Group: 文章状态分组（可选，默认0）
//	  • Category: 分类ID（可选）
//
// 返回值：
//
//	*schema.BiliBiliArticleListRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 返回消息
//	  • TTL: 生存时间
//	  • Data: 文章列表数据
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliArticleClient) List(ctx context.Context, data *schema.BiliBiliArticleListReq) (*schema.BiliBiliArticleListRes, error) {
	result := &schema.BiliBiliArticleListRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/article/list", nil, data, nil, result)
	return result, err
}

// ## Categories 获取文章的分类列表
// 接口文档参考：
// https://open.bilibili.com/doc/4/120b97a6-24a1-4e9f-47b7-9a201196e5d3#h1-u83B7u53D6u6587u7AE0u7684u5206u7C7Bu5217u8868
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//
// 返回值：
//
//	*schema.BiliBiliArticleCategoriesRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 返回消息
//	  • RequestID: 请求ID
//	  • TTL: 生存时间
//	  • Data: 分类列表数据
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliArticleClient) Categories(ctx context.Context, data *schema.BiliBiliArticleCategoriesReq) (*schema.BiliBiliArticleCategoriesRes, error) {
	result := &schema.BiliBiliArticleCategoriesRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/article/categories", nil, data, nil, result)
	return result, err
}

// ## Card 获取视频、文章卡片信息
// 接口文档参考：
// https://open.bilibili.com/doc/4/120b97a6-24a1-4e9f-47b7-9a201196e5d3#h1--
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • ResourceID: 资源id（视频BV号/文章cv号）：BV开头为视频，cv开头为文章
//
// 返回值：
//
//	*schema.BiliBiliArticleCardRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 返回消息
//	  • RequestID: 请求ID
//	  • TTL: 生存时间
//	  • Data: 卡片片段
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliArticleClient) Card(ctx context.Context, data *schema.BiliBiliArticleCardReq) (*schema.BiliBiliArticleCardRes, error) {
	result := &schema.BiliBiliArticleCardRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/article/cards", nil, data, nil, result)
	return result, err
}

// ## UploadImage 专栏稿件图片上传
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/0eaa4d3e-c4c0-f874-6f3c-e083aa939a1b#h1-u6587u96C6u4FE1u606Fu7F16u8F91
//
// 参数：
//
//	ctx - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • Files: 图片文件
//	  • Form: 表单数据
//	  • Watermark: 是否带水印
//
// 返回值：
//
//	*schema.BiliBiliArticleUploadImageRes 包含以下字段：
//	  • Code: 响应状态码
//	  • Message: 响应消息
//	  • Data: 图片上传结果数据
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliArticleClient) UploadImage(ctx context.Context, data *schema.BiliBiliArticleUploadImageReq) (*schema.BiliBiliArticleUploadImageRes, error) {
	result := &schema.BiliBiliArticleUploadImageRes{}
	_, err := c.BaseClient.HttpUpload(ctx, "/arcopen/fn/article/upload/image",
		data.Files, data.Form,
		nil, &object.HashMap{
			"watermark": data.Watermark,
		}, nil, result)
	return result, err
}
