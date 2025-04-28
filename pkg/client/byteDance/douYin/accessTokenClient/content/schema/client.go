package schema

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinContentSchemaClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinContentSchemaClient {
	return &DouYinContentSchemaClient{
		BaseClient: c,
	}
}
