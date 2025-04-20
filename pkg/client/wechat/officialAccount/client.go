package officialAccount

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/material"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/media"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/publish"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type WeChatOfficialAccountClient struct {
	Logger                *logger.Logger
	Cache                 cache.ICache
	WeChatClient          *core.WeChatClient
	OfficialAccountConfig *config.WeChatOfficialAccountConfig
	AccessTokenHandler    *core.WeChatAccessTokenHandler

	// clients
	media    *media.OfficialAccountMediaClient
	material *material.OfficialAccountMaterialClient
	publish  *publish.OfficialAccountPublishClient
}

func NewWeChatOfficialAccountClient(cfg *config.WeChatOfficialAccountConfig, logger *logger.Logger, cache cache.ICache) (*WeChatOfficialAccountClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://api.weixin.qq.com"
	}
	c, err := core.NewWeChatClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewWeChatAccessTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.AccessTokenHandler

	return &WeChatOfficialAccountClient{
		Logger:                logger,
		Cache:                 cache,
		WeChatClient:          c,
		OfficialAccountConfig: cfg,
		AccessTokenHandler:    handler,
	}, nil
}

func (c *WeChatOfficialAccountClient) GetMediaClient() *media.OfficialAccountMediaClient {
	if c.media == nil {
		c.media = media.NewClient(c.WeChatClient.BaseClient)
	}
	return c.media
}

func (c *WeChatOfficialAccountClient) GetMaterialClient() *material.OfficialAccountMaterialClient {
	if c.material == nil {
		c.material = material.NewClient(c.WeChatClient.BaseClient)
	}
	return c.material
}

func (c *WeChatOfficialAccountClient) GetPublishClient() *publish.OfficialAccountPublishClient {
	if c.publish == nil {
		c.publish = publish.NewClient(c.WeChatClient.BaseClient)
	}
	return c.publish
}
