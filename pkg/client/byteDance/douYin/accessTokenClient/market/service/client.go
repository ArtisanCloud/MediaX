package service

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/market/service/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// DouYinMarketServiceClient 抖音市场服务类
type DouYinMarketServiceClient struct {
	*kernel.BaseClient
}

// NewClient 创建抖音市场服务类
func NewClient(c *kernel.BaseClient) *DouYinMarketServiceClient {
	return &DouYinMarketServiceClient{
		BaseClient: c,
	}
}

// ## GetUserServicePurchaseList 查询用户的服务购买信息
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/service-market/service-relationship/get-user-service-purchase-list
//
// 参数：
//   ctx context.Context 请求上下文
//   serviceId string 服务ID
//
// 返回值：
//   *schema.DouYinMarketServiceGetUserServicePurchaseListRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - PurchaseInfoList: 购买信息列表
//   error 调用过程中遇到的错误（如有）
func (c *DouYinMarketServiceClient) GetUserServicePurchaseList(ctx context.Context, serviceId string) (*schema.DouYinMarketServiceGetUserServicePurchaseListRes, error) {
	result := &schema.DouYinMarketServiceGetUserServicePurchaseListRes{}
	params := &object.StringMap{
		"serviceId": serviceId,
	}
	_, err := c.BaseClient.HttpGet(ctx, "/market/service/user/purchase/list/", params, nil, result)
	return result, err
}

// ## DecrUserRemainTimes 扣除用户服务的剩余使用次数/条数
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/service-market/service-relationship/decr-product-user-remain-times
//
// 参数：
//   ctx context.Context 请求上下文
//   data *schema.DouYinMarketServiceDecrUserRemainTimesReq 扣除请求参数，包含以下字段：
//     • Count: 要扣除的次数/条数，要求大于0
//     • IsTestEnv: 是否为测试环境
//     • OpenID: 查询订阅信息目标用户的唯一标识  *** 注意，可以根据Token对象，获取到授权的openId ***
//     • OutTradeNo: 外部开发者单号id，用于保证当前扣除操作的幂等
//     • ServiceID: 服务id，服务的唯一标识
//     • ServiceModeID: 服务规格id，由服务商在创建服务时指定
//
// 返回值：
//   *schema.DouYinMarketServiceDecrUserRemainTimesRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinMarketServiceClient) DecrUserRemainTimes(ctx context.Context, data *schema.DouYinMarketServiceDecrUserRemainTimesReq) (*schema.DouYinMarketServiceDecrUserRemainTimesRes, error) {
	result := &schema.DouYinMarketServiceDecrUserRemainTimesRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/market/service/user/remaintimes/decr/", nil, data, nil, result)
	return result, err

}

// ## InsertPurchaseInfo 导入外部订购数据
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/service-market/service-relationship/InsertPurchaseInfo
//
// 参数：
//   ctx context.Context 请求上下文
//   data *schema.DouYinMarketServiceInsertPurchaseInfoReq 导入请求参数，包含以下字段：
//     • OpenID: 查询订阅信息目标用户的唯一标识
//     • OutTradeNo: 外部开发者单号id，用于唯一标识同一服务同一规格当前导入的订购数据
//     • PeriodType: 服务规格的周期类型
//     • PurchaseTime: 服务订购时间的毫秒级时间戳
//     • ServiceID: 服务id，服务的唯一标识
//     • ServiceModeID: 用户订购的服务规格id
//     • Duration: 订购时长值，适用于时间类型的服务规格
//     • UsageTimes: 用户在服务商侧购买的服务有效使用次数/条数
//
// 返回值：
//   *schema.DouYinMarketServiceInsertPurchaseInfoRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinMarketServiceClient) InsertPurchaseInfo(ctx context.Context, data *schema.DouYinMarketServiceInsertPurchaseInfoReq) (*schema.DouYinMarketServiceInsertPurchaseInfoRes, error) {
	result := &schema.DouYinMarketServiceInsertPurchaseInfoRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/market/service/user/insert/purchase/info/", nil, data, nil, result)
	return result, err
}

// ## DeletePurchaseInfo 删除导入的外部订购数据
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/service-market/service-relationship/DeletePurchaseInfo
//
// 参数：
//   ctx context.Context 请求上下文
//   data *schema.DouYinMarketServiceDeletePurchaseInfoReq 删除请求参数，包含以下字段：
//     • OpenID: 查询订阅信息目标用户的唯一标识
//     • OutTradeNo: 外部开发者单号id，用于唯一标识同一服务同一规格下要删除的订购数据
//     • ServiceID: 服务id，服务的唯一标识
//     • ServiceModeID: 用户订购的服务规格id
//
// 返回值：
//   *schema.DouYinMarketServiceDeletePurchaseInfoRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinMarketServiceClient) DeletePurchaseInfo(ctx context.Context, data *schema.DouYinMarketServiceDeletePurchaseInfoReq) (*schema.DouYinMarketServiceDeletePurchaseInfoRes, error) {
	result := &schema.DouYinMarketServiceDeletePurchaseInfoRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/market/service/user/delete/purchase/info/", nil, data, nil, result)
	return result, err
}
