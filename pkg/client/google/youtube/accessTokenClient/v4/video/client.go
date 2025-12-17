package video

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v4/video/schema"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubeVideoClient 是 YouTube 视频相关 API 的客户端
type YoutubeVideoClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 YoutubeVideoClient 实例
func NewClient(c *kernel.BaseClient) *YoutubeVideoClient {
	return &YoutubeVideoClient{
		BaseClient: c,
	}
}

// ## List 获取视频列表
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/videos/list?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet,contentDetails 等）
//	  • id: 视频ID列表（可选，最多50个）
//	  • chart: 指定要检索的图表类型（可选，如 mostPopular）
//	  • maxResults: 返回的最大结果数（可选，默认5，最大50）
//	  • pageToken: 分页令牌（可选）
//	  • regionCode: 区域代码（可选）
//	  • videoCategoryId: 视频分类ID（可选）
//
// 返回值：
//
//	*schema.YouTubeVideoListRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的 ETag
//	  • Items: 视频列表
//	  • PageInfo: 分页信息
//	  • NextPageToken: 下一页令牌
//	  • PrevPageToken: 上一页令牌
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeVideoClient) List(ctx context.Context, data *schema.YouTubeVideoListReq) (*schema.YouTubeVideoListRes, error) {
	result := &schema.YouTubeVideoListRes{}

	params, err := utils.StructToQueryParams(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/videos", params, nil, nil, result)
	return result, err
}

// ## Insert 上传视频到 YouTube
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videos/insert
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet,contentDetails 等）
//	  • onBehalfOfContentOwner: 内容所有者（可选）
//	  • onBehalfOfContentOwnerChannel: 内容所有者频道（可选）
//	  • notifySubscribers: 是否通知订阅者（可选）
//	  • autoLevels: 是否自动调整视频质量（可选）
//	  • stabilize: 是否稳定视频（可选）
//
// 返回值：
//
//	*schema.YouTubeVideoInsertRes 包含以下字段：
//	  • Video: 视频信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeVideoClient) Insert(ctx context.Context, data *schema.YouTubeVideoInsertReq) (*schema.YouTubeVideoInsertRes, error) {
	result := &schema.YouTubeVideoInsertRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	// 假设 BaseClient 有 HttpPost 方法用于处理上传请求
	// 这里需要根据实际情况处理媒体上传
	_, err = c.BaseClient.HttpPost(ctx, "/upload/youtube/v3/videos", params, nil, nil, result)
	return result, err
}

// ## Update 更新视频元数据
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videos/update
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet,contentDetails 等）
//	  • onBehalfOfContentOwner: 内容所有者（可选）
//	  • video: 视频元数据（必填）
//
// 返回值：
//
//	*schema.YouTubeVideoUpdateRes 包含以下字段：
//	  • Video: 视频信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeVideoClient) Update(ctx context.Context, data *schema.YouTubeVideoUpdateReq) (*schema.YouTubeVideoUpdateRes, error) {
	result := &schema.YouTubeVideoUpdateRes{}

	_, err := c.BaseClient.HttpPut(ctx, "/youtube/v3/video", nil, data, nil, result)
	return result, err
}

// ## Update 更新视频元数据
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videos/update
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • ID: 视频ID（必填）
//	  • onBehalfOfContentOwner: 内容所有者（可选）
//
// 返回值：
//
//	*schema.YouTubeVideoUpdateRes HTTP 204 返回码
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeVideoClient) Delete(ctx context.Context, data *schema.YouTubeVideoDeleteReq) (*schema.YouTubeVideoDeleteRes, error) {
	result := &schema.YouTubeVideoDeleteRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpDelete(ctx, "/youtube/v3/video", params, nil, nil, result)
	return result, err
}

// ## Rate 为视频评分
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videos/rate?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 视频ID（必填）
//	  • rating: 评分类型（必填，如 like,dislike,none）
//
// 返回值：
//
//	*schema.YouTubeVideoRateRes HTTP 204 返回码
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeVideoClient) Rate(ctx context.Context, data *schema.YouTubeVideoRateReq) (*schema.YouTubeVideoRateRes, error) {
	result := &schema.YouTubeVideoRateRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpPost(ctx, "/videos/rate", params, nil, nil, result)
	return result, err
}

// ## GetRating 获取视频评分
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videos/getRating?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 视频ID列表（必填，最多50个）
//	  • OnBehalfOfContentOwner: 内容所有者（可选）
//
// 返回值：
//
//	*schema.YouTubeVideoGetRatingRes 包含以下字段：
//	  • Items: 视频评分列表，每个视频包含以下字段：
//	    • Kind: 资源类型
//	    • ETag: 资源的 ETag
//	    • Items: 评分信息列表
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeVideoClient) GetRating(ctx context.Context, data *schema.YouTubeVideoGetRatingReq) (*schema.YouTubeVideoGetRatingRes, error) {
	result := &schema.YouTubeVideoGetRatingRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpPost(ctx, "/videos/getRating", params, nil, nil, result)
	return result, err
}

// ## ReportAbuse 举报视频
//
// 接口文档参考：
// https://developers.google.com/youtube/v3/docs/videos/reportAbuse?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 视频ID（必填）
//	  • reasonId: 举报原因ID（必填）
//	  • secondaryReasonId: 次要举报原因ID（可选）
//	  • comments: 举报说明（可选）
//	  • language: 语言代码（可选）
//
// 返回值：
//
//	*schema.YouTubeVideoReportAbuseRes HTTP 204 返回码
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeVideoClient) ReportAbuse(ctx context.Context, data *schema.YouTubeVideoReportAbuseReq) (*schema.YouTubeVideoReportAbuseRes, error) {
	result := &schema.YouTubeVideoReportAbuseRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = c.BaseClient.HttpPost(ctx, "/videos/reportAbuse", params, nil, nil, result)
	return result, err
}
