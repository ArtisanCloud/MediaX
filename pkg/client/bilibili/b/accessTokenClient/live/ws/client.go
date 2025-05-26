package ws

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BiliBiliLiveWSClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliLiveWSClient {
	return &BiliBiliLiveWSClient{
		BaseClient: c,
	}
}
