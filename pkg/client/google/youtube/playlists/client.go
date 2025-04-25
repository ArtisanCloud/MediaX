package playlists

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubePlaylistsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubePlaylistsClient {
	return &YoutubePlaylistsClient{
		BaseClient: c,
	}
}
