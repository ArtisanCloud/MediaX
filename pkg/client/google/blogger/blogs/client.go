package blogs

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type BloggerBlogsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerBlogsClient {
	return &BloggerBlogsClient{
		BaseClient: c,
	}
}
