package account

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type JuGuangAccountClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangAccountClient {
	return &JuGuangAccountClient{
		BaseClient: c,
	}
}

// api/open/jg/
