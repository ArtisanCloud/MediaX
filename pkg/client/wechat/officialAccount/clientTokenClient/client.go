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

// WeChatOfficialAccountCTClient 微信公众号客户端Token客户端
// 提供微信公众号平台的各种功能接口
type WeChatOfficialAccountCTClient struct {
	Logger                *logger.Logger                      // 日志记录器
	Cache                 cache.ICache                        // 缓存接口
	WeChatClient          *core.WeChatClient                  // 微信基础客户端
	OfficialAccountConfig *config.WeChatOfficialAccountConfig // 微信公众号配置
	AccessTokenHandler    *core.WeChatAccessTokenHandler      // 访问Token处理器

	// clients
	base     *base.OfficialAccountBaseClient         // 基础功能客户端
	media    *media.OfficialAccountMediaClient       // 媒体管理客户端
	material *material.OfficialAccountMaterialClient // 素材管理客户端
	publish  *publish.OfficialAccountPublishClient   // 发布管理客户端
}

// NewWeChatOfficialAccountCTClient 创建新的微信公众号客户端Token客户端实例
// cfg: 微信公众号配置
// logger: 日志记录器
// cache: 缓存接口
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

// GetBaseClient 获取微信公众号基础功能客户端
func (c *WeChatOfficialAccountCTClient) GetBaseClient() *base.OfficialAccountBaseClient {
	if c.base == nil {
		c.base = base.NewClient(c.WeChatClient.BaseClient)
	}
	return c.base
}

// GetMediaClient 获取微信公众号媒体管理客户端
func (c *WeChatOfficialAccountCTClient) GetMediaClient() *media.OfficialAccountMediaClient {
	if c.media == nil {
		c.media = media.NewClient(c.WeChatClient.BaseClient)
	}
	return c.media
}

// GetMaterialClient 获取微信公众号素材管理客户端
func (c *WeChatOfficialAccountCTClient) GetMaterialClient() *material.OfficialAccountMaterialClient {
	if c.material == nil {
		c.material = material.NewClient(c.WeChatClient.BaseClient)
	}
	return c.material
}

// GetPublishClient 获取微信公众号发布管理客户端
func (c *WeChatOfficialAccountCTClient) GetPublishClient() *publish.OfficialAccountPublishClient {
	if c.publish == nil {
		c.publish = publish.NewClient(c.WeChatClient.BaseClient)
	}
	return c.publish
}
