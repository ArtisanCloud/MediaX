package schema

type YoutubeCommentsUpdateReq struct {
	Part    string         `json:"part"` // 指定返回的资源部分（必填，如 snippet）
	Snippet CommentSnippet `json:"snippet"`
}

type YoutubeCommentsUpdateRes struct {
	CommentSnippet
}
