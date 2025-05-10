package comments

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/comments/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type YoutubeCommentsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeCommentsClient {
	return &YoutubeCommentsClient{
		BaseClient: c,
	}
}

// Comments: list
// https://developers.google.cn/youtube/v3/docs/comments/list?hl=zh-cn
func (c *YoutubeCommentsClient) List(ctx context.Context, data *schema.YoutubeCommentsListReq) (*schema.YoutubeCommentsListRes, error) {
	result := &schema.YoutubeCommentsListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/comments", params, nil, result)
	return result, err
}

// Comments: insert
// https://developers.google.cn/youtube/v3/docs/comments/insert?hl=zh-cn
func (c *YoutubeCommentsClient) Insert(ctx context.Context, data *schema.YoutubeCommentsInsertReq) (*schema.YoutubeCommentsInsertRes, error) {
	result := &schema.YoutubeCommentsInsertRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/youtube/v3/comments", data, nil, result)
	return result, err
}

// Comments: update - 更新评论
// https://developers.google.cn/youtube/v3/docs/comments/update?hl=zh-cn
func (c *YoutubeCommentsClient) Update(ctx context.Context, data *schema.YoutubeCommentsUpdateReq) (*schema.YoutubeCommentsUpdateRes, error) {
	result := &schema.YoutubeCommentsUpdateRes{}
	_, err := c.BaseClient.HttpPut(ctx, "/youtube/v3/comments", data, nil, result)
	return result, err
}
