package data

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/data/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// BiliBiliDataClient 数据中心客户端
type BiliBiliDataClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 BiliBiliDataClient 实例。
func NewClient(c *kernel.BaseClient) *BiliBiliDataClient {
	return &BiliBiliDataClient{
		BaseClient: c,
	}
}

// ## GetUserStat 获取用户数据
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/90936ab5-7c06-e24f-2ad0-0fd6e10c7386#h1-u83B7u53D6u7528u6237u6570u636E
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*schema.BiliBiliDataGetUserStatRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 消息
//	  • Data: 用户数据，包含关注数、粉丝数、投稿数等
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliDataClient) GetUserStat(ctx context.Context) (*schema.BiliBiliDataGetUserStatRes, error) {
	result := &schema.BiliBiliDataGetUserStatRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/data/user/stat", nil, nil, nil, result)
	return result, err
}

// ## GetVideoStat 获取单个稿件数据
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/3f46ac2e-1318-3aa0-5548-0d9fd624d520#h1-u83B7u53D6u5355u4E2Au7A3Fu4EF6u6570u636E
//
// 参数：
//
//	ctx - 请求上下文
//	resourceId - 稿件唯一ID
//
// 返回值：
//
//	*schema.BiliBiliDataGetVideoStatRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 消息
//	  • Data: 稿件统计数据，包含：
//	    - Title: 标题
//	    - Ptime: 发布时间
//	    - View: 播放数
//	    - Danmaku: 弹幕数
//	    - Reply: 评论数
//	    - Favorite: 收藏数
//	    - Coin: 投币数
//	    - Share: 分享数
//	    - Like: 点赞数
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliDataClient) GetVideoStat(ctx context.Context, resourceId string) (*schema.BiliBiliDataGetVideoStatRes, error) {
	result := &schema.BiliBiliDataGetVideoStatRes{}
	params := &object.StringMap{
		"resource_id": resourceId,
	}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/data/arc/stat", params, nil, nil, result)
	return result, err
}

// ## GetVideoStatIncrement 获取整体稿件增量数据
//
// 接口文档参考：
// https://member.bilibili.com/arcopen/fn/data/arc/inc-stats
//
// 返回值：
//
//	*schema.BiliBiliDataGetVideoStatIncrementRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 消息
//	  • Data: 稿件增量统计数据，包含：
//	    - IncClick: 播放数增量
//	    - IncDm: 弹幕数增量
//	    - IncReply: 评论数增量
//	    - IncFav: 收藏数增量
//	    - IncCoin: 投币数增量
//	    - IncShare: 分享数增量
//	    - IncLike: 点赞数增量
//	    - IncElec: 充电数增量
//	error 调用过程中遇到的错误（如有）
//
// https://open.bilibili.com/doc/4/708a4ce5-5a6e-d8ce-b2e8-ff8a26a83b05#h1-u83B7u53D6u6574u4F53u7A3Fu4EF6u589Eu91CFu6570u636E
func (c *BiliBiliDataClient) GetVideoStatIncrement(ctx context.Context) (*schema.BiliBiliDataGetVideoStatIncrementRes, error) {
	result := &schema.BiliBiliDataGetVideoStatIncrementRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/data/arc/inc-stats", nil, nil, nil, result)
	return result, err
}

// ## GetColumnStat 获取单一专栏数据
//
// 接口文档参考：
// `https://open.bilibili.com/doc/4/46903ddc-1264-ddac-a3c7-559d12f32118#h1-u83B7u53D6u5355u4E00u4E13u680Fu6570u636E`
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • ids: 文章ids（必填，多个id用英文逗号分割）
//
// 返回值：
//
//	*schema.BiliBiliDataGetColumnStatRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 消息
//	  • Data: 专栏数据，包含：
//	    - id: 文章id
//	    - category: 分类信息
//	    - title: 文章标题
//	    - stats: 统计数据（浏览数、收藏数等）
//	    - 其他字段详见官方文档
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliDataClient) GetColumnStat(ctx context.Context, data *schema.BiliBiliDataGetColumnStatReq) (*schema.BiliBiliDataGetColumnStatRes, error) {
	result := &schema.BiliBiliDataGetColumnStatRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/data/art/stat", nil, data, nil, result)
	return result, err
}

// 获取整体专栏增量数据
// ## GetColumnStatIncrement 获取整体专栏增量数据
//
// 接口文档参考：
// `https://open.bilibili.com/doc/4/f5d81e34-ab1b-ac92-2c5e-34c3c0d33828#h1-u83B7u53D6u6574u4F53u4E13u680Fu589Eu91CFu6570u636E`
//
// 返回值：
//
//	*schema.BiliBiliDataGetColumnStatIncrementRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 消息
//	  • Data: 专栏增量统计数据，包含：
//	    - IncReply: 评论数增量
//	    - IncRead: 阅读数增量
//	    - IncFav: 收藏数增量
//	    - IncLikes: 点赞数增量
//	    - IncShare: 分享数增量
//	    - IncCoin: 投币数增量
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliDataClient) GetColumnStatIncrement(ctx context.Context) (*schema.BiliBiliDataGetColumnStatIncrementRes, error) {
	result := &schema.BiliBiliDataGetColumnStatIncrementRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/data/art/inc-stats", nil, nil, nil, result)
	return result, err
}
