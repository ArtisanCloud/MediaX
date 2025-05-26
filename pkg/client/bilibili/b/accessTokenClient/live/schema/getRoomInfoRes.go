package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type LiveGetRoomInfoData struct {
	OpenId      string `json:"open_id"`      // 用户OpenID
	RoomId      int    `json:"room_id"`      // 直播房间ID
	Title       string `json:"title"`        // 直播房间标题
	IsStreaming bool   `json:"is_streaming"` // 当前是否开播
	IsBanned    bool   `json:"is_banned"`    // 房间是否被封禁
}

// BiliBiliLiveGetRoomInfoRes 获取直播房间信息响应
// 接口文档参考：https://member.bilibili.com/arcopen/fn/live/room/info
type BiliBiliLiveGetRoomInfoRes struct {
	response.BiliBiliRes
	Data LiveGetRoomInfoData `json:"data"`
}
