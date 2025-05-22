package video

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BiliBiliDataVideoClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliDataVideoClient {
	return &BiliBiliDataVideoClient{
		BaseClient: c,
	}
}
