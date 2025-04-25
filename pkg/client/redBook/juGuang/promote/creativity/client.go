package creativity

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type JuGuangPromoteCreativityClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangPromoteCreativityClient {
	return &JuGuangPromoteCreativityClient{
		BaseClient: c,
	}
}
