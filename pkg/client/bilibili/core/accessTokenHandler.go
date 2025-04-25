package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type BiliBiliAccessTokenHandler struct {
	Config             *config.ClientConfig
	AccessTokenHandler *kernel.AccessTokenHandler
}

func NewBiliBiliAccessTokenHandler(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*BiliBiliAccessTokenHandler, error) {
	handler, err := kernel.NewAccessTokenHandler(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	biliBiliHandler := &BiliBiliAccessTokenHandler{
		Config:             cfg,
		AccessTokenHandler: handler,
	}

	//if cfg.OAuthUrl != "" {
	//	biliBiliHandler.AccessTokenHandler.EndpointToGetToken = cfg.OAuthUrl
	//} else {
	//	biliBiliHandler.AccessTokenHandler.EndpointToGetToken = "https://oauth2.biliBiliapis.com/token"
	//}
	biliBiliHandler.OverrideGetCredentials()

	return biliBiliHandler, nil
}

func (acHandler *BiliBiliAccessTokenHandler) OverrideGetCredentials() {

	acHandler.AccessTokenHandler.GetCredentials = func() *object.StringMap {
		return &object.StringMap{
			"grant_type":    "authorization_code",
			"client_id":     acHandler.Config.ClientID,
			"client_secret": acHandler.Config.ClientSecret,
			//"neededText": "",
		}
	}
}
