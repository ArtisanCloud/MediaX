package activity

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinContentActivityClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinContentActivityClient {
	return &DouYinContentActivityClient{
		BaseClient: c,
	}
}
