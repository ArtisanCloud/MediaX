package taskBox

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinTaskBoxClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinTaskBoxClient {
	return &DouYinTaskBoxClient{
		BaseClient: c,
	}
}
