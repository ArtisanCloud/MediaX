package videoAbuseReportReasons

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel"
)

type YoutubeVideoAbuseReportReasonsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeVideoAbuseReportReasonsClient {
	return &YoutubeVideoAbuseReportReasonsClient{
		BaseClient: c,
	}
}
