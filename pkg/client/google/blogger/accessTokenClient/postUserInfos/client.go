package postUserInfos

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BloggerPostUserInfosClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerPostUserInfosClient {
	return &BloggerPostUserInfosClient{
		BaseClient: c,
	}
}
