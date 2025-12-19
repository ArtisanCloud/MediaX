package schema

type Position struct {
	Type           string `json:"type"`
	CornerPosition string `json:"cornerPosition"`
}

type Timing struct {
	Type       string `json:"type"`
	OffsetMs   int    `json:"offsetMs"`
	DurationMs int    `json:"durationMs"`
}

type Watermark struct {
	Timing          Timing   `json:"timing"`
	Position        Position `json:"position"`
	ImageUrl        string   `json:"imageUrl"`
	ImageBytes      string   `json:"imageBytes"`
	TargetChannelId string   `json:"targetChannelId"`
}

// YouTubeWatermarksSetReq 设置水印请求参数
type YouTubeWatermarksSetReq struct {
	ChannelId              string `json:"channelId"`                        // 频道ID（必填）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	Watermark
}

type YouTubeWatermarksSetRes struct {
}
