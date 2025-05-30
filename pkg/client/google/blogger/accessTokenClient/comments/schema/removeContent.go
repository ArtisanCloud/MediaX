package schema

// BloggerCommentsRemoveContentReq 表示 POST /blogs/{blogId}/posts/{postId}/comments/{commentId}/removecontent API 的请求参数
type BloggerCommentsRemoveContentReq struct {
	BlogId    string `json:"blogId"`    // 必填，博客的ID
	CommentId string `json:"commentId"` // 必填，要从中删除内容的评论ID
	PostId    string `json:"postId"`    // 必填，帖子的ID
}

// BloggerCommentsRemoveContentRes 表示 POST /blogs/{blogId}/posts/{postId}/comments/{commentId}/removecontent API 的响应
type BloggerCommentsRemoveContentRes struct {
	CommentResource
}
