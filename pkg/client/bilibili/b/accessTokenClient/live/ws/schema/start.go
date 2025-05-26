package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// WebSocketInfo 包含WebSocket连接所需的信息
type WebSocketInfo struct {
	AuthBody string   `json:"auth_body"` // WebSocket连接认证信息
	WssLink  []string `json:"wss_link"`  // WebSocket连接地址列表
}

// LiveWSStartData 包含直播WebSocket启动返回的数据
type LiveWSStartData struct {
	ConnID        string        `json:"conn_id"`        // 心跳ID (需要在30秒内使用心跳接口保持连接)
	WebSocketInfo WebSocketInfo `json:"websocket_info"` // WebSocket连接信息
}

// BiliBiliLiveWSStartRes 表示启动直播WebSocket连接的响应结构
// 接口文档参考：https://open.bilibili.com/doc/4/67eaa648-3f67-f2bc-0fac-efa5fb922305
type BiliBiliLiveWSStartRes struct {
	response.BiliBiliRes
	// LiveWSStartData 包含直播WebSocket启动返回的数据
	Data LiveWSStartData `json:"data"` // 响应数据
}
