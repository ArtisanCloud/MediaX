package micApp

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinMicAppClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinMicAppClient {
	return &DouYinMicAppClient{
		BaseClient: c,
	}
}
