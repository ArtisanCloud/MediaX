package video

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BiliBiliDataArticleClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliDataArticleClient {
	return &BiliBiliDataArticleClient{
		BaseClient: c,
	}
}
