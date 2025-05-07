package i18nRegions

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeI18nRegionsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeI18nRegionsClient {
	return &YoutubeI18nRegionsClient{
		BaseClient: c,
	}
}
