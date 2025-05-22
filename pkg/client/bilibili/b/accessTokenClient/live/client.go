package video

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BiliBiliLiveClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliLiveClient {
	return &BiliBiliLiveClient{
		BaseClient: c,
	}
}
