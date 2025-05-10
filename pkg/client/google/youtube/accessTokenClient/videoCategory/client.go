// Package videoCategory 提供了 YouTube 视频分类 API 的客户端实现。
// 该包允许开发者获取 YouTube 视频分类信息，包括分类 ID、标题、可分配性等。
// 视频分类是 YouTube 内容组织的重要组成部分，用于对视频进行分类和筛选。
package videoCategory

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/videoCategory/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// ## YoutubeVideoCategoryClient 是 YouTube 视频类别相关 API 的客户端。
// 该客户端提供了访问和管理 YouTube 视频分类的功能，包括：
// - 获取视频分类列表
// - 按区域和语言筛选分类
// - 获取特定分类的详细信息
type YoutubeVideoCategoryClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 YoutubeVideoCategoryClient 实例。
// 参数：
//
//	c - 基础客户端实例，提供 HTTP 请求、认证等基础功能
//
// 返回值：
//
//	*YoutubeVideoCategoryClient - 新创建的视频分类客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeVideoCategoryClient {
	return &YoutubeVideoCategoryClient{
		BaseClient: c,
	}
}

// ## List 获取视频分类列表。
// 该方法允许获取 YouTube 视频分类的列表，支持按区域和语言进行筛选。
// 视频分类是 YouTube 平台的基础功能之一，用于对视频进行分类和组织。
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videoCategories/list?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文，用于控制请求的生命周期
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet）
//	  • id: 视频分类ID（可选，用于获取特定分类的信息）
//	  • regionCode: 区域代码（可选，用于获取特定地区的分类列表）
//	  • hl: 语言代码（可选，用于本地化分类名称和描述）
//
// 返回值：
//
//	*schema.YouTubeVideoCategoriesRes 包含以下字段：
//	  • Kind: 资源类型，标识响应的数据类型
//	  • ETag: 资源的 ETag，用于缓存控制
//	  • NextPageToken: 下一页令牌，用于分页请求
//	  • PrevPageToken: 上一页令牌，用于分页请求
//	  • PageInfo: 分页信息
//	    • TotalResults: 总结果数
//	    • ResultsPerPage: 每页结果数
//	  • Items: 视频分类列表，包含分类的详细信息
//	error 调用过程中遇到的错误（如有）
//
// 使用示例：
//
//	client := videoCategory.NewClient(baseClient)
//	res, err := client.List(ctx, &schema.YouTubeVideoCategoriesReq{
//	    Part: "snippet",
//	    RegionCode: "US",
//	    HL: "en",
//	})
func (c *YoutubeVideoCategoryClient) List(ctx context.Context, data *schema.YouTubeVideoCategoriesReq) (*schema.YouTubeVideoCategoriesRes, error) {
	result := &schema.YouTubeVideoCategoriesRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/videoCategories", params, nil, result)
	return result, err
}
