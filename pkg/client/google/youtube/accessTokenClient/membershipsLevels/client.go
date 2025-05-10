// Package membershipsLevels 提供了 YouTube 会员等级相关功能的客户端实现
// 该包用于管理和操作 YouTube 频道的会员等级系统，包括获取会员等级列表、详情等功能
package membershipsLevels

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/membershipsLevels/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeMembershipsLevelsClient 是 YouTube 会员等级管理的客户端结构体
// 继承自 kernel.BaseClient，用于处理 YouTube 会员等级相关的 API 请求
type YoutubeMembershipsLevelsClient struct {
	// BaseClient 提供基础的 HTTP 请求、认证、日志等功能
	*kernel.BaseClient
}

// NewClient 创建一个新的 YouTube 会员等级客户端实例
// 参数 c *kernel.BaseClient: 基础客户端实例，提供 HTTP 请求、认证等基础功能
// 返回值 *YoutubeMembershipsLevelsClient: 初始化后的会员等级客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeMembershipsLevelsClient {
	return &YoutubeMembershipsLevelsClient{
		BaseClient: c,
	}
}

// ## List 获取频道会员等级列表
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/membershipsLevels/list?hl=zh-cn
//
// 注意事项：
//   • 仅限个人创作者使用此接口
//   • 需要频道已启用会员功能
//   • 调用此方法的配额费用为 1 个单位
//
// 参数：
//   ctx  - 请求上下文
//   data - 请求参数，包含以下字段：
//     • part: 指定返回的资源部分（必填，如 id,snippet）
//
// 返回值：
//   *schema.YoutubeMembershipsLevelsListRes 包含以下字段：
//     • Kind: 资源类型（youtube#membershipsLevelListResponse）
//     • ETag: 资源的 ETag
//     • PageInfo: 分页信息，包含以下字段：
//       • TotalResults: 总结果数
//       • ResultsPerPage: 每页结果数
//     • Items: 会员等级列表，按显示顺序排列
//   error 可能的错误：
//     • channelMembershipsNotEnabled: 频道未启用会员功能
//     • 其他 API 调用相关错误
func (c *YoutubeMembershipsLevelsClient) List(ctx context.Context, data *schema.YoutubeMembershipsLevelsListReq) (*schema.YoutubeMembershipsLevelsListRes, error) {
	result := &schema.YoutubeMembershipsLevelsListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/membershipsLevels", params, nil, result)
	return result, err
}