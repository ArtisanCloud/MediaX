package core

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type RedBookClient struct {
	*kernel.BaseClient
	ClientConfig *config.ClientConfig
}

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
