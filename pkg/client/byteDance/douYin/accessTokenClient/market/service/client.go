package service

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinMarketServiceClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinMarketServiceClient {
	return &DouYinMarketServiceClient{
		BaseClient: c,
	}
}
