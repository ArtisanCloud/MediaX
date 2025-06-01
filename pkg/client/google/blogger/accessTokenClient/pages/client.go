package pages

import (
	"context"
	"fmt"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/accessTokenClient/pages/schema"
)

type BloggerPagesClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerPagesClient {
	return &BloggerPagesClient{
		BaseClient: c,
	}
}

// ## List 获取页面列表
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/pages/list?hl=zh-cn
//
// 功能：
// 检索指定博客的页面列表，支持按状态(status)过滤和设置返回视图(view)
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BlogId: 博客ID(必填)
//	  • FetchBodies: 是否返回页面内容(可选)
//	  • Status: 页面状态过滤(draft/imported/live)(可选)
//	  • View: 返回视图级别(ADMIN/AUTHOR/READER)(可选)
//
// 返回值：
//
//	*schema.BloggerPagesListRes 包含以下字段：
//	  • Kind: 资源类型(blogger#pageList)
//	  • Items: 页面资源列表
//	error 调用过程中遇到的错误(如有)
//
// 授权范围：
// 对于私有博客需要 https://www.googleapis.com/auth/blogger 权限
func (c *BloggerPagesClient) List(ctx context.Context, data *schema.BloggerPagesListReq) (*schema.BloggerPagesListRes, error) {
	result := &schema.BloggerPagesListRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/pages", data.BlogId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// ## Get 获取单个页面详情
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/pages/get?hl=zh-cn
//
// 功能：
// 获取指定博客中单个页面的详细信息
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BlogId: 包含相应页面的博客ID(必填)
//	  • PageId: 要获取的页面ID(必填)
//	  • View: 可选，指定返回视图级别(ADMIN/AUTHOR/READER)
//
// 返回值：
//
//	*schema.BloggerPagesGetRes 包含页面详情，字段包括：
//	  • Kind: 资源类型(blogger#page)
//	  • Id: 页面ID
//	  • 其他页面字段
//	error 调用过程中遇到的错误(如有)
//
// 授权范围：
// 对于私有博客需要 https://www.googleapis.com/auth/blogger 权限
func (c *BloggerPagesClient) Get(ctx context.Context, data *schema.BloggerPagesGetReq) (*schema.BloggerPagesGetRes, error) {
	result := &schema.BloggerPagesGetRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/pages/%s", data.BlogId, data.PageId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// ## Delete 删除页面
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/pages/delete?hl=zh-cn
//
// 功能：
// 删除指定博客中的页面，可选择将其移至回收站
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BlogId: 包含页面的博客ID(必填)
//	  • PageId: 要删除的页面ID(必填)
//	  • UseTrash: 可选，是否将页面移至回收站
//
// 返回值：
//
//	*schema.BloggerPagesDeleteRes: 成功时返回空响应体
//	error: 调用过程中遇到的错误(如有)
//
// 授权范围：
// 需要 https://www.googleapis.com/auth/blogger 权限
func (c *BloggerPagesClient) Delete(ctx context.Context, data *schema.BloggerPagesDeleteReq) (*schema.BloggerPagesDeleteRes, error) {
	result := &schema.BloggerPagesDeleteRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/pages/%s", data.BlogId, data.PageId)
	_, err := c.BaseClient.HttpDelete(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// ## Insert 创建新页面
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/pages/insert?hl=zh-cn
//
// 功能：
// 在指定博客中创建新页面
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BlogId: 目标博客ID(必填)
//	  • Page: 要创建的页面资源(必填)
//
// 返回值：
//
//	*schema.BloggerPagesInsertRes: 包含新创建的页面详细信息
//	error: 调用过程中遇到的错误(如有)
//
// 授权范围：
// 需要 https://www.googleapis.com/auth/blogger 权限
func (c *BloggerPagesClient) Insert(ctx context.Context, data *schema.BloggerPagesInsertReq) (*schema.BloggerPagesInsertRes, error) {
	result := &schema.BloggerPagesInsertRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/pages", data.BlogId)
	_, err := c.BaseClient.HttpPost(ctx, endpoint, nil, data, nil, result)
	return result, err
}

// ## Patch 更新页面
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/pages/patch?hl=zh-cn
//
// 功能：
// 使用PATCH语义更新指定博客中的页面
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BlogId: 目标博客ID(必填)
//	  • PageId: 目标页面ID(必填)
//	  • Page: 要更新的页面资源(必填)
//
// 返回值：
//
//	*schema.BloggerPagesPatchRes: 包含更新后的页面详细信息
//	error: 调用过程中遇到的错误(如有)
//
// 授权范围：
// 需要 https://www.googleapis.com/auth/blogger 权限
func (c *BloggerPagesClient) Patch(ctx context.Context, data *schema.BloggerPagesPatchReq) (*schema.BloggerPagesPatchRes, error) {
	result := &schema.BloggerPagesPatchRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/pages/%s", data.BlogId, data.PageId)
	_, err := c.BaseClient.HttpPatch(ctx, endpoint, nil, data, nil, result)
	return result, err
}

// ## Update 更新页面
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/pages/update?hl=zh-cn
//
// 功能：
// 使用PUT语义完全更新指定博客中的页面
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BlogId: 目标博客ID(必填)
//	  • PageId: 目标页面ID(必填)
//	  • Page: 完整的页面资源(必填)
//
// 返回值：
//
//	*schema.BloggerPagesUpdateRes: 包含更新后的页面详细信息
//	error: 调用过程中遇到的错误(如有)
//
// 授权范围：
// 需要 https://www.googleapis.com/auth/blogger 权限
func (c *BloggerPagesClient) Update(ctx context.Context, data *schema.BloggerPagesUpdateReq) (*schema.BloggerPagesUpdateRes, error) {
	result := &schema.BloggerPagesUpdateRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/pages/%s", data.BlogId, data.PageId)
	_, err := c.BaseClient.HttpPut(ctx, endpoint, nil, data, nil, result)
	return result, err
}
