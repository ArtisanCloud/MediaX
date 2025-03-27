package blogger

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/post"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/core"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type GoogleBloggerClient struct {
	Logger             *logger.Logger
	Cache              cache.ICache
	GoogleClient       *core.GoogleClient
	BloggerConfig      *config.GoogleBloggerConfig
	AccessTokenHandler *core.GoogleAccessTokenHandler

	// clients
	post *post.BloggerPostClient
}

func NewGoogleBloggerClient(cfg *config.GoogleBloggerConfig, logger *logger.Logger, cache cache.ICache) (*GoogleBloggerClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://www.googleapis.com/blogger/v3"
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

	return &GoogleBloggerClient{
		Logger:             logger,
		Cache:              cache,
		GoogleClient:       c,
		BloggerConfig:      cfg,
		AccessTokenHandler: handler,
	}, nil
}

func (client *GoogleBloggerClient) GetVideoClient() *post.BloggerPostClient {
	if client.post == nil {
		client.post = post.NewClient(client.GoogleClient.BaseClient)
	}
	return client.post
}
