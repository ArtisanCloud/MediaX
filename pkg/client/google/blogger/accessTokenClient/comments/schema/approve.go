package schema

// BloggerCommentsApproveReq 表示批准评论的请求参数
type BloggerCommentsApproveReq struct {
	BlogId    string `json:"blogId"`    // 必需，博客ID
	CommentId string `json:"commentId"` // 必需，要批准的评论ID
	PostId    string `json:"postId"`    // 必需，帖子ID
}

// BloggerCommentsApproveRes 表示批准评论的响应结构
type BloggerCommentsApproveRes struct {
	CommentResource
}
