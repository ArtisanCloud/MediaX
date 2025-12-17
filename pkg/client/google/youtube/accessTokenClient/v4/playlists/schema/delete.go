package schema

type YouTubePlaylistsDeleteReq struct {
	Id                     string `json:"id"`                     // 播放列表ID（必填）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
}

// YouTubePlaylistsDeleteRes HTTP 204 返回码
type YouTubePlaylistsDeleteRes struct {
}
