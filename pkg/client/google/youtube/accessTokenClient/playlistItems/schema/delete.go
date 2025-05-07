package schema

// YouTubePlaylistItemsDeleteReq 删除播放列表项请求参数
type YouTubePlaylistItemsDeleteReq struct {
	Id                     string `json:"id"`                     // 播放列表项ID
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
}

// YouTubePlaylistItemsDeleteRes 删除播放列表项返回结果 HTTP 204 返回码
type YouTubePlaylistItemsDeleteRes struct {
}
