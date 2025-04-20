package video

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/video/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
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
func (comp *YoutubeVideoClient) List(ctx context.Context, data *schema.YouTubeVideoListReq) (*schema.YouTubeVideoListRes, error) {
	result := &schema.YouTubeVideoListRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = comp.BaseClient.HttpGet(ctx, "videos", params, nil, result)
	return result, err
}

// Videos: insert 将视频上传到 YouTube，并可选择设置视频的元数据
// https://developers.google.com/youtube/v3/docs/videos/insert
func (comp *YoutubeVideoClient) Insert(ctx context.Context, data *schema.YouTubeVideoInsertReq) (*schema.YouTubeVideoInsertRes, error) {
	result := &schema.YouTubeVideoInsertRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	// 假设 BaseClient 有 HttpPost 方法用于处理上传请求
	// 这里需要根据实际情况处理媒体上传
	_, err = comp.BaseClient.HttpPost(ctx, "https://www.googleapis.com/upload/youtube/v3/videos", params, nil, result)
	return result, err
}

// Videos: update 更新视频的元数据
// https://developers.google.com/youtube/v3/docs/videos/update
func (comp *YoutubeVideoClient) Update(ctx context.Context, data *schema.YouTubeVideoUpdateReq) (*schema.YouTubeVideoUpdateRes, error) {
	result := &schema.YouTubeVideoUpdateRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = comp.BaseClient.HttpPut(ctx, "videos", params, nil, result)
	return result, err
}

// Videos: delete 删除 YouTube 视频
// https://developers.google.com/youtube/v3/docs/videos/delete
func (comp *YoutubeVideoClient) Delete(ctx context.Context, data *schema.YouTubeVideoDeleteReq) (*schema.YouTubeVideoDeleteRes, error) {
	result := &schema.YouTubeVideoDeleteRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = comp.BaseClient.HttpDelete(ctx, "videos", params, nil, result)
	return result, err
}
