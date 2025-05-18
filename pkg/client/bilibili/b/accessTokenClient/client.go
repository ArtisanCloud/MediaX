package accessTokenClient

import (
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

	// YouTubeConfig YouTube相关配置，用于跨平台功能集成
	YouTubeConfig *config.BiliBiliConfig

	// AccessTokenHandler 访问令牌处理器，负责令牌的获取和刷新
	AccessTokenHandler *core.BiliBiliAccessTokenHandler

	// clients 其他客户端实例，可根据需要扩展
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
		YouTubeConfig:      cfg,
		AccessTokenHandler: handler,
	}, nil
}
