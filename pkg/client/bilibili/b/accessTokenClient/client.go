package accessTokenClient

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/article"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/data"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/live"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/live/thirdParty"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/live/ws"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/user"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/video"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// BiliBiliACClient 是B站访问令牌客户端
// 用于管理与B站API交互所需的访问令牌和相关配置
// 示例用法：
//
//	cfg := &config.BiliBiliConfig{...}
//	client, err := NewBiliBiliACClient(cfg, logger, cache)
type BiliBiliACClient struct {
	// BiliBiliClient B站基础客户端，提供核心API调用功能
	BiliBiliClient *core.BiliBiliClient

	// BiliBiliConfig BiliBili相关配置，用于跨平台功能集成
	BiliBiliConfig *config.BiliBiliConfig

	// AccessTokenHandler 访问令牌处理器，负责令牌的获取和刷新
	AccessTokenHandler *core.BiliBiliAccessTokenHandler

	// clients 其他客户端实例，可根据需要扩展
	video          *video.BiliBiliVideoClient               // 视频客户端
	user           *user.BiliBiliUserClient                 // 用户客户端
	article        *article.BiliBiliArticleClient           // 文章客户端
	data           *data.BiliBiliDataClient                 // 数据文章客户端
	live           *live.BiliBiliLiveClient                 // 直播客户端
	liveWS         *ws.BiliBiliLiveWSClient                 // 直播WS客户端
	liveThirdParty *thirdParty.BiliBiliLiveThirdPartyClient // 直播第三方客户端
}

func NewBiliBiliACClient(cfg *config.BiliBiliConfig, logger *logger.Logger, cache cache.ICache) (*BiliBiliACClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = config.BiliBiliAPIUrl
	}
	c, err := core.NewBiliBiliClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}
	handler, err := core.NewBiliBiliAccessTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}
	c.TokenHandler = handler.AccessTokenHandler
	c.TokenHandler.GetCustomToken = cfg.GetOAuthToken

	return &BiliBiliACClient{
		BiliBiliClient:     c,
		BiliBiliConfig:     cfg,
		AccessTokenHandler: handler,
	}, nil
}

// GetVideoClient 获取视频客户端
func (client *BiliBiliACClient) GetVideoClient() *video.BiliBiliVideoClient {
	if client.video == nil {
		client.video = video.NewClient(client.BiliBiliClient.BaseClient)
	}
	return client.video
}

// GetUserClient 获取用户客户端
func (client *BiliBiliACClient) GetUserClient() *user.BiliBiliUserClient {
	if client.user == nil {
		client.user = user.NewClient(client.BiliBiliClient.BaseClient)
	}
	return client.user
}

// GetArticleClient 获取文章客户端
func (client *BiliBiliACClient) GetArticleClient() *article.BiliBiliArticleClient {
	if client.article == nil {
		client.article = article.NewClient(client.BiliBiliClient.BaseClient)
	}
	return client.article
}

// GetDataClient 获取数据文章客户端
func (client *BiliBiliACClient) GetDataClient() *data.BiliBiliDataClient {
	if client.data == nil {
		client.data = data.NewClient(client.BiliBiliClient.BaseClient)
	}
	return client.data
}

// GetLiveClient 获取直播客户端
func (client *BiliBiliACClient) GetLiveClient() *live.BiliBiliLiveClient {
	if client.live == nil {
		client.live = live.NewClient(client.BiliBiliClient.BaseClient)
	}
	return client.live
}

// GetLiveWSClient 获取直播WS客户端
func (client *BiliBiliACClient) GetLiveWSClient() *ws.BiliBiliLiveWSClient {
	if client.liveWS == nil {
		client.liveWS = ws.NewClient(client.BiliBiliClient.BaseClient)
	}
	return client.liveWS
}

// GetLiveThirdPartyClient 获取直播第三方客户端
func (client *BiliBiliACClient) GetLiveThirdPartyClient() *thirdParty.BiliBiliLiveThirdPartyClient {
	if client.liveThirdParty == nil {
		client.liveThirdParty = thirdParty.NewClient(client.BiliBiliClient.BaseClient)
	}
	return client.liveThirdParty
}
