package schema

// YoutubeCommentsInsertReq 表示 POST /youtube/v3/comments API 的请求参数
type YoutubeCommentsInsertReq struct {
	Part    string         `json:"part"`    // 必填，指定 API 响应包含的属性
	Snippet CommentSnippet `json:"snippet"` // 评论资源信息
}

// YoutubeCommentsInsertRes 表示 POST /youtube/v3/comments API 的响应
type YoutubeCommentsInsertRes struct {
	Kind    string         `json:"kind"`    // 资源类型
	Etag    string         `json:"etag"`    // 资源的 ETag
	Id      string         `json:"id"`      // 评论ID
	Snippet CommentSnippet `json:"snippet"` // 评论信息
}
