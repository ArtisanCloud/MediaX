package thumbnails

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeThumbnailsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeThumbnailsClient {
	return &YoutubeThumbnailsClient{
		BaseClient: c,
	}
}
