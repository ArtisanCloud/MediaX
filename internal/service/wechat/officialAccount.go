package wechat

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/service/wechat/response"
)

func (srv *WeChatService) Publish() {
	srv.Logger.Info("publish official article")
}

func (srv *WeChatService) GetCallbackIP(ctx context.Context) (*response.GetCallBackIPRes, error) {

	result := &response.GetCallBackIPRes{}

	_, err := srv.Client.HttpGet(ctx, "cgi-bin/getcallbackip", nil, nil, result)

	return result, err
}
