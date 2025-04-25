package officialAccount

import (
	"context"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/response"
)

// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_the_WeChat_server_IP_address.html#2.%20%E8%8E%B7%E5%8F%96%E5%BE%AE%E4%BF%A1callback%20IP%E5%9C%B0%E5%9D%80
func (c *WeChatOfficialAccountClient) GetCallbackIp(ctx context.Context) (*response.GetCallBackIPRes, error) {

	result := &response.GetCallBackIPRes{}

	_, err := c.WeChatClient.HttpGet(ctx, "cgi-bin/getcallbackip", nil, nil, result)

	return result, err
}
