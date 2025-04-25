package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type ByteDanceAccessTokenHandler struct {
	Config             *config.ClientConfig
	AccessTokenHandler *kernel.AccessTokenHandler
}

func NewByteDanceAccessTokenHandler(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*ByteDanceAccessTokenHandler, error) {
	handler, err := kernel.NewAccessTokenHandler(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	byteDanceHandler := &ByteDanceAccessTokenHandler{
		Config:             cfg,
		AccessTokenHandler: handler,
	}

	if cfg.AccessTokenUrl != "" {
		byteDanceHandler.AccessTokenHandler.EndpointToGetToken = cfg.AccessTokenUrl
	} else {
		byteDanceHandler.AccessTokenHandler.EndpointToGetToken = config.ByteDanceDouYinAuthTokenUrl
	}
	byteDanceHandler.OverrideGetCredentials()

	return byteDanceHandler, nil
}

func (acHandler *ByteDanceAccessTokenHandler) OverrideGetCredentials() {

	acHandler.AccessTokenHandler.GetCredentials = func() *object.StringMap {
		return &object.StringMap{
			"grant_type":    "authorization_code",
			"client_id":     acHandler.Config.ClientID,
			"client_secret": acHandler.Config.ClientSecret,
			//"neededText": "",
		}
	}
}
