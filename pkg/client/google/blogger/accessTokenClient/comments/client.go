package comments

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type BloggerCommentsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerCommentsClient {
	return &BloggerCommentsClient{
		BaseClient: c,
	}
}
