// Package members 提供了 YouTube 频道会员相关功能的客户端实现
// 该包用于管理和操作 YouTube 频道的会员系统，包括获取会员列表等功能
package members

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v4/members/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeMembersClient 是 YouTube 频道会员管理的客户端结构体
// 继承自 kernel.BaseClient，用于处理 YouTube 频道会员相关的 API 请求
type YoutubeMembersClient struct {
	// BaseClient 提供基础的 HTTP 请求、认证、日志等功能
	*kernel.BaseClient
}

// NewClient 创建一个新的 YouTube 频道会员客户端实例
// 参数 c *kernel.BaseClient: 基础客户端实例，提供 HTTP 请求、认证等基础功能
// 返回值 *YoutubeMembersClient: 初始化后的频道会员客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeMembersClient {
	return &YoutubeMembersClient{
		BaseClient: c,
	}
}

// ## List 获取频道会员列表
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/members/list?hl=zh-cn
//
// 注意事项：
//   - 仅限个人创作者使用此接口
//   - 需要频道已启用会员功能
//   - 需要 OAuth2 授权，scope 为 youtube.channel-memberships.view
//   - 调用此方法的配额费用为 1 个单位
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet）
//	  • mode: 会员筛选模式（可选，all 或 updates）
//	  • filterByMemberChannelId: 按会员频道ID筛选（可选）
//	  • hasAccessToLevel: 按会员等级筛选（可选）
//	  • maxResults: 返回的最大结果数（可选，默认5，最大1000）
//	  • pageToken: 分页令牌（可选）
//
// 返回值：
//
//	*schema.YoutubeMembersListRes 包含以下字段：
//	  • Kind: 资源类型（youtube#memberListResponse）
//	  • ETag: 资源的 ETag
//	  • NextPageToken: 下一页令牌
//	  • PrevPageToken: 上一页令牌
//	  • PageInfo: 分页信息
//	  • Items: 会员列表
//	error 可能的错误：
//	  • channelMembershipsNotEnabled: 频道未启用会员功能
//	  • insufficientPermissions: 权限不足
//	  • 其他 API 调用相关错误
func (c *YoutubeMembersClient) List(ctx context.Context, data *schema.YoutubeMembersListReq) (*schema.YoutubeMembersListRes, error) {
	result := &schema.YoutubeMembersListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/members", params, nil, nil, result)
	return result, err
}
