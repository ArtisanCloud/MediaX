package oauth

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinOAuthClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinOAuthClient {
	return &DouYinOAuthClient{
		BaseClient: c,
	}
}
