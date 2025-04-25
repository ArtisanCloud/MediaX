package unit

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type JuGuangPromoteUnitClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangPromoteUnitClient {
	return &JuGuangPromoteUnitClient{
		BaseClient: c,
	}
}
