package postUserInfos

import (
	"context"
	"fmt"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/accessTokenClient/postUserInfos/schema"
)

type BloggerPostUserInfosClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerPostUserInfosClient {
	return &BloggerPostUserInfosClient{
		BaseClient: c,
	}
}

// ## Get 获取帖子用户信息
//
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/postUserInfos/get
//
// 功能：
//
//	根据用户ID、博客ID和帖子ID获取帖子用户信息
//
// 参数:
//
//	ctx: 请求上下文
//	data: 包含请求参数的 BloggerPostUserInfosGetReq 结构体
//	  - UserId: (必需) 用户ID（"self"或用户个人资料标识符）
//	  - BlogId: (必需) 博客的ID
//	  - PostId: (必需) 要获取的帖子的ID
//	  - MaxComments: (可选) 要为帖子检索的评论数量上限
//
// 返回值:
//
//	*schema.BloggerPostUserInfosGetRes: 包含帖子用户信息的响应结构体
//	error: 如果请求失败，则返回错误信息
//
// 授权范围:
//
//	此请求需要获得以下至少一个范围的授权：
//	- https://www.googleapis.com/auth/blogger
//	- https://www.googleapis.com/auth/blogger.readonly
func (c *BloggerPostUserInfosClient) Get(ctx context.Context, data *schema.BloggerPostUserInfosGetReq) (*schema.BloggerPostUserInfosGetRes, error) {
	result := &schema.BloggerPostUserInfosGetRes{}
	endpoint := fmt.Sprintf("/blogger/v3/users/%s/blogs/%s/posts/%s", data.UserId, data.BlogId, data.PostId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}

// PostUserInfos：list
// https://developers.google.com/blogger/docs/3.0/reference/postUserInfos/list?hl=zh-cn
func (c *BloggerPostUserInfosClient) List(ctx context.Context, data *schema.BloggerPostUserInfosListReq) (*schema.BloggerPostUserInfosListRes, error) {
	result := &schema.BloggerPostUserInfosListRes{}
	endpoint := fmt.Sprintf("/blogger/v3/users/%s/blogs/%s/posts", data.UserId, data.BlogId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}
