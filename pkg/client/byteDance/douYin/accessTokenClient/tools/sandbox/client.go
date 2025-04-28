package sandbox

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinSandboxClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinSandboxClient {
	return &DouYinSandboxClient{
		BaseClient: c,
	}
}
