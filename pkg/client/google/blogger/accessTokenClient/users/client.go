package users

import (
	"context"
	"fmt"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/accessTokenClient/users/schema"
)

type BloggerUsersClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerUsersClient {
	return &BloggerUsersClient{
		BaseClient: c,
	}
}

// ## Get 获取用户信息
//
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/users?hl=zh-cn
//
// 功能：
//
//	根据用户ID获取Blogger用户的基本信息
//
// 参数:
//
//	ctx: 请求上下文
//	data: 包含请求参数的 BloggerUsersGetReq 结构体
//	  - UserId: (必需) 当前用户的数字ID，或者为"self"表示当前认证用户
//
// 返回值:
//
//	*schema.BloggerUsersGetRes: 包含用户资源的响应结构体
//	error: 如果请求失败，则返回错误信息
//
// 授权范围:
//
//	此请求需要获得以下至少一个范围的授权：
//	- https://www.googleapis.com/auth/blogger
//	- https://www.googleapis.com/auth/blogger.readonly
func (c *BloggerUsersClient) Get(ctx context.Context, data *schema.BloggerUsersGetReq) (*schema.BloggerUsersGetRes, error) {
	result := &schema.BloggerUsersGetRes{}
	endpoint := fmt.Sprintf("/blogger/v3/users/%s", data.UserId)
	_, err := c.BaseClient.HttpGet(ctx, endpoint, nil, nil, nil, result)
	return result, err
}
