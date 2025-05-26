package ws

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/live/ws/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type BiliBiliLiveWSClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliLiveWSClient {
	return &BiliBiliLiveWSClient{
		BaseClient: c,
	}
}

// ## Start 启动直播WebSocket连接
//
// 接口文档参考：
// `https://open.bilibili.com/doc/4/da2b13dc-7f7a-0025-be11-0b677e793baa#h1--id`
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*schema.BiliBiliLiveWSStartRes 包含以下字段：
//	  • Code: 状态码
//	  • Message: 消息
//	  • Data: 连接数据，包含以下字段：
//	    - ConnID: 心跳ID (需要在30秒内使用心跳接口保持连接)
//	    - WebSocketInfo: 连接地址
//		error 调用过程中遇到的错误（如有）
func (c *BiliBiliLiveWSClient) Start(ctx context.Context) (*schema.BiliBiliLiveWSStartRes, error) {
	result := &schema.BiliBiliLiveWSStartRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/live/room/ws-start", nil, nil, nil, result)
	return result, err
}

// ## HeartBeat 发送直播房间长链心跳
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/8f57c906-15e6-4e61-e850-da0911be1cf1#h1-u76F4u64ADu623Fu95F4u957Fu94FEu5FC3u8DF3
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • ClientId: 申请应用时获得的client_id(必填)
//	  • ConnID: 获取长链地址时得到的心跳ID(必填)
//
// 返回值：
//
//	*schema.BiliBiliLiveWSHeartBeatRes 包含以下字段：
//	  • Code: 状态码(0表示成功)
//	  • Message: 状态信息
//	  • Data: 响应数据(通常为空)
//	error 调用过程中遇到的错误(如有)
//
// 注意事项：
// 1. 必须每30秒发送一次心跳以保持连接
// 2. 请求URL为https://member.bilibili.com/arcopen/fn/live/room/ws-heartbeat
// 3. 请求头需包含Host: member.bilibili.com和Connection: keep-alive
func (c *BiliBiliLiveWSClient) HeartBeat(ctx context.Context, data *schema.BiliBiliLiveWSHeartBeatReq) (*schema.BiliBiliLiveWSHeartBeatRes, error) {
	result := &schema.BiliBiliLiveWSHeartBeatRes{}
	params := &object.StringMap{
		"client_id": data.ClientId,
	}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/live/room/ws-heartbeat", params, data, nil, result)
	return result, err
}

// ## BatchHeartBeat 发送直播房间批量长链心跳
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/6ea09a64-8ac2-b650-9580-866850432c2c#h1-u76F4u64ADu623Fu95F4u6279u91CFu957Fu94FEu5FC3u8DF3
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • ConnIDs: 心跳ID列表(必填)
//
// 返回值：
//
//	*schema.BiliBiliLiveWSBatchHeartBeatRes 包含以下字段：
//	  • Code: 状态码(0表示成功)
//	  • Message: 状态信息
//	  • Data.FailedConnIDs: 心跳失败的心跳ID列表
//	error 调用过程中遇到的错误(如有)
//
// 注意事项：
// 1. 必须每30秒发送一次心跳以保持连接
// 2. 请求URL为https://member.bilibili.com/arcopen/fn/live/room/ws-batch-heartbeat
// 3. 请求头需包含Host: member.bilibili.com和Connection: keep-alive
func (c *BiliBiliLiveWSClient) BatchHeartBeat(ctx context.Context, data *schema.BiliBiliLiveWSBatchHeartBeatReq) (*schema.BiliBiliLiveWSBatchHeartBeatRes, error) {
	result := &schema.BiliBiliLiveWSBatchHeartBeatRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/live/room/ws-batch-heartbeat", nil, data, nil, result)
	return result, err
}
