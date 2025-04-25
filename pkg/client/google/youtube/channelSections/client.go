package channelSections

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeChannelSectionsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeChannelSectionsClient {
	return &YoutubeChannelSectionsClient{
		BaseClient: c,
	}
}
