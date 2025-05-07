package membershipsLevels

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeMembershipsLevelsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeMembershipsLevelsClient {
	return &YoutubeMembershipsLevelsClient{
		BaseClient: c,
	}
}
