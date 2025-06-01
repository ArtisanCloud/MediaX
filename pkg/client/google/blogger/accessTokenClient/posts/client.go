package posts

import (
	"context"
	"fmt"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/accessTokenClient/posts/schema"
)

// BloggerPostsClient 博客文章客户端
type BloggerPostsClient struct {
	*kernel.BaseClient
}

// NewClient 初始化博客文章客户端
func NewClient(c *kernel.BaseClient) *BloggerPostsClient {
	return &BloggerPostsClient{
		BaseClient: c,
	}
}

// ## List 获取博客文章列表
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/posts/list?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • blogId: 要从中提取帖子的博客的ID（必填）
//	  • endDate: 要提取的最新发布日期，RFC 3339格式（可选）
//	  • fetchBodies: 是否包含帖子正文内容，默认true（可选）
//	  • fetchImages: 是否包含每个帖子的图片URL元数据（可选）
//	  • labels: 要搜索的标签列表，逗号分隔（可选）
//	  • maxResults: 要提取的帖子数量上限（可选）
//	  • orderBy: 排序字段，可选值：published（发布日期）、updated（更新时间）
//	  • sortOption: 排序方向，可选值：descending（降序，默认）、ascending（升序）
//	  • pageToken: 分页令牌（可选）
//	  • startDate: 要提取的最早发布日期，RFC 3339格式（可选）
//	  • status: 帖子状态，可选值：draft（草稿）、live（已发布）、scheduled（计划发布）
//	  • view: 详细信息级别，可选值：ADMIN、AUTHOR、READER
//
// 返回值：
//
//	*schema.BloggerPostsListRes 包含以下字段：
//	  • Kind: 资源类型，固定值：blogger#postList
//	  • NextPageToken: 用于提取下一页的分页令牌
//	  • Items: 博文列表
//	error 调用过程中遇到的错误（如有）
//
// 如果博文位于私享博客上，则需要进行授权。如果博文位于公开的博客上，则用户无需授权即可调用此方法。
func (c *BloggerPostsClient) List(ctx context.Context, data *schema.BloggerPostsListReq) (*schema.BloggerPostsListRes, error) {
	result := &schema.BloggerPostsListRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts", data.BlogId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// ## Get 根据ID获取博客文章
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/posts/get?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • blogId: 要从中提取博文的博客的ID（必填）
//	  • postId: 帖子的ID（必填）
//	  • maxComments: 要检索的评论数量上限，不指定则不返回评论（可选）
//	  • view: 详细信息级别，可选值：ADMIN（管理员级别）、AUTHOR（作者级别）、READER（读者级别）
//
// 返回值：
//
//	*schema.BloggerPostsGetRes 包含以下字段：
//	  • 继承自 Post 结构体的所有博客文章字段
//	error 调用过程中遇到的错误（如有）
//
// 如果博文位于私享博客上，则需要进行授权。如果博文位于公开的博客上，则用户无需授权即可调用此方法。
func (c *BloggerPostsClient) Get(ctx context.Context, data *schema.BloggerPostsGetReq) (*schema.BloggerPostsGetRes, error) {
	result := &schema.BloggerPostsGetRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts/%s", data.BlogId, data.PostId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// ## Search 搜索博客文章
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/posts/search?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • blogId: 要搜索的博客 ID（必填）
//	  • q: 搜索查询字词（必填）
//	  • fetchBodies: 是否包含帖子正文内容（可选，默认为 true）
//	  • orderBy: 搜索结果排序顺序（可选，可选值："published"、"updated"）
//
// 返回值：
//
//	*schema.BloggerPostsSearchRes 包含以下字段：
//	  • Kind: 资源类型，固定为 blogger#postList
//	  • NextPageToken: 下一页分页令牌
//	  • Items: 博文列表
//	error 调用过程中遇到的错误（如有）
//
// 注意：如果要搜索的博客是私人博客，则需要进行授权
func (c *BloggerPostsClient) Search(ctx context.Context, data *schema.BloggerPostsSearchReq) (*schema.BloggerPostsSearchRes, error) {
	result := &schema.BloggerPostsSearchRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts/search", data.BlogId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// ## Insert 创建新博客文章
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/posts/insert?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BlogId: 要向其添加博文的博客的 ID（必填）
//	  • IsDraft: 是否以草稿形式创建博文（可选）
//	  • Post: 嵌入的 Post 结构体，包含文章内容
//
// 返回值：
//
//	*schema.BloggerPostsInsertRes 包含新创建的博文资源
//	error 调用过程中遇到的错误（如有）
//
// 授权要求：
// 此请求需要获得以下至少一个范围的授权：
// https://www.googleapis.com/auth/blogger
func (c *BloggerPostsClient) Insert(ctx context.Context, data *schema.BloggerPostsInsertReq) (*schema.BloggerPostsInsertRes, error) {
	result := &schema.BloggerPostsInsertRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts", data.BlogId)
	_, err := c.BaseClient.HttpPost(ctx, endpoint, nil, data, nil, result)
	return result, err
}

// ## Delete 删除指定博客文章
//
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/posts/delete?hl=zh-cn
//
// 参数:
//
//	ctx: 请求上下文。
//	data: 包含请求参数的 BloggerPostsDeleteReq 结构体。
//	  - BlogId: (必需) 要删除文章的博客ID。
//	  - PostId: (必需) 要删除的文章ID。
//	  - UseTrash: (可选) 是否将文章移至回收站而非直接删除。
//
// 返回值:
//
//	*schema.BloggerPostsDeleteRes: 空响应结构体。
//	error: 如果请求失败，则返回错误信息。
//
// 授权范围:
//
//	此请求需要获得以下至少一个范围的授权：
//	- https://www.googleapis.com/auth/blogger
func (c *BloggerPostsClient) Delete(ctx context.Context, data *schema.BloggerPostsDeleteReq) (*schema.BloggerPostsDeleteRes, error) {
	result := &schema.BloggerPostsDeleteRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts/%s", data.BlogId, data.PostId)
	_, err := c.BaseClient.HttpDelete(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// ## GetByPath 通过路径获取博客文章
//
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/posts/getByPath?hl=zh-cn
//
// 参数:
//
//	ctx: 请求上下文
//	data: 包含请求参数的 BloggerPostsGetByPathReq 结构体
//	  - BlogId: (必需) 要从中提取博文的博客ID
//	  - Path: (必需) 要检索的博文路径
//	  - MaxComments: (可选) 为帖子检索的评论数量上限
//	  - View: (可选) 视图级别(ADMIN/AUTHOR/READER)
//
// 返回值:
//
//	*schema.BloggerPostsGetByPathRes: 包含博文资源的响应结构体
//	error: 如果请求失败，则返回错误信息
//
// 注意:
//
//	如果博文位于私享博客上，则需要授权。博文路径是博文网址中位于主机之后的部分
func (c *BloggerPostsClient) GetByPath(ctx context.Context, data *schema.BloggerPostsGetByPathReq) (*schema.BloggerPostsGetByPathRes, error) {
	result := &schema.BloggerPostsGetByPathRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts/bypath", data.BlogId, data.Path)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// ## Patch 更新博客文章
//
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/posts/patch?hl=zh-cn
//
// 参数:
//
//	ctx: 请求上下文
//	data: 包含请求参数的 BloggerPostsPatchReq 结构体
//	  - BlogId: (必需) 要更新的博客ID
//	  - PostId: (必需) 要更新的帖子ID
//	  - Post: (必需) 包含更新内容的Posts资源
//
// 返回值:
//
//	*schema.BloggerPostsPatchRes: 包含更新后帖子资源的响应结构体
//	error: 如果请求失败，则返回错误信息
//
// 授权范围:
//
//	此请求需要获得以下至少一个范围的授权：
//	- https://www.googleapis.com/auth/blogger
func (c *BloggerPostsClient) Patch(ctx context.Context, data *schema.BloggerPostsPatchReq) (*schema.BloggerPostsPatchRes, error) {
	result := &schema.BloggerPostsPatchRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts/%s", data.BlogId, data.PostId)
	_, err := c.BaseClient.HttpPatch(ctx, endpoint, nil, data, nil, result)
	return result, err
}

// ## Update 更新博客文章
//
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/posts/update?hl=zh-cn
//
// 参数:
//
//	ctx: 请求上下文
//	data: 包含请求参数的 BloggerPostsUpdateReq 结构体
//	  - BlogId: (必需) 要更新的博客ID
//	  - PostId: (必需) 要更新的帖子ID
//	  - Post: (必需) 包含更新内容的Posts资源
//
// 返回值:
//
//	*schema.BloggerPostsUpdateRes: 包含更新后帖子资源的响应结构体
//	error: 如果请求失败，则返回错误信息
//
// 授权范围:
//
//	此请求需要获得以下至少一个范围的授权：
//	- https://www.googleapis.com/auth/blogger
func (c *BloggerPostsClient) Update(ctx context.Context, data *schema.BloggerPostsUpdateReq) (*schema.BloggerPostsUpdateRes, error) {
	result := &schema.BloggerPostsUpdateRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts/%s", data.BlogId, data.PostId)
	_, err := c.BaseClient.HttpPut(ctx, endpoint, nil, data, nil, result)
	return result, err
}

// ## Publish 发布博客文章
//
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/posts/publish?hl=zh-cn
//
// 参数:
//
//	ctx: 请求上下文
//	data: 包含请求参数的 BloggerPostsPublishReq 结构体
//	  - BlogId: (必需) 要发布文章的博客ID
//	  - PostId: (必需) 要发布的帖子ID
//	  - PublishDate: (可选) 安排发布的日期时间
//
// 返回值:
//
//	*schema.BloggerPostsPublishRes: 包含发布后帖子资源的响应结构体
//	error: 如果请求失败，则返回错误信息
//
// 授权范围:
//
//	此请求需要获得以下至少一个范围的授权：
//	- https://www.googleapis.com/auth/blogger
func (c *BloggerPostsClient) Publish(ctx context.Context, data *schema.BloggerPostsPublishReq) (*schema.BloggerPostsPublishRes, error) {
	result := &schema.BloggerPostsPublishRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts/%s/publish", data.BlogId, data.PostId)
	_, err := c.BaseClient.HttpPost(ctx, endpoint, nil, data, nil, result)
	return result, err
}

// ## Revert 还原博客文章为草稿
//
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/posts/revert?hl=zh-cn
//
// 功能：
//
//	将已发布或计划发布的博文还原为草稿状态，这会从可公开查看的内容中移除相应博文
//
// 参数:
//
//	ctx: 请求上下文
//	data: 包含请求参数的 BloggerPostsRevertReq 结构体
//	  - BlogId: (必需) 博客的ID
//	  - PostId: (必需) 帖子的ID
//
// 返回值:
//
//	*schema.BloggerPostsRevertRes: 包含还原后帖子资源的响应结构体
//	error: 如果请求失败，则返回错误信息
//
// 授权范围:
//
//	此请求需要获得以下至少一个范围的授权：
//	- https://www.googleapis.com/auth/blogger
func (c *BloggerPostsClient) Revert(ctx context.Context, data *schema.BloggerPostsRevertReq) (*schema.BloggerPostsRevertRes, error) {
	result := &schema.BloggerPostsRevertRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/posts/%s/revert", data.BlogId, data.PostId)
	_, err := c.BaseClient.HttpPost(ctx, endpoint, nil, data, nil, result)
	return result, err
}
