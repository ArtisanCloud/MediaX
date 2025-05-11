package schema

// YoutubeCaptionsUpdateReq 定义更新YouTube字幕的请求参数
type YoutubeCaptionsUpdateReq struct {
	Part                   string  `json:"part"`                             // 指定返回的资源部分（必填，如 snippet）
	OnBehalfOfContentOwner string  `json:"onBehalfOfContentOwner,omitempty"` // 代表的内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	Caption                Caption `json:"caption"`                          // 字幕内容
}

// YoutubeCaptionsUpdateRes 定义更新YouTube字幕的返回结果
type YoutubeCaptionsUpdateRes struct {
	Caption // 更新后的字幕信息
}
