package campaign

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type JuGuangPromoteCampaignClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangPromoteCampaignClient {
	return &JuGuangPromoteCampaignClient{
		BaseClient: c,
	}
}
