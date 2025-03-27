package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type GoogleAccessTokenHandler struct {
	Config             *config.ClientConfig
	AccessTokenHandler *kernel.AccessTokenHandler
}

func NewGoogleAccessTokenHandler(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*GoogleAccessTokenHandler, error) {
	handler, err := kernel.NewAccessTokenHandler(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	googleHandler := &GoogleAccessTokenHandler{
		Config:             cfg,
		AccessTokenHandler: handler,
	}

	if cfg.OAuthUrl != "" {
		googleHandler.AccessTokenHandler.EndpointToGetToken = cfg.OAuthUrl
	} else {
		googleHandler.AccessTokenHandler.EndpointToGetToken = "https://oauth2.googleapis.com/token"
	}
	googleHandler.OverrideGetCredentials()

	return googleHandler, nil
}

func (acHandler *GoogleAccessTokenHandler) OverrideGetCredentials() {

	acHandler.AccessTokenHandler.GetCredentials = func() *object.StringMap {
		return &object.StringMap{
			"grant_type":    "authorization_code",
			"client_id":     acHandler.Config.ClientID,
			"client_secret": acHandler.Config.ClientSecret,
			//"neededText": "",
		}
	}
}
