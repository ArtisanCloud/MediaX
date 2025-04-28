package tool

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinIMToolClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinIMToolClient {
	return &DouYinIMToolClient{
		BaseClient: c,
	}
}
