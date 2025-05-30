package activities

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type YoutubeActivitiesClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeActivitiesClient {
	return &YoutubeActivitiesClient{
		BaseClient: c,
	}
}
