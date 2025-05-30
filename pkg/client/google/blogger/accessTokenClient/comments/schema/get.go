package schema

type BloggerCommentsGetReq struct {
	BlogId    string `json:"blogId"`         // 必需，包含评论的博客ID
	CommentId string `json:"commentId"`      // 必需，要获取的评论ID
	PostId    string `json:"postId"`         // 必需，所属帖子ID
	View      string `json:"view,omitempty"` // 可选，视图级别（ADMIN/AUTHOR/READER）
}

type BloggerCommentsGetRes struct {
	CommentResource
}
