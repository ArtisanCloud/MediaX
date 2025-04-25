package video

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/video/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

const videoUri = "api/douyin/v1/video/"

type DouYinContentVideoClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinContentVideoClient {
	return &DouYinContentVideoClient{
		BaseClient: c,
	}
}

// 查询授权账号视频列表
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/douyin/search-video/account-video-list
func (comp *DouYinContentVideoClient) List(ctx context.Context, data *schema.DouYinContentVideoListReq) (*schema.DouYinContentVideoListRes, error) {
	result := &schema.DouYinContentVideoListRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = comp.BaseClient.HttpGet(ctx, videoUri+"video_list/", params, nil, result)
	return result, err
}

// 查询特定视频的视频数据
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/douyin/search-video/video-data
func (comp *DouYinContentVideoClient) Data(ctx context.Context, data *schema.DouYinContentVideoDataReq) (*schema.DouYinContentVideoDataRes, error) {
	result := &schema.DouYinContentVideoDataRes{}

	_, err := comp.BaseClient.HttpPost(ctx, videoUri+"video_data/", nil, data, nil, result)
	return result, err
}
