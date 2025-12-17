// Package commentThreads 提供与YouTube评论会话相关的API客户端
// 该包实现了YouTube Data API v3中评论会话相关接口的封装
package commentThreads

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v4/commentThreads/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeCommentThreadsClient 封装了YouTube评论会话相关的API操作
// 继承自kernel.BaseClient，提供基础的HTTP请求能力
// 包含评论会话的列表获取、创建等操作
type YoutubeCommentThreadsClient struct {
	*kernel.BaseClient
}

// NewClient 创建并返回一个新的YoutubeCommentThreadsClient实例
// 参数：
//
//	c - 基础客户端实例，提供HTTP请求等基础能力
//
// 返回值：
//
//	*YoutubeCommentThreadsClient 新的评论会话客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeCommentThreadsClient {
	return &YoutubeCommentThreadsClient{
		BaseClient: c,
	}
}

// ## List 获取YouTube视频的评论会话列表
// 接口文档参考：https://developers.google.cn/youtube/v3/docs/commentThreads/list?hl=zh-cn
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如snippet,replies等）
//	  • videoId: 视频ID（可选）
//	  • maxResults: 返回的最大结果数（可选，默认20，最大100）
//	  • pageToken: 分页令牌（可选）
//	  • textFormat: 评论文本格式（可选，如plainText,html）
//
// 返回值：
//
//	*schema.YoutubeCommentThreadsListRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Items: 评论会话列表
//	  • PageInfo: 分页信息
//	  • Items : 评论会话列表
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCommentThreadsClient) List(ctx context.Context, data *schema.YoutubeCommentThreadsListReq) (*schema.YoutubeCommentThreadsListRes, error) {
	result := &schema.YoutubeCommentThreadsListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/commentThreads", params, nil, nil, result)
	return result, err
}

// ## Insert 创建新的评论会话
// 接口文档参考：https://developers.google.cn/youtube/v3/docs/commentThreads/insert?hl=zh-cn
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如snippet）
//	  • snippet.videoId: 视频ID（可选）
//	  • snippet.topLevelComment: 顶级评论内容（必填）
//	  • snippet.canReply: 是否允许回复（可选）
//	  • snippet.isPublic: 是否公开（可选）
//
// 返回值：
//
//	*schema.YoutubeCommentThreadsInsertRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 新创建的评论会话ID
//	  • Snippet: 评论会话的基本信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCommentThreadsClient) Insert(ctx context.Context, data *schema.YoutubeCommentThreadsInsertReq) (*schema.YoutubeCommentThreadsInsertRes, error) {
	result := &schema.YoutubeCommentThreadsInsertRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/youtube/v3/commentThreads", &object.StringMap{
		"part": data.Part,
	}, data, nil, result)
	return result, err
}
