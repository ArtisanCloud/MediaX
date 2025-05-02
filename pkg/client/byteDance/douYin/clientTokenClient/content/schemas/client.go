package schemas

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/schemas/schema"
)

// DouYinContentSchemasClient 抖音内容Schema客户端
// 提供与抖音内容Schema相关的接口封装
type DouYinContentSchemasClient struct {
	*kernel.BaseClient
}

// NewClient 初始化并返回一个新的 DouYinContentSchemasClient 实例
func NewClient(c *kernel.BaseClient) *DouYinContentSchemasClient {
	return &DouYinContentSchemasClient{
		BaseClient: c,
	}
}

// ## GetH5Share 获取H5分享跳转链接
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/jump/h5-share
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • ClientTicket: 客户端票据
//	       • ExpireAt: 过期时间戳
//	       • HashtagList: 话题标签列表
//	       • MicroAppInfo: 微应用信息
//	       • PoiId: POI ID
//	       • ShareToPublish: 是否分享后发布
//	       • State: 状态信息
//	       • Title: 标题
//	       • VideoPath: 视频路径
//
// 返回值：
//
//	*schema.DouYinContentSchemasGetH5ShareRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - Schema: H5分享链接
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentSchemasClient) GetH5Share(ctx context.Context, data *schema.DouYinContentSchemasGetH5ShareReq) (*schema.DouYinContentSchemasGetH5ShareRes, error) {
	result := &schema.DouYinContentSchemasGetH5ShareRes{}

	_, err := c.BaseClient.HttpPost(ctx, "/api/douyin/v1/schema/get_share/", nil, data, nil, result)

	return result, err
}

// ## GetUserProfile 获取个人页跳转链接
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/jump/get-user-profile
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • ExpireAt: 过期时间戳
//	       • OpenId: 用户OpenId
//
// 返回值：
//
//	*schema.DouYinContentSchemasGetUserProfileRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - Schema: 个人页跳转链接
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentSchemasClient) GetUserProfile(ctx context.Context, data *schema.DouYinContentSchemasGetUserProfileReq) (*schema.DouYinContentSchemasGetUserProfileRes, error) {
	result := &schema.DouYinContentSchemasGetUserProfileRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/douyin/v1/schema/get_user_profile/", nil, data, nil, result)
	return result, err
}

// ## GetChat 获取个人会话页跳转链接
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/jump/chat
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • ExpireAt: 过期时间戳
//	       • OpenId: 用户OpenId
//
// 返回值：
//
//	*schema.DouYinContentSchemasGetChatRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - Schema: 个人会话页跳转链接
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentSchemasClient) GetChat(ctx context.Context, data *schema.DouYinContentSchemasGetChatReq) (*schema.DouYinContentSchemasGetChatRes, error) {
	result := &schema.DouYinContentSchemasGetChatRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/douyin/v1/schema/get_chat/", nil, data, nil, result)
	return result, err
}

// ## GetItemInfo 获取视频详情页跳转链接
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/jump/item-info
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • ExpireAt: 过期时间戳
//	       • ItemId: 视频id
//	       • VideoId: 视频id
//
// 返回值：
//
//	*schema.DouYinContentSchemasGetItemInfoRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - Schema: 视频详情页跳转链接
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentSchemasClient) GetItemInfo(ctx context.Context, data *schema.DouYinContentSchemasGetItemInfoReq) (*schema.DouYinContentSchemasGetItemInfoRes, error) {
	result := &schema.DouYinContentSchemasGetItemInfoRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/douyin/v1/schema/get_item_info/", nil, data, nil, result)
	return result, err
}

// ## GetLive 获取直播间跳转链接
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/jump/live
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • ExpireAt: 过期时间戳
//	       • RoomId: 直播间ID
//
// 返回值：
//
//	*schema.DouYinContentSchemasGetLiveRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - Schema: 直播间跳转链接
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentSchemasClient) GetLive(ctx context.Context, data *schema.DouYinContentSchemasGetLiveReq) (*schema.DouYinContentSchemasGetLiveRes, error) {
	result := &schema.DouYinContentSchemasGetLiveRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/douyin/v1/schema/get_live/", nil, data, nil, result)
	return result, err
}
