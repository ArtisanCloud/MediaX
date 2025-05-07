package playlistItems

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubePlaylistItemsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubePlaylistItemsClient {
	return &YoutubePlaylistItemsClient{
		BaseClient: c,
	}
}
