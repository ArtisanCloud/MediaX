package search

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeSearchClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeSearchClient {
	return &YoutubeSearchClient{
		BaseClient: c,
	}
}
