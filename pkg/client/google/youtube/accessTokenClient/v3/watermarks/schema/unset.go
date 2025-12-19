package schema

// YouTubeWatermarksUnsetReq 视频水印解除请求参数
type YouTubeWatermarksUnsetReq struct {
	ChannelId              string `json:"channelId"`                        // 频道ID（必填）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
}

// YouTubeWatermarksUnsetRes 视频水印解除返回结果 204 No Content
type YouTubeWatermarksUnsetRes struct {
}
