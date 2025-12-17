package schema

// YouTubePlaylistsUpdateReq 更新播放列表请求参数
type YouTubePlaylistsUpdateReq struct {
	Part                   string `json:"part"`                   // 指定返回的资源部分（必填，如 snippet,status 等）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	Playlist                      // 播放列表信息
}

// YouTubePlaylistsUpdateRes 更新播放列表返回结果
type YouTubePlaylistsUpdateRes struct {
	Playlist // 播放列表信息
}
