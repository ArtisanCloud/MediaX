package blogUserInfos

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type BloggerBlogUserInfosClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerBlogUserInfosClient {
	return &BloggerBlogUserInfosClient{
		BaseClient: c,
	}
}
