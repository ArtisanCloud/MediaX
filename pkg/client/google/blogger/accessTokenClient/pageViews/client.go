package pageViews

import (
	"context"
	"fmt"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/accessTokenClient/pageViews/schema"
)

// BloggerPageViewsClient 博客页面浏览量统计客户端
type BloggerPageViewsClient struct {
	*kernel.BaseClient
}

// NewClient 初始化博客页面浏览量统计客户端
func NewClient(c *kernel.BaseClient) *BloggerPageViewsClient {
	return &BloggerPageViewsClient{
		BaseClient: c,
	}
}

// ## Get 获取博客页面浏览量统计
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/pageViews/get?hl=zh-cn
//
// 功能：
// 检索指定博客的网页浏览统计信息
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BlogId: 要获取的博客ID(必填)
//	  • Range: 时间范围(可选)，可选值：
//	    - "30DAYS": 过去30天的网页浏览量
//	    - "7DAYS": 过去7天的网页浏览量
//	    - "all": 所有页面的总浏览量
//
// 返回值：
//
//	*schema.BloggerPageViewsGetRes: 包含博客页面浏览量统计信息
//	error: 调用过程中遇到的错误(如有)
//
// 授权范围：
// 需要 https://www.googleapis.com/auth/blogger 权限
func (c *BloggerPageViewsClient) Get(ctx context.Context, pageToken string, data *schema.BloggerPageViewsGetReq) (*schema.BloggerPageViewsGetRes, error) {
	result := &schema.BloggerPageViewsGetRes{}
	endpoint := fmt.Sprintf("/blogger/v3/blogs/%s/pageviews", data.BlogId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, data, nil, result)
	return result, err
}
