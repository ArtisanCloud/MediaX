package base

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/base/schema"
)

// OfficialAccountBaseClient 是一个用于操作微信公众号基础功能的客户端。
type OfficialAccountBaseClient struct {
	*kernel.BaseClient

	AllowTypes []string
}

// NewClient 创建一个新的 OfficialAccountBaseClient 实例。
func NewClient(c *kernel.BaseClient) *OfficialAccountBaseClient {

	return &OfficialAccountBaseClient{
		BaseClient: c,
	}
}

// ## GetCallbackIp 获取微信服务器IP地址
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_the_WeChat_server_IP_address.html#2.%20%E8%8E%B7%E5%8F%96%E5%BE%AE%E4%BF%A1callback%20IP%E5%9C%B0%E5%9D%80
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*schema.GetCallBackIPRes 包含以下字段：
//	  • IpList: 微信服务器IP地址列表
//
//	error 调用过程中遇到的错误（如有）
func (c *OfficialAccountBaseClient) GetCallbackIp(ctx context.Context) (*schema.GetCallBackIPRes, error) {

	result := &schema.GetCallBackIPRes{}

	_, err := c.BaseClient.HttpGet(ctx, "cgi-bin/getcallbackip", nil, nil, result)

	return result, err
}
