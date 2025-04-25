package channelBanners

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeChannelBannersClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeChannelBannersClient {
	return &YoutubeChannelBannersClient{
		BaseClient: c,
	}
}
