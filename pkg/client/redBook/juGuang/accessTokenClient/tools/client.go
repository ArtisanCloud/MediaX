package tools

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type JuGuangToolClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangToolClient {
	return &JuGuangToolClient{
		BaseClient: c,
	}
}
