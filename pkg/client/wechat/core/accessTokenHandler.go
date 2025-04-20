package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type WeChatAccessTokenHandler struct {
	Config             *config.ClientConfig
	AccessTokenHandler *kernel.AccessTokenHandler
}

func NewWeChatAccessTokenHandler(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*WeChatAccessTokenHandler, error) {
	handler, err := kernel.NewAccessTokenHandler(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	wechatHandler := &WeChatAccessTokenHandler{
		Config:             cfg,
		AccessTokenHandler: handler,
	}

	if cfg.AccessTokenUrl != "" {
		wechatHandler.AccessTokenHandler.EndpointToGetToken = cfg.AccessTokenUrl
	} else {
		wechatHandler.AccessTokenHandler.EndpointToGetToken = "https://api.weixin.qq.com/cgi-bin/token"
	}
	wechatHandler.OverrideGetCredentials()

	return wechatHandler, nil
}

func (acHandler *WeChatAccessTokenHandler) OverrideGetCredentials() {

	acHandler.AccessTokenHandler.GetCredentials = func() *object.StringMap {
		return &object.StringMap{
			"grant_type": "client_credential",
			"appid":      acHandler.Config.ClientID,
			"secret":     acHandler.Config.ClientSecret,
			"neededText": "",
		}
	}
}
