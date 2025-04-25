package subscriptions

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeSubscriptionsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeSubscriptionsClient {
	return &YoutubeSubscriptionsClient{
		BaseClient: c,
	}
}
