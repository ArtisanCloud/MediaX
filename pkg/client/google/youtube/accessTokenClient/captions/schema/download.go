package schema

// YoutubeCaptionsDownloadReq 表示 GET /youtube/v3/captions API 的请求参数
type YoutubeCaptionsDownloadReq struct {
	ID                     string  `json:"id"`                               // 必填，指定要检索的字幕轨道ID
	OnBehalfOfContentOwner *string `json:"onBehalfOfContentOwner,omitempty"` // 可选，代表内容所有者执行操作
	Tfmt                   *string `json:"tfmt,omitempty"`                   // 可选，指定字幕格式
	Tlang                  *string `json:"tlang,omitempty"`                  // 可选，指定字幕语言
}

// YoutubeCaptionsDownloadRes 表示 GET /youtube/v3/captions API 的响应
type YoutubeCaptionsDownloadRes struct {
	ContentType string `json:"-"` // 响应头 Content-Type
	Data        []byte `json:"-"` // 二进制文件数据
}
