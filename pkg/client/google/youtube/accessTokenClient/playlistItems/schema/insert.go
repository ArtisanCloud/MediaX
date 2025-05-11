package schema

// YouTubePlaylistItemsInsertReq 插入播放列表项请求参数
type YouTubePlaylistItemsInsertReq struct {
	Part                   string `json:"part"`                   // 指定返回的资源部分（必填，如 snippet,contentDetails 等）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	// PlaylistId             string `json:"playlistId"`             // 播放列表ID（可选，与 id 互斥）
	// ResourceId             ResourceId      `json:"resourceId"`             // 资源ID
	Snippet
	// Note             string `json:"note"`             // 备注
	// startAt          string `json:"startAt"`          // 开始时间
	// endAt            string `json:"endAt"`            // 结束时间
	ContentDetails
}

// YouTubePlaylistItemsInsertRes 插入播放列表项返回结果
type YouTubePlaylistItemsInsertRes struct {
	// PlaylistItem  播放列表项信息
	PlayListItem
}
