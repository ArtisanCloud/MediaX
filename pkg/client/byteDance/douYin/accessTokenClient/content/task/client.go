package task

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinContentTaskClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinContentTaskClient {
	return &DouYinContentTaskClient{
		BaseClient: c,
	}
}
