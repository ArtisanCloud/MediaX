package youtube

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/video"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type GoogleYouTubeClient struct {
	Logger             *logger.Logger
	Cache              cache.ICache
	GoogleClient       *core.GoogleClient
	YouTubeConfig      *config.GoogleYouTubeConfig
	AccessTokenHandler *core.GoogleAccessTokenHandler

	// clients
	video *video.YoutubeVideoClient
}

func NewGoogleYouTubeClient(cfg *config.GoogleYouTubeConfig, logger *logger.Logger, cache cache.ICache) (*GoogleYouTubeClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://www.googleapis.com/youtube/v3/"
	}
	c, err := core.NewGoogleClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewGoogleAccessTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.AccessTokenHandler

	// override get custom token
	c.TokenHandler.GetCustomToken = cfg.GetOAuthToken

	return &GoogleYouTubeClient{
		Logger:             logger,
		Cache:              cache,
		GoogleClient:       c,
		YouTubeConfig:      cfg,
		AccessTokenHandler: handler,
	}, nil
}

func (client *GoogleYouTubeClient) GetVideoClient() *video.YoutubeVideoClient {
	if client.video == nil {
		client.video = video.NewClient(client.GoogleClient.BaseClient)
	}
	return client.video
}
