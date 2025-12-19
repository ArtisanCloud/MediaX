package channels

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v3/channels/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeChannelsClient 是YouTube频道API的客户端
// 继承自BaseClient，提供对YouTube频道API的访问功能
type YoutubeChannelsClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的YoutubeChannelsClient实例
// 参数：
//
//	c - BaseClient实例，提供HTTP请求、日志记录等基础功能
//
// 返回值：
//
//	*YoutubeChannelsClient - 返回初始化后的客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeChannelsClient {
	return &YoutubeChannelsClient{
		BaseClient: c,
	}
}

// ## Channels: list
// https://developers.google.cn/youtube/v3/docs/channels/list?hl=zh-cn
// List 获取YouTube频道列表
// 参数：
//
//	ctx - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填）
//	  • id: 频道ID列表（可选）
//	  • mine: 是否获取当前用户的频道（可选）
//	  • maxResults: 返回的最大结果数（可选）
//	  • pageToken: 分页令牌（可选）
//
// 返回值：
//
//	*schema.YoutubeChannelsListRes - 包含频道列表数据
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • NextPageToken: 下一页令牌
//	  • PrevPageToken: 上一页令牌
//	  • PageInfo: 分页信息
//	  • Items: 频道列表数据
//
//	error - 调用过程中遇到的错误（如有）
func (c *YoutubeChannelsClient) List(ctx context.Context, data *schema.YoutubeChannelsListReq) (*schema.YoutubeChannelsListRes, error) {
	result := &schema.YoutubeChannelsListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/channels", params, nil, nil, result)
	return result, err
}

// ## Channels: update - 更新频道
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/channels/update?hl=zh-cn
// 功能：更新YouTube频道的元数据
// 注意：此方法仅支持对channel资源的brandingSettings、invideoPromotion和localizations对象及其子属性进行更新
// 配额影响：每次调用消耗50个配额单位
// 授权范围：
//   - https://www.googleapis.com/auth/youtubepartner
//   - https://www.googleapis.com/auth/youtube
//   - https://www.googleapis.com/auth/youtube.force-ssl
//
// 参数：
//
//	ctx - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定要更新的资源部分（必填，如brandingSettings,localizations等）
//	  • onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//
// 返回值：
//
//	*schema.YoutubeChannelsUpdateRes - 包含更新后的频道数据
//	  • 返回完整的channel资源对象
//	error - 调用过程中遇到的错误（如有）
func (c *YoutubeChannelsClient) Update(ctx context.Context, data *schema.YoutubeChannelsUpdateReq) (*schema.YoutubeChannelsUpdateRes, error) {
	result := &schema.YoutubeChannelsUpdateRes{}
	_, err := c.BaseClient.HttpPut(ctx, "/youtube/v3/channels", &object.StringMap{
		"part":                   data.Part,
		"onBehalfOfContentOwner": data.OnBehalfOfContentOwner,
	}, data, nil, result)
	return result, err
}
