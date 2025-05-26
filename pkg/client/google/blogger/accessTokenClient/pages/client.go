package pages

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type BloggerPagesClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerPagesClient {
	return &BloggerPagesClient{
		BaseClient: c,
	}
}
