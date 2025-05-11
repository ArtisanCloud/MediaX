package user

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/accessTokenClient/user/schema"
)

// BiliBiliUserClient 是BiliBili用户API的客户端
type BiliBiliUserClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的BiliBiliUserClient实例
func NewClient(c *kernel.BaseClient) *BiliBiliUserClient {
	return &BiliBiliUserClient{
		BaseClient: c,
	}
}

// ## AccountScopes 查询用户已授权权限列表
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/08f935c5-29f1-e646-85a3-0b11c2830558#h1-u67E5u8BE2u7528u6237u5DF2u6388u6743u6743u9650u5217u8868
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*schema.BiliBiliUserRes 包含以下字段：
//	  • openid: 用户唯一标识
//	  • scopes: 用户已授权的权限点列表
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliUserClient) AccountScopes(ctx context.Context) (*schema.BiliBiliUserRes, error) {
	result := &schema.BiliBiliUserRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/user/account/scopes", nil, nil, result)
	return result, err
}
