package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// RedBookClient 小红书客户端结构体
type RedBookClient struct {
	*kernel.BaseClient
	ClientConfig *config.ClientConfig
}

// NewRedBookClient 创建一个新的小红书客户端实例
func NewRedBookClient(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*RedBookClient, error) {
	baseClient, err := kernel.NewBaseClient(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	redBookClient := &RedBookClient{
		BaseClient:   baseClient,
		ClientConfig: cfg,
	}

	return redBookClient, nil
}
