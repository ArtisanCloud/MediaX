package micApp

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinToolMicAppClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinToolMicAppClient {
	return &DouYinToolMicAppClient{
		BaseClient: c,
	}
}
