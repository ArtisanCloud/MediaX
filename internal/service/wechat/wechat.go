package wechat

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type WeChatService struct {
	Client *WechatClient
}

func NewWeChatService(cfg *config.WeChatConfig, logger *logger.Logger, cache cache.CacheInterface) (*WeChatService, error) {

	return &WeChatService{
		Client: NewWechatClient(cfg, logger, cache),
	}, nil
}
