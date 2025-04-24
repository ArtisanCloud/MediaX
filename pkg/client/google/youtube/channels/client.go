package channels

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeChannelsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeChannelsClient {
	return &YoutubeChannelsClient{
		BaseClient: c,
	}
}
