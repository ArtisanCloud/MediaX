package search

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v3/search/schema"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
)

type YoutubeSearchClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeSearchClient {
	return &YoutubeSearchClient{
		BaseClient: c,
	}
}

// ## List 搜索
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/search/list?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet）
//	  • forContentOwner: 是否限制为内容所有者拥有的视频（可选，仅供 YouTube 内容合作伙伴使用）
//	  • forDeveloper: 是否限制为开发者上传的视频（可选）
//	  • forMine: 是否限制为经过身份验证的用户拥有的视频（可选）
//	  • channelId: 频道ID（可选）
//	  • channelType: 频道类型（可选，如 any,show）
//	  • eventType: 事件类型（可选，如 completed,live,upcoming）
//	  • location: 地理位置坐标（可选，如 (37.42307,-122.08427)）
//	  • locationRadius: 地理位置半径（可选，如 1500m,5km）
//	  • maxResults: 返回的最大结果数（可选，默认5，最大50）
//	  • onBehalfOfContentOwner: 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
//	  • q: 搜索查询字词（可选）
//	  • regionCode: 区域代码（可选，ISO 3166-1 alpha-2 国家/地区代码）
//	  • relevanceLanguage: 相关性语言（可选，如 zh-Hans,zh-Hant）
//	  • safeSearch: 安全搜索（可选，如 moderate,none,strict）
//	  • topicId: 主题ID（可选）
//	  • type: 资源类型（可选，如 video,channel,playlist）
//	  • videoCaption: 视频字幕（可选，如 any,closedCaption,none）
//	  • videoCategoryId: 视频分类ID（可选）
//	  • videoDefinition: 视频清晰度（可选，如 any,high,standard）
//	  • videoDimension: 视频维度（可选，如 2d,3d,any）
//	  • videoDuration: 视频时长（可选，如 any,long,medium,short）
//	  • videoEmbeddable: 视频是否可嵌入（可选，如 any,true）
//	  • videoLicense: 视频许可类型（可选，如 any,creativeCommon,youtube）
//	  • videoPaidProductPlacement: 视频是否包含付费宣传内容（可选，如 any,true）
//	  • videoSyndicated: 视频是否已同步（可选，如 any,true）
//	  • videoType: 视频类型（可选，如 any,episode,movie）
//
// 返回值：
//
//	*schema.YouTubeSearchListRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的 ETag
//	  • NextPageToken: 下一页令牌
//	  • PrevPageToken: 上一页令牌
//	  • RegionCode: 区域代码
//	  • PageInfo: 分页信息
//	  • Items: 搜索结果列表
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeSearchClient) List(ctx context.Context, data *schema.YouTubeSearchListReq) (*schema.YouTubeSearchListRes, error) {
	result := &schema.YouTubeSearchListRes{}
	params, err := utils.StructToQueryParams(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/search", params, nil, nil, result)
	return result, err
}
