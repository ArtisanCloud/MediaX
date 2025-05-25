package playlistImages

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/playlistImages/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// YoutubePlaylistImagesClient 是一个用于操作 YouTube 播放列表图片的客户端
type YoutubePlaylistImagesClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 YoutubePlaylistImagesClient 实例
func NewClient(c *kernel.BaseClient) *YoutubePlaylistImagesClient {
	return &YoutubePlaylistImagesClient{
		BaseClient: c,
	}
}

// ## List 获取播放列表图片
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/playlistImages/list?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet 等）
//	  • id: 播放列表图片ID列表（可选，以逗号分隔）
//	  • playlistId: 播放列表ID（可选，与 id 互斥）
//	  • maxResults: 返回的最大结果数（可选，默认5，最大50）
//	  • onBehalfOfContentOwner: 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
//	  • onBehalfOfContentOwnerChannel: 内容所有者频道（可选，仅供 YouTube 内容合作伙伴使用）
//	  • pageToken: 分页令牌（可选）
//
// 返回值：
//
//	*schema.YoutubePlaylistImagesListRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的 ETag
//	  • NextPageToken: 下一页令牌
//	  • PrevPageToken: 上一页令牌
//	  • PageInfo: 分页信息
//	  • Items: 播放列表图片列表
//	error 调用过程中遇到的错误（如有）
func (c *YoutubePlaylistImagesClient) List(ctx context.Context, data *schema.YoutubePlaylistImagesListReq) (res *schema.YoutubePlaylistImagesListRes, err error) {
	result := &schema.YoutubePlaylistImagesListRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/playlistImages", params, nil, nil, result)
	return result, err
}

// ## Insert 插入播放列表图片
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/playlistImages/insert?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet 等）
//	  • onBehalfOfContentOwner: 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
//	  • onBehalfOfContentOwnerChannel: 内容所有者频道（可选，仅供 YouTube 内容合作伙伴使用）
//	  • snippet.playlistId: 播放列表ID（必填）
//	  • snippet.resourceId: 资源ID（必填）
//	  • snippet.position: 播放列表项位置（可选）
//	  • contentDetails.note: 备注（可选）
//	  • contentDetails.startAt: 开始时间（可选）
//	  • contentDetails.endAt: 结束时间（可选）
//
// 返回值：
//
//	*schema.YoutubePlaylistImagesInsertRes 包含以下字段：
//	  • PlaylistImages: 播放列表图片信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubePlaylistImagesClient) Insert(ctx context.Context, data *schema.YoutubePlaylistImagesInsertReq) (res *schema.YoutubePlaylistImagesInsertRes, err error) {
	result := &schema.YoutubePlaylistImagesInsertRes{}
	_, err = c.BaseClient.HttpPost(ctx, "/youtube/v3/playlistImages", &object.StringMap{
		"part":                          data.Part,
		"onBehalfOfContentOwner":        data.OnBehalfOfContentOwner,
		"onBehalfOfContentOwnerChannel": data.OnBehalfOfContentOwnerChannel,
	}, data, nil, result)
	return result, err
}

// ## Update 更新播放列表图片
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/playlistImages/update?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet 等）
//	  • onBehalfOfContentOwner: 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
//	  • snippet.playlistId: 播放列表ID（必填）
//	  • snippet.type: 图片类型（必填）
//	  • snippet.width: 图片宽度（可选）
//	  • snippet.height: 图片高度（可选）
//
// 返回值：
//
//	*schema.YoutubePlaylistImagesUpdateRes 包含以下字段：
//	  • PlaylistImages: 播放列表图片信息:
//	error 调用过程中遇到的错误（如有）
func (c *YoutubePlaylistImagesClient) Update(ctx context.Context, data *schema.YoutubePlaylistImagesUpdateReq) (res *schema.YoutubePlaylistImagesUpdateRes, err error) {
	result := &schema.YoutubePlaylistImagesUpdateRes{}
	_, err = c.BaseClient.HttpPut(ctx, "/youtube/v3/playlistImages", &object.StringMap{
		"part":                   data.Part,
		"onBehalfOfContentOwner": data.OnBehalfOfContentOwner,
	}, data, nil, result)
	return result, err
}

// ## Delete 删除播放列表图片
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/playlistImages/delete?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 播放列表图片ID（必填）
//	  • onBehalfOfContentOwner: 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
//
// 返回值：
//
//	*schema.YoutubePlaylistImagesDeleteRes HTTP 204 返回码
//	error 调用过程中遇到的错误（如有）
func (c *YoutubePlaylistImagesClient) Delete(ctx context.Context, data *schema.YoutubePlaylistImagesDeleteReq) (res *schema.YoutubePlaylistImagesDeleteRes, err error) {
	result := &schema.YoutubePlaylistImagesDeleteRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpDelete(ctx, "/youtube/v3/playlistImages", params, nil, nil, result)
	return result, err
}
