package schema

// YoutubeCaptionsDeleteReq 定义删除YouTube字幕的请求结构
type YoutubeCaptionsDeleteReq struct {
	ID                     string `json:"id"`                               // id 参数用于指定要删除的字幕轨道ID（必填）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"` // 代表的内容所有者，仅适用于 YouTube 内容合作伙伴（可选）
}

// YoutubeCaptionsDeleteRes 定义删除YouTube字幕的响应结构
// 成功时返回 HTTP 204 No Content 状态代码
type YoutubeCaptionsDeleteRes struct{}
