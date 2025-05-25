package data

import (
	"context"
	"fmt"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/connection/data/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// DouYinConnectionDataClient 抖音用户数据客户端
type DouYinConnectionDataClient struct {
	*kernel.BaseClient
}

// NewClient 创建新的抖音用户数据客户端实例
func NewClient(c *kernel.BaseClient) *DouYinConnectionDataClient {
	return &DouYinConnectionDataClient{
		BaseClient: c,
	}
}

// ## GetUserVideoStatus 获取用户视频情况
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/user-data/get-user-video-status
//
// 参数：
//
//	ctx  - 请求上下文
//	dateType - 日期类型，1: 昨天，2: 最近7天，3: 最近30天
//
// 返回值：
//
//	*schema.DouYinConnectionDataUserVideoStatusRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - ResultList: 用户视频情况列表，包含多个视频的详细信息
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionDataClient) GetUserVideoStatus(ctx context.Context, dateType int64) (*schema.DouYinConnectionDataUserVideoStatusRes, error) {
	result := &schema.DouYinConnectionDataUserVideoStatusRes{}

	params := &object.StringMap{
		"date_type": fmt.Sprintf("%d", dateType),
	}

	_, err := c.BaseClient.HttpGet(ctx, "/data/external/user/item/", params, nil, nil, result)
	return result, err
}

// ## GetUserFansCount 获取用户粉丝数
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/user-data/get-user-fans-count
//
// 参数：
//
//	ctx  - 请求上下文
//	dateType - 日期类型，1: 昨天，2: 最近7天，3: 最近30天
//
// 返回值：
//
//	*schema.DouYinConnectionDataUserFansCountRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - ResultList: 用户粉丝情况列表，包含多个日期的粉丝数据
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionDataClient) GetUserFansCount(ctx context.Context, dateType int64) (*schema.DouYinConnectionDataUserFansCountRes, error) {
	result := &schema.DouYinConnectionDataUserFansCountRes{}
	params := &object.StringMap{
		"date_type": fmt.Sprintf("%d", dateType),
	}
	_, err := c.BaseClient.HttpGet(ctx, "/data/external/user/fans/", params, nil, nil, result)
	return result, err
}

// ## GetUserLikeNumber 获取用户点赞数
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/user-data/get-user-like-number
//
// 参数：
//
//	ctx  - 请求上下文
//	dateType - 日期类型，1: 昨天，2: 最近7天，3: 最近30天
//
// 返回值：
//
//	*schema.DouYinConnectionDataUserLikeNumberRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - ResultList: 用户点赞情况列表，包含多个日期的点赞数据
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionDataClient) GetUserLikeNumber(ctx context.Context, dateType int64) (*schema.DouYinConnectionDataUserLikeNumberRes, error) {
	result := &schema.DouYinConnectionDataUserLikeNumberRes{}
	params := &object.StringMap{
		"date_type": fmt.Sprintf("%d", dateType),
	}
	_, err := c.BaseClient.HttpGet(ctx, "/data/external/user/like/", params, nil, nil, result)
	return result, err
}

// ## GetUserCommentCount 获取用户评论数
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/user-data/get-user-comment-count
//
// 参数：
//
//	ctx  - 请求上下文
//	dateType - 日期类型，1: 昨天，2: 最近7天，3: 最近30天
//
// 返回值：
//
//	*schema.DouYinConnectionDataUserCommentCountRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - ResultList: 用户评论情况列表，包含多个日期的评论数据
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionDataClient) GetUserCommentCount(ctx context.Context, dateType int64) (*schema.DouYinConnectionDataUserCommentCountRes, error) {
	result := &schema.DouYinConnectionDataUserCommentCountRes{}
	params := &object.StringMap{
		"date_type": fmt.Sprintf("%d", dateType),
	}
	_, err := c.BaseClient.HttpGet(ctx, "/data/external/user/comment/", params, nil, nil, result)
	return result, err
}

// ## GetUserProfile 获取用户主页访问数
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/user-data/get-user-home-pv
//
// 参数：
//
//	ctx  - 请求上下文
//	dateType - 日期类型，1: 昨天，2: 最近7天，3: 最近30天
//
// 返回值：
//
//	*schema.DouYinConnectionDataUserProfile 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - ResultList: 用户主页访问情况列表，包含多个日期的主页访问数据
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionDataClient) GetUserProfile(ctx context.Context, dateType int64) (*schema.DouYinConnectionDataUserProfile, error) {
	result := &schema.DouYinConnectionDataUserProfile{}
	params := &object.StringMap{
		"date_type": fmt.Sprintf("%d", dateType),
	}
	_, err := c.BaseClient.HttpGet(ctx, "/data/external/user/profile/", params, nil, nil, result)
	return result, err
}
