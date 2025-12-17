package schema

// YoutubePlaylistImagesDeleteReq 删除播放列表图片请求参数
type YoutubePlaylistImagesDeleteReq struct {
	Part                   string `json:"part"`                   // 指定返回的资源部分（必填，如 snippet,contentDetails 等）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
}

// YoutubePlaylistImagesDeleteRes 删除播放列表图片返回结果 HTTP 204 返回码
type YoutubePlaylistImagesDeleteRes struct {
}
