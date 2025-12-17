package schema

// YoutubeCaptionsInsertReq 表示 POST /youtube/v3/captions API 的请求参数
type YoutubeCaptionsInsertReq struct {
	Part                   string  `json:"part"`                             // 必填，指定 API 响应包含的属性
	OnBehalfOfContentOwner string  `json:"onBehalfOfContentOwner,omitempty"` // 可选，代表内容所有者执行操作
	Caption                Caption `json:"caption"`                          // 字幕信息
}

// YoutubeCaptionsInsertRes 表示 POST /youtube/v3/captions API 的响应
type YoutubeCaptionsInsertRes struct {
	Caption // 字幕信息
}
