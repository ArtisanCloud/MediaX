package media

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinToolMediaClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinToolMediaClient {
	return &DouYinToolMediaClient{
		BaseClient: c,
	}
}
