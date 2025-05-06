package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMToolGetMessageResourcesReq 获取消息资源请求参数
type DouYinIMToolGetMessageResourcesReq struct {
	// ConversationId 会话ID
	ConversationId string `json:"conversation_id"`
	// MessageId 消息ID
	MessageId string `json:"message_id"`
}

// DouYinIMToolGetMessageResourcesRes 获取消息资源响应参数
type DouYinIMToolGetMessageResourcesRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// MediaType 媒体类型
		MediaType string `json:"media_type"`
		// Url 媒体URL
		Url string `json:"url"`
	} `json:"data"`
}
