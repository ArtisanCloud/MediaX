package note

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type JuGuangNoteClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangNoteClient {
	return &JuGuangNoteClient{
		BaseClient: c,
	}
}
