// Package channelSections 提供与YouTube频道版块相关的API客户端
// 包含列表、创建、更新和删除频道版块的功能
package channelSections

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/channelSections/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeChannelSectionsClient 是YouTube频道版块API的客户端
// 继承自BaseClient，提供HTTP请求的基础功能
type YoutubeChannelSectionsClient struct {
	*kernel.BaseClient
}

// NewClient 创建新的YoutubeChannelSectionsClient实例
// c: 基础客户端实例
// 返回: 初始化后的YoutubeChannelSectionsClient指针
func NewClient(c *kernel.BaseClient) *YoutubeChannelSectionsClient {
	return &YoutubeChannelSectionsClient{
		BaseClient: c,
	}
}

// ## List 获取频道版块列表
// 接口文档参考：https://developers.google.cn/youtube/v3/docs/channelSections/list?hl=zh-cn
// ctx: 请求上下文
// data: 请求参数，包含以下字段：
//   - part: 指定返回的资源部分（必填，如snippet,contentDetails等）
//   - id: 频道版块ID列表（可选，最多50个）
//   - hl: 本地化语言代码（可选）
//   - maxResults: 返回的最大结果数（可选，默认5，最大50）
//   - pageToken: 分页令牌（可选）
//   - onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//
// 返回值：
//
//	*schema.YoutubeChannelSectionsListRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • NextPageToken: 下一页令牌
//	  • PrevPageToken: 上一页令牌
//	  • PageInfo: 分页信息
//	  • Items: 频道版块列表
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeChannelSectionsClient) List(ctx context.Context, data *schema.YoutubeChannelSectionsListReq) (*schema.YoutubeChannelSectionsListRes, error) {
	result := &schema.YoutubeChannelSectionsListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/channelSections", params, nil, nil, result)
	return result, err
}

// ## Insert 创建新的频道版块
// 接口文档参考：https://developers.google.cn/youtube/v3/docs/channelSections/insert?hl=zh-cn
// ctx: 请求上下文
// data: 请求参数，包含以下字段：
//   - part: 指定返回的资源部分（必填，如snippet,contentDetails等）
//   - onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//   - onBehalfOfContentOwnerChannel: 代表内容所有者频道执行操作（可选）
//
// 返回值：
//
//	*schema.YoutubeChannelSectionsInsertRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 新创建的频道版块ID
//	  • Snippet: 频道版块的基本信息
//	  • ContentDetails: 频道版块的内容详情
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeChannelSectionsClient) Insert(ctx context.Context, data *schema.YoutubeChannelSectionsInsertReq) (*schema.YoutubeChannelSectionsInsertRes, error) {
	result := &schema.YoutubeChannelSectionsInsertRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/youtube/v3/channelSections", &object.StringMap{
		"part":                          data.Part,
		"onBehalfOfContentOwner":        data.OnBehalfOfContentOwner,
		"onBehalfOfContentOwnerChannel": data.OnBehalfOfContentOwnerChannel,
	}, data, nil, result)
	return result, err
}

// ## Update 更新现有频道版块
// 接口文档参考：https://developers.google.cn/youtube/v3/docs/channelSections/update?hl=zh-cn
// ctx: 请求上下文
// data: 请求参数，包含以下字段：
//   - part: 指定返回的资源部分（必填，如snippet,contentDetails等）
//   - onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//
// 返回值：
//
//	*schema.YoutubeChannelSectionsUpdateRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 更新的频道版块ID
//	  • Snippet: 更新后的频道版块基本信息
//	  • ContentDetails: 更新后的频道版块内容详情
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeChannelSectionsClient) Update(ctx context.Context, data *schema.YoutubeChannelSectionsUpdateReq) (*schema.YoutubeChannelSectionsUpdateRes, error) {
	result := &schema.YoutubeChannelSectionsUpdateRes{}
	_, err := c.BaseClient.HttpPut(ctx, "/youtube/v3/channelSections", &object.StringMap{
		"part":                   data.Part,
		"onBehalfOfContentOwner": data.OnBehalfOfContentOwner,
	}, data, nil, result)
	return result, err
}

// ## Delete 删除指定的频道版块
// 接口文档参考：https://developers.google.cn/youtube/v3/docs/channelSections/delete?hl=zh-cn
// ctx: 请求上下文
// data: 请求参数，包含以下字段：
//   - id: 要删除的频道版块ID
//   - onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//
// 返回值：
//
//	*schema.YoutubeChannelSectionsDeleteRes 包含以下字段：
//	  • ChannelSection: 已删除的频道版块信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeChannelSectionsClient) Delete(ctx context.Context, data *schema.YoutubeChannelSectionsDeleteReq) (*schema.YoutubeChannelSectionsDeleteRes, error) {
	result := &schema.YoutubeChannelSectionsDeleteRes{}
	_, err := c.BaseClient.HttpDelete(ctx, "/youtube/v3/channelSections", &object.StringMap{
		"id":                     data.Id,
		"onBehalfOfContentOwner": data.OnBehalfOfContentOwner,
	}, nil, nil, result)
	return result, err
}
