package post

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/blogger/post/schema"
)

type BloggerPostClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BloggerPostClient {
	return &BloggerPostClient{
		BaseClient: c,
	}
}

// Posts: list 返回与 API 请求参数匹配的视频列表
// https://developers.google.cn/youtube/v3/docs/posts/list?hl=zh-cn
func (comp *BloggerPostClient) List(ctx context.Context, data *schema.ListReq) (*schema.ListRes, error) {
	result := &schema.ListRes{}

	_, err := comp.BaseClient.HttpPost(ctx, "posts", data, nil, result)
	return result, err
}
