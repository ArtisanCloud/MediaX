package video

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BiliBiliVideoClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliVideoClient {
	return &BiliBiliVideoClient{
		BaseClient: c,
	}
}
