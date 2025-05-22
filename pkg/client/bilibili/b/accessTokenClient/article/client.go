package video

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BiliBiliArticleClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliArticleClient {
	return &BiliBiliArticleClient{
		BaseClient: c,
	}
}
