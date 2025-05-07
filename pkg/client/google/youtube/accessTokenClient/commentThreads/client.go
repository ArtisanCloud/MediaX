package commentThreads

import "github.com/ArtisanCloud/MediaX/internal/kernel"

type YoutubeCommentThreadsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeCommentThreadsClient {
	return &YoutubeCommentThreadsClient{
		BaseClient: c,
	}
}
