package users

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type BloggerUsersClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerUsersClient {
	return &BloggerUsersClient{
		BaseClient: c,
	}
}
