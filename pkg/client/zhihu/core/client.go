package core

import (
	"errors"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// CoreClient 封装知乎 SessionToken 相关的基础能力（BaseClient、配置、日志等）。
type CoreClient struct {
	BaseClient *kernel.BaseClient
	Config     *config.ZhihuSessionTokenConfig
	Logger     *logger.Logger
	Cache      cache.ICache
}

// NewClient 根据配置创建知乎 SessionToken 核心客户端。
func NewClient(cfg *config.ZhihuSessionTokenConfig, log *logger.Logger, cache cache.ICache) (*CoreClient, error) {
	if cfg == nil {
		return nil, errors.New("zhihu.sessiontoken: config is nil")
	}
	if cfg.Service.BaseURL == "" {
		return nil, errors.New("zhihu.sessiontoken: service.base_url is empty")
	}
	timeout := cfg.Service.Timeout
	if timeout <= 0 {
		timeout = 30
	}
	baseCfg := &config.ClientConfig{
		BaseConfig: &config.BaseConfig{
			ApiUrl:    cfg.Service.BaseURL,
			Timeout:   float64(timeout),
			HttpDebug: cfg.Service.HTTPDebug,
		},
	}
	baseClient, err := kernel.NewBaseClient(baseCfg, log, cache)
	if err != nil {
		return nil, err
	}
	return &CoreClient{
		BaseClient: baseClient,
		Config:     cfg,
		Logger:     log,
		Cache:      cache,
	}, nil
}
