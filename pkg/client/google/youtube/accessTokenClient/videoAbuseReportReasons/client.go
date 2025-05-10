// Package videoAbuseReportReasons 提供了 YouTube 视频滥用报告原因相关的 API 功能。
//
// 该包实现了 YouTube Data API v3 中视频滥用报告原因相关的接口，主要用于获取可用于报告视频滥用行为的原因列表。
// 通过这些 API，开发者可以：
// 1. 获取 YouTube 平台预定义的视频滥用报告原因
// 2. 支持多语言环境下的原因描述
// 3. 用于构建视频举报功能的用户界面
//
// 使用示例：
//
//	client := videoAbuseReportReasons.NewClient(baseClient)
//	res, err := client.List(ctx, &schema.YouTubeVideoAbuseReportReasonsListReq{
//	  Part: "id,snippet",
//	  HL:   "zh-CN",
//	})
package videoAbuseReportReasons

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/videoAbuseReportReasons/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeVideoAbuseReportReasonsClient 是 YouTube 视频滥用报告原因相关 API 的客户端。
// 该客户端封装了与视频滥用报告原因相关的所有 API 操作，提供了简单易用的接口。
// 它继承自 BaseClient，具备基础的 HTTP 请求能力和认证功能。
type YoutubeVideoAbuseReportReasonsClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 YoutubeVideoAbuseReportReasonsClient 实例。
// 参数:
//   - c: 基础客户端实例，提供 HTTP 请求和认证功能
//
// 返回值:
//   - *YoutubeVideoAbuseReportReasonsClient: 新创建的客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeVideoAbuseReportReasonsClient {
	return &YoutubeVideoAbuseReportReasonsClient{
		BaseClient: c,
	}
}

// ## List 获取视频滥用报告原因列表。
//
// 该方法从 YouTube Data API 获取预定义的视频滥用报告原因列表。这些原因可用于：
// 1. 构建视频举报功能的下拉菜单或选项列表
// 2. 在不同语言环境下显示标准化的举报原因
// 3. 确保举报原因符合 YouTube 的政策和规范
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videoAbuseReportReasons/list?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文，用于控制请求的生命周期
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 id,snippet）
//	         id: 返回举报原因的唯一标识符
//	         snippet: 返回举报原因的详细信息，如标题和描述
//	  • hl: 语言代码（可选，默认 en_US）
//	       用于指定返回结果的语言，如 zh-CN 表示简体中文
//
// 返回值：
//
//	*schema.YouTubeVideoAbuseReportReasonsListRes 包含以下字段：
//	  • Kind: 资源类型，标识响应的数据类型
//	  • ETag: 资源的 ETag，用于缓存控制
//	  • Items: 视频滥用报告原因列表，每个项目包含 ID 和描述信息
//	error 调用过程中遇到的错误（如有）
//
// 使用示例：
//
//	res, err := client.List(ctx, &schema.YouTubeVideoAbuseReportReasonsListReq{
//	  Part: "id,snippet",
//	  HL:   "zh-CN",
//	})
//	if err != nil {
//	  log.Printf("获取视频滥用报告原因失败：%v", err)
//	  return
//	}
//	for _, item := range res.Items {
//	  fmt.Printf("原因ID：%s，描述：%s\n", item.ID, item.Snippet.Label)
//	}
func (c *YoutubeVideoAbuseReportReasonsClient) List(ctx context.Context, data *schema.YouTubeVideoAbuseReportReasonsListReq) (*schema.YouTubeVideoAbuseReportReasonsListRes, error) {
	result := &schema.YouTubeVideoAbuseReportReasonsListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/videoAbuseReportReasons", params, nil, result)
	return result, err
}
