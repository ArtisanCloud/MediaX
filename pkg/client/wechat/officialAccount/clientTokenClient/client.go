package officialAccount

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/base"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/material"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/media"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/publish"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type WeChatOfficialAccountCTClient struct {
	Logger                *logger.Logger
	Cache                 cache.ICache
	WeChatClient          *core.WeChatClient
	OfficialAccountConfig *config.WeChatOfficialAccountConfig
	AccessTokenHandler    *core.WeChatAccessTokenHandler

	// clients
	base     *base.OfficialAccountBaseClient
	media    *media.OfficialAccountMediaClient
	material *material.OfficialAccountMaterialClient
	publish  *publish.OfficialAccountPublishClient
}

func NewWeChatOfficialAccountCTClient(cfg *config.WeChatOfficialAccountConfig, logger *logger.Logger, cache cache.ICache) (*WeChatOfficialAccountCTClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = config.WechatAppAPIUrl
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
	c.TokenHandler = handler.ClientTokenHandler

	return &WeChatOfficialAccountCTClient{
		Logger:                logger,
		Cache:                 cache,
		WeChatClient:          c,
		OfficialAccountConfig: cfg,
		AccessTokenHandler:    handler,
	}, nil
}

func (c *WeChatOfficialAccountCTClient) GetBaseClient() *base.OfficialAccountBaseClient {
	if c.base == nil {
		c.base = base.NewClient(c.WeChatClient.BaseClient)
	}
	return c.base
}

func (c *WeChatOfficialAccountCTClient) GetMediaClient() *media.OfficialAccountMediaClient {
	if c.media == nil {
		c.media = media.NewClient(c.WeChatClient.BaseClient)
	}
	return c.media
}

func (c *WeChatOfficialAccountCTClient) GetMaterialClient() *material.OfficialAccountMaterialClient {
	if c.material == nil {
		c.material = material.NewClient(c.WeChatClient.BaseClient)
	}
	return c.material
}

func (c *WeChatOfficialAccountCTClient) GetPublishClient() *publish.OfficialAccountPublishClient {
	if c.publish == nil {
		c.publish = publish.NewClient(c.WeChatClient.BaseClient)
	}
	return c.publish
}
