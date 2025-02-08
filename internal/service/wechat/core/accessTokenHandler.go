package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type WeChatAccessTokenHandler struct {
	Config             *config.WeChatConfig
	AccessTokenHandler *kernel.AccessTokenHandler
}

func NewWeChatAccessTokenHandler(cfg *config.WeChatConfig, logger *logger.Logger, cache cache.CacheInterface) (*WeChatAccessTokenHandler, error) {
	handler, err := kernel.NewAccessTokenHandler(&cfg.AppConfig, logger, cache)
	if err != nil {
		return nil, err
	}
	wechatHandler := &WeChatAccessTokenHandler{
		Config:             cfg,
		AccessTokenHandler: handler,
	}

	wechatHandler.AccessTokenHandler.EndpointToGetToken = "https://api.weixin.qq.com/cgi-bin/token"
	wechatHandler.OverrideGetCredentials()

	return wechatHandler, nil
}

func (acHandler *WeChatAccessTokenHandler) OverrideGetCredentials() {

	acHandler.AccessTokenHandler.GetCredentials = func() *object.StringMap {
		return &object.StringMap{
			"grant_type": "client_credential",
			"appid":      acHandler.Config.AppID,
			"secret":     acHandler.Config.AppSecret,
			"neededText": "",
		}
	}
}
