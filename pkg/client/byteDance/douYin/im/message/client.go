package message

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinIMMessageClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinIMMessageClient {
	return &DouYinIMMessageClient{
		BaseClient: c,
	}
}
