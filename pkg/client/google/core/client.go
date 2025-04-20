package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type GoogleClient struct {
	*kernel.BaseClient
	ClientConfig *config.ClientConfig
}

func NewGoogleClient(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*GoogleClient, error) {

	baseClient, err := kernel.NewBaseClient(cfg, logger, cache)
	if err != nil {
		return nil, err
	}
	youtubeClient := &GoogleClient{
		BaseClient:   baseClient,
		ClientConfig: cfg,
	}

	youtubeClient.OverrideCheckTokenNeedRefresh()

	return youtubeClient, nil
}
