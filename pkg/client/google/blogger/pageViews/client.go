package pageViews

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type BloggerPageViewsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerPageViewsClient {
	return &BloggerPageViewsClient{
		BaseClient: c,
	}
}
