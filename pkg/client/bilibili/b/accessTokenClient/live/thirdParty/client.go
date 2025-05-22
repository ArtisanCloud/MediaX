package video

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type BiliBiliLiveThirdPartyClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliLiveThirdPartyClient {
	return &BiliBiliLiveThirdPartyClient{
		BaseClient: c,
	}
}
