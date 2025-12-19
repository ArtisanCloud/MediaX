package schema

// YoutubeCommentThreadsInsertReq 表示 POST /youtube/v3/commentThreads API 的请求参数
type YoutubeCommentThreadsInsertReq struct {
	Part           string           `json:"part"` // 必填，指定 API 响应包含的属性
	CommentSnippet `json:"snippet"` // 评论的信息
}

// YoutubeCommentThreadsInsertRes 表示 POST /youtube/v3/commentThreads API 的响应
type YoutubeCommentThreadsInsertRes struct {
	Kind           string           `json:"kind"` // 资源类型
	Etag           string           `json:"etag"` // 资源的 ETag
	Id             string           `json:"id"`   // 评论会话ID
	CommentSnippet `json:"snippet"` // 评论会话的信息
}
