package schema

// BloggerCommentsMarkAsSpamReq 表示 POST /blogs/{blogId}/posts/{postId}/comments/{commentId}/spam API 的请求参数
type BloggerCommentsMarkAsSpamReq struct {
	BlogId    string `json:"blogId"`    // 必填，博客的ID
	CommentId string `json:"commentId"` // 必填，要标记为垃圾内容的评论ID
	PostId    string `json:"postId"`    // 必填，帖子的ID
}

// BloggerCommentsMarkAsSpamRes 表示 POST /blogs/{blogId}/posts/{postId}/comments/{commentId}/spam API 的响应
type BloggerCommentsMarkAsSpamRes struct {
	CommentResource
}
