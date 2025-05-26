package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliLiveWSHeartBeatReq represents the WebSocket heartbeat request
// Required fields:
// - ConnID: Heartbeat ID obtained from WebSocket start response (must be sent every 30 seconds)
type BiliBiliLiveWSHeartBeatReq struct {
	ClientId string `json:"client_id"` // Client ID 表示客户端的唯一标识符
	ConnID   string `json:"conn_id"`   // ConnID 表示连接的唯一标识符
}

// BiliBiliLiveWSHeartBeatRes represents the WebSocket heartbeat response
// Standard response format with code/message/data fields
type BiliBiliLiveWSHeartBeatRes struct {
	response.BiliBiliRes
	Data map[string]interface{} `json:"data"`
}
