package juGuang

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/core"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type RedBookJuGuangClient struct {
	RedBookClient      *core.RedBookClient
	JuGuangConfig      *config.RedBookJuGuangConfig
	AccessTokenHandler *core.RedBookAccessTokenHandler

	// clients

}

func NewRedBookJuGuangClient(cfg *config.RedBookJuGuangConfig, logger *logger.Logger, cache cache.ICache) (*RedBookJuGuangClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://adapi.xiaohongshu.com/api/open/jg/"
	}
	c, err := core.NewRedBookClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewRedBookAccessTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.AccessTokenHandler

	// override get custom token
	c.TokenHandler.GetCustomToken = cfg.GetOAuthToken

	return &RedBookJuGuangClient{
		RedBookClient:      c,
		JuGuangConfig:      cfg,
		AccessTokenHandler: handler,
	}, nil
}
