package schema

type BrandingSettings struct {
	ChannelItem `json:"channel"`
}

// YoutubeChannelsUpdateReq 更新频道信息的请求参数
type YoutubeChannelsUpdateReq struct {
	// Part 指定返回的资源部分（必填，如 brandingSettings）
	Part string `json:"part"`
	// OnBehalfOfContentOwner 代表内容所有者执行操作（可选，仅供 YouTube 内容合作伙伴使用）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"`
	// BrandingSettings 频道品牌设置信息（可选）
	BrandingSettings BrandingSettings `json:"brandingSettings,omitempty"`
}
type YoutubeChannelsUpdateRes struct {
	ChannelItem
}
