package videoCategory

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/videoCategory/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeVideoCategoryClient 是 YouTube 视频类别相关 API 的客户端
type YoutubeVideoCategoryClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 YoutubeVideoCategoryClient 实例
func NewClient(c *kernel.BaseClient) *YoutubeVideoCategoryClient {
	return &YoutubeVideoCategoryClient{
		BaseClient: c,
	}
}

// ## List 获取视频分类列表
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videoCategories/list?hl=zh-cn
//
// 参数：
//   ctx  - 请求上下文
//   data - 请求参数，包含以下字段：
//     • part: 指定返回的资源部分（必填，如 snippet）
//     • id: 视频分类ID（可选）
//     • regionCode: 区域代码（可选）
//     • hl: 语言代码（可选）
//
// 返回值：
//   *schema.YouTubeVideoCategoriesRes 包含以下字段：
//     • Kind: 资源类型
//     • ETag: 资源的 ETag
//     • NextPageToken: 下一页令牌
//     • PrevPageToken: 上一页令牌
//     • PageInfo: 分页信息
//       • TotalResults: 总结果数
//       • ResultsPerPage: 每页结果数
//     • Items: 视频分类列表
//   error 调用过程中遇到的错误（如有）
func (c *YoutubeVideoCategoryClient) List(ctx context.Context, data *schema.YouTubeVideoCategoriesReq) (*schema.YouTubeVideoCategoriesRes, error) {
	result := &schema.YouTubeVideoCategoriesRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/videoCategories", params, nil, result)
	return result, err
}
