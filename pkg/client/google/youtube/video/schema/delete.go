package schema

// YouTubeVideoDeleteReq 定义删除 YouTube 视频的请求结构
type YouTubeVideoDeleteReq struct {
	ID                     string `json:"id"`                               // id 参数用于指定要删除的资源的 YouTube 视频 ID
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"` // 代表的内容所有者，仅适用于 YouTube 内容合作伙伴
}

// YouTubeVideoDeleteRes 定义删除 YouTube 视频的响应结构
// 成功时返回 HTTP 204 响应代码，这里可留空结构体
type YouTubeVideoDeleteRes struct{}
