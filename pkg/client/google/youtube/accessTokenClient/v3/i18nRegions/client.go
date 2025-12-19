// Package i18nRegions 提供 YouTube i18nRegions API 的客户端实现
// 主要用于获取 YouTube 支持的区域列表信息
package i18nRegions

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v3/i18nRegions/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeI18nRegionsClient 是 YouTube i18nRegions API 的客户端
// 继承自 kernel.BaseClient，提供基础的 HTTP 请求能力
// 用于获取 YouTube 支持的区域列表信息
// 参考文档：https://developers.google.cn/youtube/v3/docs/i18nRegions/list?hl=zh-cn
type YoutubeI18nRegionsClient struct {
	*kernel.BaseClient
}

// NewClient 创建并返回一个新的 YoutubeI18nRegionsClient 实例
// 参数：
//
//	c - 基础客户端实例，提供 HTTP 请求能力
//
// 返回值：
//
//	*YoutubeI18nRegionsClient - 新的客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeI18nRegionsClient {
	return &YoutubeI18nRegionsClient{
		BaseClient: c,
	}
}

// ## List 获取 YouTube 支持的区域列表
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/i18nRegions/list?hl=zh-cn
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • hl: 指定本地化语言代码（可选）
//
// 返回值：
//
//	*schema.YoutubeI18nRegionsListRes 包含以下字段：
//	  • Items: 区域列表
//	  • Kind: 资源类型
//	  • ETag: 资源的 ETag
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeI18nRegionsClient) List(ctx context.Context, data *schema.YoutubeI18nRegionsListReq) (*schema.YoutubeI18nRegionsListRes, error) {
	result := &schema.YoutubeI18nRegionsListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/i18nRegions", params, nil, nil, result)
	return result, err
}
