package ticket

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type DouYinTicketClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinTicketClient {
	return &DouYinTicketClient{
		BaseClient: c,
	}
}
