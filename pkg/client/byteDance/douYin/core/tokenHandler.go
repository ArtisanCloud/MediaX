package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	config2 "github.com/ArtisanCloud/MediaXCore/pkg/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type ByteDanceTokenHandler struct {
	Config       *config.ClientConfig
	TokenHandler *kernel.TokenHandler
}

func NewByteDanceTokenHandler(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*ByteDanceTokenHandler, error) {
	handler, err := kernel.NewTokenHandler(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	byteDanceHandler := &ByteDanceTokenHandler{
		Config:       cfg,
		TokenHandler: handler,
	}

	if cfg.AccessTokenUrl != "" {
		byteDanceHandler.TokenHandler.EndpointToGetToken = cfg.AccessTokenUrl
	} else {
		byteDanceHandler.TokenHandler.EndpointToGetToken = config.ByteDanceDouYinAuthTokenUrl
	}
	byteDanceHandler.OverrideGetCredentials()

	return byteDanceHandler, nil
}

func (acHandler *ByteDanceTokenHandler) OverrideGetCredentials() {

	acHandler.TokenHandler.GetCredentials = func() *object.StringMap {
		// this is only for client credential flow
		return &object.StringMap{
			"grant_type":    string(config2.AuthFlowClientCred),
			"client_id":     acHandler.Config.ClientID,
			"client_secret": acHandler.Config.ClientSecret,
			//"neededText": "",
		}
	}
}
