package schema

// YoutubePlaylistImagesUpdateReq 更新播放列表图片请求参数
type YoutubePlaylistImagesUpdateReq struct {
	Part                   string `json:"part"`                             // 指定返回的资源部分（必填，如 snippet 等）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	// Snippet.PlaylistId             string `json:"playlistId"`             // 播放列表ID（可选，与 id 互斥）
	// Snippet.Type       string `json:"type"`       // 类型（可选）
	// Snippet.Width      string `json:"width"`      // 宽度（可选）
	// Snippet.Height     string `json:"height"`     // 高度（可选）
	Snippet
}

// YoutubePlaylistImagesUpdateRes 更新播放列表图片返回结果
type YoutubePlaylistImagesUpdateRes struct {
	// PlaylistItem  播放列表项信息
	PlaylistImages
}
