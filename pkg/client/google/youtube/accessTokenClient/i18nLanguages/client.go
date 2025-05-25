package i18nLanguages

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/i18nLanguages/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// Package i18nLanguages 提供与YouTube i18nLanguages API交互的客户端实现
// 该包主要用于获取YouTube支持的语言列表

type YoutubeI18nLanguagesClient struct {
	// BaseClient 嵌入基础客户端，继承其HTTP请求和认证功能
	*kernel.BaseClient
}

// NewClient 创建并返回一个新的YoutubeI18nLanguagesClient实例
// 参数：
//
//	c - 基础客户端实例
//
// 返回值：
//
//	*YoutubeI18nLanguagesClient - 新的i18nLanguages客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeI18nLanguagesClient {
	return &YoutubeI18nLanguagesClient{
		BaseClient: c,
	}
}

// ## List 获取YouTube支持的语言列表
// 接口文档参考：https://developers.google.cn/youtube/v3/docs/i18nLanguages/list?hl=zh-cn
// 参数：
//
//	ctx - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • hl: 指定本地化语言代码（可选）
//
// 返回值：
//
//		*schema.YoutubeI18nLanguagesListRes - 包含语言列表的响应结果
//	  • Kind - 资源类型
//	  • ETag - 资源的ETag
//	  • Items - 语言列表
//		error - 调用过程中遇到的错误（如有）
func (c *YoutubeI18nLanguagesClient) List(ctx context.Context, data *schema.YoutubeI18nLanguagesListReq) (*schema.YoutubeI18nLanguagesListRes, error) {
	result := &schema.YoutubeI18nLanguagesListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/i18nLanguages", params, nil, nil, result)
	return result, err
}
