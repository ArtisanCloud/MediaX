package search

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/search/schema"
)

type DouYinSearchClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinSearchClient {
	return &DouYinSearchClient{
		BaseClient: c,
	}
}

// ## SearchVideo 抖音视频搜索
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/douyin-search-capability/aweme-dy-video-search
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • Keyword: 搜索关键词
//	       • Count: 返回数量
//	       • Cursor: 游标，用于分页
//	       • DeviceId: 设备ID
//
// 返回值：
//
//	*schema.DouYinSearchVideoRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - VideoList: 视频列表 []schema.Video
//	      - SearchId: 搜索ID
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinSearchClient) SearchVideo(ctx context.Context, data *schema.DouYinSearchVideoReq) (*schema.DouYinSearchVideoRes, error) {
	result := &schema.DouYinSearchVideoRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/dy_open_api/v1/search/video/", nil, data, nil, result)
	return result, err
}
