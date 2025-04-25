package offline

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type JuGuangDataReportOfflineClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangDataReportOfflineClient {
	return &JuGuangDataReportOfflineClient{
		BaseClient: c,
	}
}
