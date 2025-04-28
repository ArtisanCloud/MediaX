package fanData

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinConnectionFanDataClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinConnectionFanDataClient {
	return &DouYinConnectionFanDataClient{
		BaseClient: c,
	}
}
