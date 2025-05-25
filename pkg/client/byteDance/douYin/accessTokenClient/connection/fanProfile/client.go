package fanProfile

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/connection/fanProfile/schema"
)

// DouYinConnectionFanProfileClient 抖音粉丝关系数据接口。
// 提供与抖音粉丝画像数据相关的接口封装
type DouYinConnectionFanProfileClient struct {
	*kernel.BaseClient
}

// NewClient 初始化并返回一个新的 DouYinConnectionFanProfileClient 实例
// 参数：
//
//	c - 基础客户端实例
//
// 返回值：
//
//	*DouYinConnectionFanProfileClient: 新的抖音粉丝关系数据客户端实例
func NewClient(c *kernel.BaseClient) *DouYinConnectionFanProfileClient {
	return &DouYinConnectionFanProfileClient{
		BaseClient: c,
	}
}

// ## GetUserFansData 获取用户粉丝数据
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/fans-portrait-data/get-user-fans-data
//
// 参数：
//
//	ctx  - 请求上下文
//
// 返回值：
//
//	*schema.DouYinConnectionFanProfileGetUserFansDataRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - FansData: 粉丝画像数据，包含活跃天数分布、年龄分布、设备分布等信息
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionFanProfileClient) GetUserFansData(ctx context.Context) (*schema.DouYinConnectionFanProfileGetUserFansDataRes, error) {
	result := &schema.DouYinConnectionFanProfileGetUserFansDataRes{}

	_, err := c.BaseClient.HttpGet(ctx, "/api/douyin/v1/user/fans_data/", nil, nil, nil, result)
	return result, err
}

// ## GetUserFansSource 获取用户粉丝来源
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/fans-portrait-data/get-user-fans-origin
//
// 参数：
//
//	ctx  - 请求上下文
//
// 返回值：
//
//	*schema.DouYinConnectionFanProfileFansSourceRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - List: 粉丝来源列表，包含来源类型及其占比
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionFanProfileClient) FansSource(ctx context.Context) (*schema.DouYinConnectionFanProfileFansSourceRes, error) {
	result := &schema.DouYinConnectionFanProfileFansSourceRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/data/extern/fans/source/", nil, nil, nil, result)
	return result, err
}

// ## FansFavourite 获取用户粉丝喜好
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/fans-portrait-data/get-user-fans-like
//
// 参数：
//
//	ctx  - 请求上下文
//
// 返回值：
//
//	*schema.DouYinConnectionFanProfileFansFavouriteRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - List: 粉丝喜好数据列表，包含关键词、排名和热度值
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionFanProfileClient) FansFavourite(ctx context.Context) (*schema.DouYinConnectionFanProfileFansFavouriteRes, error) {
	result := &schema.DouYinConnectionFanProfileFansFavouriteRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/data/extern/fans/favourite/", nil, nil, nil, result)
	return result, err
}

// ## FansComment 获取用户粉丝热评
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/fans-portrait-data/get-user-fans-hot-comment
//
// 参数：
//
//	ctx  - 请求上下文
//
// 返回值：
//
//	*schema.DouYinConnectionFanProfileFansCommentRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - List: 粉丝热评列表，包含评论内容和热度值
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//	error 调用过程中遇到的错误（如有）
func (c *DouYinConnectionFanProfileClient) FansComment(ctx context.Context) (*schema.DouYinConnectionFanProfileFansCommentRes, error) {
	result := &schema.DouYinConnectionFanProfileFansCommentRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/data/extern/fans/comment/", nil, nil, nil, result)
	return result, err
}
