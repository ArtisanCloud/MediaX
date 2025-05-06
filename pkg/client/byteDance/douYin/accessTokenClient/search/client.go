package search

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinSearchClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinSearchClient {
	return &DouYinSearchClient{
		BaseClient: c,
	}
}
