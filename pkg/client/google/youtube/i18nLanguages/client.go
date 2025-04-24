package i18nLanguages

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeI18nLanguagesClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeI18nLanguagesClient {
	return &YoutubeI18nLanguagesClient{
		BaseClient: c,
	}
}
