package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type RedBookAccessTokenHandler struct {
	Config             *config.ClientConfig
	AccessTokenHandler *kernel.TokenHandler
}

func NewRedBookAccessTokenHandler(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*RedBookAccessTokenHandler, error) {
	handler, err := kernel.NewTokenHandler(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	redBookHandler := &RedBookAccessTokenHandler{
		Config:             cfg,
		AccessTokenHandler: handler,
	}

	//if cfg.OAuthUrl != "" {
	//	redBookHandler.AccessTokenHandler.EndpointToGetToken = cfg.OAuthUrl
	//} else {
	//	redBookHandler.AccessTokenHandler.EndpointToGetToken = "https://oauth2.redBookapis.com/token"
	//}
	redBookHandler.OverrideGetCredentials()

	return redBookHandler, nil
}

func (acHandler *RedBookAccessTokenHandler) OverrideGetCredentials() {
	acHandler.AccessTokenHandler.GetCredentials = func() *object.StringMap {
		return &object.StringMap{
			"grant_type":    "authorization_code",
			"client_id":     acHandler.Config.ClientID,
			"client_secret": acHandler.Config.ClientSecret,
			//"neededText": "",
		}
	}
}
