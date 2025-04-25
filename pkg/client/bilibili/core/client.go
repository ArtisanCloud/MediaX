package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type BiliBiliClient struct {
	*kernel.BaseClient
	ClientConfig *config.ClientConfig
}

func NewBiliBiliClient(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*BiliBiliClient, error) {

	baseClient, err := kernel.NewBaseClient(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	youtubeClient := &BiliBiliClient{
		BaseClient:   baseClient,
		ClientConfig: cfg,
	}

	youtubeClient.OverrideCheckTokenNeedRefresh()

	return youtubeClient, nil
}
