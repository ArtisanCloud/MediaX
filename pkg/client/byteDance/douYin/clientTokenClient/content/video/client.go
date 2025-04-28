package video

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/video/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type DouYinContentVideoClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinContentVideoClient {
	return &DouYinContentVideoClient{
		BaseClient: c,
	}
}

// 查询视频发布结果
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/douyin/search-video/video-share-result
func (comp *DouYinContentVideoClient) ShareResult(ctx context.Context, data *schema.DouYinContentVideoShareResultReq) (*schema.DouYinContentVideoShareResultRes, error) {
	result := &schema.DouYinContentVideoShareResultRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = comp.BaseClient.HttpGet(ctx, "share-id/", params, nil, result)
	return result, err
}

// 查询视频携带的地点信息
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/douyin/search-video/video-poi
func (comp *DouYinContentVideoClient) PoiSearch(ctx context.Context, data *schema.DouYinContentVideoPoiSearchReq) (*schema.DouYinContentVideoPoiSearchRes, error) {
	result := &schema.DouYinContentVideoPoiSearchRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = comp.BaseClient.HttpGet(ctx, "poi/search/keyword/", params, nil, result)
	return result, err
}
