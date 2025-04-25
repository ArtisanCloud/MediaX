package watermarks

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeWatermarksClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeWatermarksClient {
	return &YoutubeWatermarksClient{
		BaseClient: c,
	}
}
