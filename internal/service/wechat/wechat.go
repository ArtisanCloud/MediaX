package wechat

import (
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/core"
	"github.com/ArtisanCloud/MediaX/v1/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type WeChatService struct {
	Logger             *logger.Logger
	Cache              cache.CacheInterface
	Client             *core.WeChatClient
	AccessTokenHandler *core.WeChatAccessTokenHandler
}

func NewWeChatService(cfg *config.WeChatConfig, logger *logger.Logger, cache cache.CacheInterface) (*WeChatService, error) {
	c, err := core.NewWeChatClient(cfg, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewWeChatAccessTokenHandler(cfg, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.AccessTokenHandler

	return &WeChatService{
		Logger:             logger,
		Cache:              cache,
		Client:             c,
		AccessTokenHandler: handler,
	}, nil
}
