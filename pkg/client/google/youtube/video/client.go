package video

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/video/schema"
)

type YoutubeVideoClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeVideoClient {
	return &YoutubeVideoClient{
		BaseClient: c,
	}
}

// Videos: list 返回与 API 请求参数匹配的视频列表
// https://developers.google.cn/youtube/v3/docs/videos/list?hl=zh-cn
func (comp *YoutubeVideoClient) List(ctx context.Context, data *schema.ListReq) (*schema.ListRes, error) {
	result := &schema.ListRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "videos", data, nil, result)
	return result, err
}
