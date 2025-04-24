package comments

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeCommentsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeCommentsClient {
	return &YoutubeCommentsClient{
		BaseClient: c,
	}
}
