package live

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/live/schema"
)

type BiliBiliLiveClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *BiliBiliLiveClient {
	return &BiliBiliLiveClient{
		BaseClient: c,
	}
}

// ## GetRoomInfo 获取直播间基础信息
//
// 接口文档参考：
// `https://open.bilibili.com/doc/4/67eaa648-3f67-f2bc-0fac-efa5fb922305#h1-u83B7u53D6u7528u6237u6570u636E`
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*schema.BiliBiliLiveGetRoomInfoRes 包含以下字段：
//	  • open_id: 用户OpenID
//	  • room_id: 直播房间ID
//	  • title: 直播房间标题
//	  • is_streaming: 当前是否开播
//	  • is_banned: 房间是否被封禁
//	error 调用过程中遇到的错误（如有）
func (c *BiliBiliLiveClient) GetRoomInfo(ctx context.Context) (*schema.BiliBiliLiveGetRoomInfoRes, error) {
	result := &schema.BiliBiliLiveGetRoomInfoRes{}
	_, err := c.BaseClient.HttpGet(ctx, "arcopen/fn/live/room/info", nil, nil, nil, result)
	return result, err
}
