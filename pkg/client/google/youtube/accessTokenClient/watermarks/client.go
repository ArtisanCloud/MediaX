package watermarks

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/watermarks/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeWatermarksClient 是 YouTube 视频相关 API 的客户端
type YoutubeWatermarksClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 YoutubeWatermarksClient 实例
func NewClient(c *kernel.BaseClient) *YoutubeWatermarksClient {
	return &YoutubeWatermarksClient{
		BaseClient: c,
	}
}

// ## Set 设置水印
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/watermarks/set?hl=zh-cn
//
// 参数：
//   ctx  - 请求上下文
//   data - 请求参数，包含以下字段：
//     • channelId: 频道ID（必填）
//     • onBehalfOfContentOwner: 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
//     • watermark: 水印数据（必填）
//
// 返回值：
//   *schema.YouTubeWatermarksSetRes HTTP 204 返回码
//   error 调用过程中遇到的错误（如有）
func (c *YoutubeWatermarksClient) Set(ctx context.Context, data *schema.YouTubeWatermarksSetReq) (*schema.YouTubeWatermarksSetRes, error) {
	result := &schema.YouTubeWatermarksSetRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/upload/youtube/v3/watermarks/set", &object.StringMap{
		"channelId":              data.ChannelId,
		"onBehalfOfContentOwner": data.OnBehalfOfContentOwner,
	}, data.Watermark, nil, result)
	return result, err
}

// ## Unset 删除水印
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/watermarks/unset?hl=zh-cn
//
// 参数：
//   ctx  - 请求上下文
//   data - 请求参数，包含以下字段：
//     • channelId: 频道ID（必填）
//     • onBehalfOfContentOwner: 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
//
// 返回值：
//   *schema.YouTubeWatermarksUnsetRes HTTP 204 返回码
//   error 调用过程中遇到的错误（如有）
func (c *YoutubeWatermarksClient) Unset(ctx context.Context, data *schema.YouTubeWatermarksUnsetReq) (*schema.YouTubeWatermarksUnsetRes, error) {
	result := &schema.YouTubeWatermarksUnsetRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpPost(ctx, "/youtube/v3/watermarks/unset", params, nil, nil, result)
	return result, err
}
