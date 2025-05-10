package schema

// YoutubeCommentsDeleteReq 定义删除 YouTube 评论的请求结构
type YoutubeCommentsDeleteReq struct {
	ID string `json:"id"` // 评论ID（必填）
}

// YoutubeCommentsDeleteRes 定义删除 YouTube 评论的响应结构
// 成功时返回 HTTP 204 响应代码，这里可留空结构体
type YoutubeCommentsDeleteRes struct{}
