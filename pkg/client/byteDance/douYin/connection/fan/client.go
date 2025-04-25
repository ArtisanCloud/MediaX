package fan

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinConnectionFanClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinConnectionFanClient {
	return &DouYinConnectionFanClient{
		BaseClient: c,
	}
}
