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
func (c *YoutubeVideoClient) List(ctx context.Context, data *schema.YouTubeVideoListReq) (*schema.YouTubeVideoListRes, error) {
	result := &schema.YouTubeVideoListRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpGet(ctx, "youtube/v3/videos", params, nil, result)
	return result, err
}

// Videos: insert 将视频上传到 YouTube，并可选择设置视频的元数据
// https://developers.google.com/youtube/v3/docs/videos/insert
func (c *YoutubeVideoClient) Insert(ctx context.Context, data *schema.YouTubeVideoInsertReq) (*schema.YouTubeVideoInsertRes, error) {
	result := &schema.YouTubeVideoInsertRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	// 假设 BaseClient 有 HttpPost 方法用于处理上传请求
	// 这里需要根据实际情况处理媒体上传
	_, err = c.BaseClient.HttpPost(ctx, "youtube/v3/video", params, nil, nil, result)
	return result, err
}

// Videos: update 更新视频的元数据
// https://developers.google.com/youtube/v3/docs/videos/update
func (c *YoutubeVideoClient) Update(ctx context.Context, data *schema.YouTubeVideoUpdateReq) (*schema.YouTubeVideoUpdateRes, error) {
	result := &schema.YouTubeVideoUpdateRes{}

	_, err := c.BaseClient.HttpPut(ctx, "youtube/v3/video", nil, data, nil, result)
	return result, err
}

// Videos: delete 删除 YouTube 视频
// https://developers.google.com/youtube/v3/docs/videos/delete
func (c *YoutubeVideoClient) Delete(ctx context.Context, data *schema.YouTubeVideoDeleteReq) (*schema.YouTubeVideoDeleteRes, error) {
	result := &schema.YouTubeVideoDeleteRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpDelete(ctx, "youtube/v3/video", params, nil, nil, result)
	return result, err
}

// Videos:Rate 为视频添加“顶”或“踩”评分，或者删除视频的评分。
// https://developers.google.com/youtube/v3/docs/videos/rate?hl=zh-cn
func (c *YoutubeVideoClient) Rate(ctx context.Context, data *schema.YouTubeVideoRateReq) error {
	params, err := object.StructToStringMap(data)
	if err != nil {
		return err
	}

	_, err = c.BaseClient.HttpPost(ctx, "videos/rate", params, nil, nil, nil)
	return err
}

// Videos: getRating 检索授权用户对指定视频列表的评分。
// https://developers.google.com/youtube/v3/docs/videos/getRating?hl=zh-cn
func (c *YoutubeVideoClient) GetRating(ctx context.Context, data *schema.YouTubeVideoGetRatingReq) (*schema.YouTubeVideoGetRatingRes, error) {
	result := &schema.YouTubeVideoGetRatingRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpPost(ctx, "videos/getRating", params, nil, nil, result)
	return result, err
}

// Videos:reportAbuse 举报包含侮辱性内容的视频。
// https://developers.google.com/youtube/v3/docs/videos/reportAbuse?hl=zh-cn
func (c *YoutubeVideoClient) ReportAbuse(ctx context.Context, data *schema.YouTubeVideoReportAbuseReq) error {
	params, err := object.StructToStringMap(data)
	if err != nil {
		return err
	}

	_, err = c.BaseClient.HttpPost(ctx, "videos/reportAbuse", params, nil, nil, nil)
	return err
}
