package group

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinIMGroupClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinIMGroupClient {
	return &DouYinIMGroupClient{
		BaseClient: c,
	}
}
