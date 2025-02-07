package wechat

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type WeChatClient struct {
	kernel.BaseClient
	Config *config.WeChatConfig
}

func NewWeChatClient(cfg *config.WeChatConfig, logger *logger.Logger, cache cache.CacheInterface) (*WechatClient, error) {
	baseClient := &kernel.BaseClient{
		Config: cfg,
		Logger: logger,
		Cache:  cache,
	}
	return &WechatClient{
		BaseClient: *baseClient,
	}, nil
}
