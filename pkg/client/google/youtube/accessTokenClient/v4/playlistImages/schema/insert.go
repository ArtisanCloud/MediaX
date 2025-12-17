package schema

type YoutubePlaylistImagesInsertReq struct {
	Part                          string `json:"part"`                          // 指定返回的资源部分（必填，如 snippet,contentDetails 等）
	OnBehalfOfContentOwner        string `json:"onBehalfOfContentOwner"`        // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	OnBehalfOfContentOwnerChannel string `json:"onBehalfOfContentOwnerChannel"` // 内容所有者频道（可选，仅供 YouTube 内容合作伙伴使用）
	// Snippet.PlaylistId             string `json:"playlistId"`             // 播放列表ID（可选，与 id 互斥）
	// Snippet.Type                   string `json:"type"`                   // 播放列表项类型（可选）
	// Snippet.Width                  string `json:"width"`                  // 播放列表项宽度（可选）
	// Snippet.Height                 string `json:"height"`                 // 播放列表项高度（可选）
	Snippet
}

// PlaylistImages 播放列表图片
type YoutubePlaylistImagesInsertRes struct {
	// PlaylistImages 播放列表图片信息
	PlaylistImages
}
