package members

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeMembersClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeMembersClient {
	return &YoutubeMembersClient{
		BaseClient: c,
	}
}
