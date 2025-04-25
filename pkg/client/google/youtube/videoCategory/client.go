package videoCategory

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/videoCategory/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type YoutubeVideoCategoryClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeVideoCategoryClient {
	return &YoutubeVideoCategoryClient{
		BaseClient: c,
	}
}

// VideoCategories:list 返回可与 YouTube 视频相关联的类别列表。
// https://developers.google.com/youtube/v3/docs/videoCategories/list?hl=zh-cn
func (comp *YoutubeVideoCategoryClient) List(ctx context.Context, data *schema.YouTubeVideoCategoriesReq) (*schema.YouTubeVideoCategoriesRes, error) {
	result := &schema.YouTubeVideoCategoriesRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = comp.BaseClient.HttpGet(ctx, "videoCategories", params, nil, result)
	return result, err
}
