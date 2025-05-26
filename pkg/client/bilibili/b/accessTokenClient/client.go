package accessTokenClient

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/article"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/data"
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
	liveWS         *ws.BiliBiliLiveWSClient                 // 直播客户端
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
