package thirdParty

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/live/thirdParty/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type BiliBiliLiveThirdPartyClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliLiveThirdPartyClient {
	return &BiliBiliLiveThirdPartyClient{
		BaseClient: c,
	}
}

// ## GetLiveGrantUrl 获取开播授权链接
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/5827c4a4-aab6-235e-624b-a47248d712e3#h2-u83B7u53D6u5F00u64ADu6388u6743u94FEu63A5
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BizCode: 申请接入时获得的biz_code(必填)
//	  • OpenID: 开播用户的open_id(必填)
//	  • LiveAreaID: 目标开播分区(必填)
//	  • ThirdLiveUUID: 第三方提供本次开播的唯一ID(必填)
//
// 返回值：
//
//	*schema.BiliBiliLiveThirdPartyGetLiveGrantUrlRes 包含以下字段：
//	  • Code: 状态码(0表示成功)
//	  • Message: 状态信息
//	  • Data.GrantURL: 授权链接地址
//	error 调用过程中遇到的错误(如有)
//
// 注意事项：
// 1. 需要申请LIVE_ROOM_DATA权限
// 2. 请求URL为https://member.bilibili.com/liveopen/fn/live/thirdPartyLive/grantUrl
// 3. 需要添加公共请求头(Request Header)和签名
// 4. 返回的授权地址需要生成为二维码供用户扫码
func (c *BiliBiliLiveThirdPartyClient) GetLiveGrantUrl(ctx context.Context, data *schema.BiliBiliLiveThirdPartyGetLiveGrantUrlReq) (*schema.BiliBiliLiveThirdPartyGetLiveGrantUrlRes, error) {
	result := &schema.BiliBiliLiveThirdPartyGetLiveGrantUrlRes{}
	params, _ := object.StructToStringMap(data)
	_, err := c.BaseClient.HttpGet(ctx, "/liveopen/fn/live/thirdPartyLive/grantUrl", params, nil, nil, result)
	return result, err
}

// ## StartLive 开始第三方直播
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/cafd6cf4-be34-d2c6-4d43-ca309b7ad346#h2-u5F00u59CBu76F4u64AD
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BizCode: 业务唯一固定码(必填)
//	  • OpenID: 用户唯一ID(必填)
//	  • LiveToken: 一次性开播token(必填)
//
// 返回值：
//
//	*schema.BiliBiliLiveThirdPartyStartLiveRes 包含以下字段：
//	  • Code: 状态码(0表示成功)
//	  • Message: 错误信息
//	  • Data.PushURL: 推流地址
//	  • Data.LiveKey: 本次开播唯一ID
//	  • Data.StartTime: 开播时间戳(UTC+8)
//	error 调用过程中遇到的错误(如有)
//
// 注意事项：
// 1. 需要先调用grantUrl接口获取授权
// 2. 需要LIVE_ROOM_DATA权限
// 3. 请求头必须包含签名信息(版本2.0)
func (c *BiliBiliLiveThirdPartyClient) StartLive(ctx context.Context, data *schema.BiliBiliLiveThirdPartyStartLiveReq) (*schema.BiliBiliLiveThirdPartyStartLiveRes, error) {
	result := &schema.BiliBiliLiveThirdPartyStartLiveRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/liveopen/fn/live/thirdPartyLive/startLive", nil, data, nil, result)
	return result, err
}

// ## EndLive 结束第三方直播
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/e84cd5ba-68fe-effe-1835-c2a2ff22a72a#h2-u7ED3u675Fu76F4u64AD
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • BizCode: 业务唯一固定码(必填)
//	  • OpenID: 用户唯一ID(必填)
//	  • LiveKey: 开播后返回的直播live_key(必填)
//	  • ThirdLiveUUID: 本次开播调用grantUrl时给到的唯一id(必填)
//
// 返回值：
//
//	*schema.BiliBiliLiveThirdPartyEndLiveRes 包含以下字段：
//	  • Code: 状态码(0表示成功)
//	  • Message: 错误信息
//	  • Data.OpenID: 直播用户唯一ID
//	  • Data.EndTime: 停止直播的时间戳(UTC+8)
//	error 调用过程中遇到的错误(如有)
//
// 注意事项：
// 1. 需要LIVE_ROOM_DATA权限
// 2. 请求头必须包含签名信息(版本2.0)
// 3. 调用前需确保直播已通过StartLive接口成功开启
func (c *BiliBiliLiveThirdPartyClient) EndLive(ctx context.Context, data *schema.BiliBiliLiveThirdPartyEndLiveReq) (*schema.BiliBiliLiveThirdPartyEndLiveRes, error) {
	result := &schema.BiliBiliLiveThirdPartyEndLiveRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/liveopen/fn/live/thirdPartyLive/endLive", nil, data, nil, result)
	return result, err
}
