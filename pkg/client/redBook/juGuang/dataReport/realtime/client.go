package realtime

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type JuGuangDataReportRealtimeClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *JuGuangDataReportRealtimeClient {
	return &JuGuangDataReportRealtimeClient{
		BaseClient: c,
	}
}
