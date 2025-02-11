package officialAccount

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type WeChatOfficialAccountClient struct {
	Logger             *logger.Logger
	Cache              cache.CacheInterface
	Client             *core.WeChatClient
	AccessTokenHandler *core.WeChatAccessTokenHandler
}

func NewWeChatOfficialAccountClient(cfg *config.WeChatOfficialAccountConfig, logger *logger.Logger, cache cache.CacheInterface) (*WeChatOfficialAccountClient, error) {
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

	return &WeChatOfficialAccountClient{
		Logger:             logger,
		Cache:              cache,
		Client:             c,
		AccessTokenHandler: handler,
	}, nil
}
