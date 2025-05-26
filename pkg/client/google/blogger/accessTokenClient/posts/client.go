package posts

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type BloggerPostsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerPostsClient {
	return &BloggerPostsClient{
		BaseClient: c,
	}
}
