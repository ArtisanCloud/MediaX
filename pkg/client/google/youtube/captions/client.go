package captions

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeCaptionsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeCaptionsClient {
	return &YoutubeCaptionsClient{
		BaseClient: c,
	}
}
