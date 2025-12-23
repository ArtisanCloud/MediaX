package oauth

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/oauth/schema"
)

type DouYinOAuthClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinOAuthClient {
	return &DouYinOAuthClient{
		BaseClient: c,
	}
}

// ## GetOpenIdByC 移动网站应用C端联合授权获取OpenID
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/union_auth/get-openid-by-c
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 获取OpenID请求参数，包含以下字段：
//	  • BClientKey: 绑定的B端应用ID
//	  • CClientKey: C端移动网站应用对应的应用ID
//	  • OpenID: 用户在当前C端移动网站应用下的OpenID
//
// 返回值：
//
//	*schema.DouYinOAuthGetOpenIdByCRes 包含以下字段：
//	  • OpenID: 用户在B端应用下的OpenID
//	error 调用过程中遇到的错误（如有）
func (c *DouYinOAuthClient) GetOpenIdByC(ctx context.Context, data *schema.DouYinOAuthGetOpenIdByCReq) (*schema.DouYinOAuthGetOpenIdByCRes, error) {
	result := &schema.DouYinOAuthGetOpenIdByCRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/douyin/v1/auth/get_openid_by_c/", nil, data, nil, result)
	return result, err
}

// ## GetOpenIdByB 移动网站应用B端联合授权获取OpenID
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/union_auth/get-openid-by-b
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 获取OpenID请求参数，包含以下字段：
//	  • BClientKey: B端移动网站应用对应的应用ID
//	  • CClientKey: 绑定的C端应用ID
//	  • OpenID: 用户在绑定的C端应用下的OpenID
//
// 返回值：
//
//	*schema.DouYinOAuthGetOpenIdByBRes 包含以下字段：
//	  • BClientKey: B端移动网站应用对应的应用ID
//	  • CClientKey: 绑定的C端应用ID
//	  • OpenID: 用户在B端应用下的OpenID
//	error 调用过程中遇到的错误（如有）
func (c *DouYinOAuthClient) GetOpenIdByB(ctx context.Context, data *schema.DouYinOAuthGetOpenIdByBReq) (*schema.DouYinOAuthGetOpenIdByBRes, error) {
	result := &schema.DouYinOAuthGetOpenIdByBRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/douyin/v1/auth/get_openid_by_b/", nil, data, nil, result)
	return result, err
}
