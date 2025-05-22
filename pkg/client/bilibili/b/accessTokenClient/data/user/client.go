package video

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BiliBiliDataUserClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliDataUserClient {
	return &BiliBiliDataUserClient{
		BaseClient: c,
	}
}
