package data

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinConnectionDataClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinConnectionDataClient {
	return &DouYinConnectionDataClient{
		BaseClient: c,
	}
}
