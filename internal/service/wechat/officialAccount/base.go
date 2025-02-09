package officialAccount

import (
	"context"
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/officialAccount/response"
)

func (srv *WeChatOfficialAccountService) GetCallbackIP(ctx context.Context) (*response.GetCallBackIPRes, error) {

	result := &response.GetCallBackIPRes{}

	_, err := srv.Client.HttpGet(ctx, "cgi-bin/getcallbackip", nil, nil, result)

	return result, err
}
