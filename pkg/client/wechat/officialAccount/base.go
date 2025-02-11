package officialAccount

import (
	"context"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/response"
)

func (srv *WeChatOfficialAccountService) GetCallbackIP(ctx context.Context) (*response.GetCallBackIPRes, error) {

	result := &response.GetCallBackIPRes{}

	_, err := srv.Client.HttpGet(ctx, "cgi-bin/getcallbackip", nil, nil, result)

	return result, err
}
