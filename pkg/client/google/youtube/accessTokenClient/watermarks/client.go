// Package watermarks 提供了 YouTube 水印管理相关的功能
//
// 本包实现了 YouTube Data API v3 中的水印管理功能，包括：
// - 为频道设置水印
// - 删除频道水印
//
// 使用说明：
// 1. 需要先获取有效的 YouTube API 访问令牌
// 2. 仅频道所有者可以管理水印
// 3. 每个频道只能设置一个水印
// 4. 水印图片必须符合 YouTube 的规格要求
//
// 详细文档请参考：
// https://developers.google.cn/youtube/v3/docs/watermarks?hl=zh-cn
package watermarks

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/watermarks/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeWatermarksClient 是 YouTube 水印管理 API 的客户端
// 提供设置和删除频道水印的方法
type YoutubeWatermarksClient struct {
	*kernel.BaseClient // 继承基础客户端功能
}

// NewClient 创建一个新的 YoutubeWatermarksClient 实例
// 参数:
//   - c: 基础客户端实例，用于处理 HTTP 请求和认证
//
// 返回值:
//   - *YoutubeWatermarksClient: 水印管理客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeWatermarksClient {
	return &YoutubeWatermarksClient{
		BaseClient: c,
	}
}

// Set 为 YouTube 频道设置水印
//
// 该方法用于为指定的 YouTube 频道设置水印。水印将显示在该频道所有视频的右下角。
// 注意事项：
// 1. 每个频道同一时间只能设置一个水印
// 2. 水印图片必须是透明背景的 PNG 格式
// 3. 图片尺寸建议为 150x150 像素
// 4. 图片大小不得超过 1MB
// 5. 需要频道所有者权限
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/watermarks/set?hl=zh-cn
//
// 参数：
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - data: 水印设置请求参数，包含：
//   - channelId: 要设置水印的频道ID（必填）
//   - onBehalfOfContentOwner: 内容所有者ID（可选，仅限内容合作伙伴使用）
//   - watermark: 水印配置数据（必填），包含图片URL、显示时间等
//
// 返回值：
//   - *schema.YouTubeWatermarksSetRes: 成功时返回 HTTP 204
//   - error: 设置过程中的错误信息，如权限不足、参数无效等
func (c *YoutubeWatermarksClient) Set(ctx context.Context, data *schema.YouTubeWatermarksSetReq) (*schema.YouTubeWatermarksSetRes, error) {
	result := &schema.YouTubeWatermarksSetRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/upload/youtube/v3/watermarks/set", &object.StringMap{
		"channelId":              data.ChannelId,
		"onBehalfOfContentOwner": data.OnBehalfOfContentOwner,
	}, data.Watermark, nil, result)
	return result, err
}

// ## Unset 删除 YouTube 频道的水印
//
// 该方法用于移除指定 YouTube 频道当前设置的水印。
// 注意事项：
// 1. 删除后水印将立即从所有视频中移除
// 2. 需要频道所有者权限
// 3. 删除操作不可撤销
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/watermarks/unset?hl=zh-cn
//
// 参数：
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - data: 删除水印请求参数，包含：
//   - channelId: 要删除水印的频道ID（必填）
//   - onBehalfOfContentOwner: 内容所有者ID（可选，仅限内容合作伙伴使用）
//
// 返回值：
//   - *schema.YouTubeWatermarksUnsetRes: 成功时返回 HTTP 204
//   - error: 删除过程中的错误信息，如权限不足、频道不存在等
func (c *YoutubeWatermarksClient) Unset(ctx context.Context, data *schema.YouTubeWatermarksUnsetReq) (*schema.YouTubeWatermarksUnsetRes, error) {
	result := &schema.YouTubeWatermarksUnsetRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpPost(ctx, "/youtube/v3/watermarks/unset", params, nil, nil, result)
	return result, err
}
