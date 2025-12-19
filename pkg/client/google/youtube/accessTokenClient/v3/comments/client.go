// Package comments 提供与YouTube评论相关的API客户端实现
// 包含评论的列表、插入、更新、删除等操作
package comments

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v3/comments/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeCommentsClient 是YouTube评论API的客户端实现
// 继承自kernel.BaseClient，提供HTTP请求的基础功能
type YoutubeCommentsClient struct {
	*kernel.BaseClient
}

// NewClient 创建新的YoutubeCommentsClient实例
// 参数:
//
//	c - 基础客户端实例
//
// 返回:
//
//	*YoutubeCommentsClient - 新的评论客户端实例
func NewClient(c *kernel.BaseClient) *YoutubeCommentsClient {
	return &YoutubeCommentsClient{
		BaseClient: c,
	}
}

// ## List 获取YouTube评论列表
// 该接口用于获取YouTube视频或频道的评论列表
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如snippet,id等）
//	  • id: 评论ID列表（可选，最多100个）
//	  • parentId: 父评论ID（可选，用于获取回复评论）
//	  • maxResults: 返回的最大结果数（可选，默认20，最大100）
//	  • pageToken: 分页令牌（可选）
//	  • textFormat: 评论格式（可选，html或plainText）
//
// 返回值：
//
//	*schema.YouTubeCommentsListRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Items: 评论列表
//	  • PageInfo: 分页信息
//	  • NextPageToken: 下一页令牌
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCommentsClient) List(ctx context.Context, data *schema.YoutubeCommentsListReq) (*schema.YoutubeCommentsListRes, error) {
	result := &schema.YoutubeCommentsListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/comments", params, nil, nil, result)
	return result, err
}

// ## Insert 创建新的YouTube评论
// 该接口用于在指定视频或评论下创建新的评论
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如snippet,id等）
//	  • snippet: 评论内容，包含以下子字段：
//	    - videoId: 视频ID（可选，用于视频评论）
//	    - channelId: 频道ID（可选，用于频道评论）
//	    - textOriginal: 评论内容（必填）
//	    - parentId: 父评论ID（可选，用于回复评论）
//
// 返回值：
//
//	*schema.YouTubeCommentsInsertRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 新评论的唯一标识符
//	  • Snippet: 评论的基本信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCommentsClient) Insert(ctx context.Context, data *schema.YoutubeCommentsInsertReq) (*schema.YoutubeCommentsInsertRes, error) {
	result := &schema.YoutubeCommentsInsertRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/youtube/v3/comments", nil, data, nil, result)
	return result, err
}

// ## Update 更新现有的YouTube评论
// 该接口用于更新指定评论的内容
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如snippet,id等）
//	  • id: 要更新的评论ID（必填）
//	  • snippet: 更新后的评论内容，包含以下子字段：
//	    - textOriginal: 新的评论内容（必填）
//
// 返回值：
//
//	*schema.YouTubeCommentsUpdateRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 更新后的评论ID
//	  • Snippet: 更新后的评论信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCommentsClient) Update(ctx context.Context, data *schema.YoutubeCommentsUpdateReq) (*schema.YoutubeCommentsUpdateRes, error) {
	result := &schema.YoutubeCommentsUpdateRes{}
	_, err := c.BaseClient.HttpPut(ctx, "/youtube/v3/comments", nil, data, nil, result)
	return result, err
}

// ## SetModerationStatus 设置评论的审核状态
// 该接口用于设置指定评论的审核状态
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 要设置审核状态的评论ID（必填）
//	  • moderationStatus: 新的审核状态（必填，可选值：heldForReview, likelySpam, published）
//	  • banAuthor: 是否禁止评论作者（可选，默认false）
//
// 返回值：
//
//	*schema.YouTubeCommentsSetModerationStatusRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 评论ID
//	  • Snippet: 评论的基本信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCommentsClient) SetModerationStatus(ctx context.Context, data *schema.YoutubeCommentsSetModerationStatusReq) (*schema.YoutubeCommentsSetModerationStatusRes, error) {
	result := &schema.YoutubeCommentsSetModerationStatusRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/youtube/v3/comments/setModerationStatus", nil, data, nil, result)
	return result, err
}

// ## Delete 删除指定的YouTube评论
// 该接口用于删除指定的评论
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 要删除的评论ID（必填）
//
// 返回值：
//
//	*schema.YouTubeCommentsDeleteRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 已删除的评论ID
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCommentsClient) Delete(ctx context.Context, data *schema.YoutubeCommentsDeleteReq) (*schema.YoutubeCommentsDeleteRes, error) {
	result := &schema.YoutubeCommentsDeleteRes{}
	_, err := c.BaseClient.HttpDelete(ctx, "/youtube/v3/comments", nil, data, nil, result)
	return result, err
}
