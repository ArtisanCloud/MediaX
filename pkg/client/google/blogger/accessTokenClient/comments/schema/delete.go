package schema

// BloggerCommentsDeleteReq 表示删除评论的请求参数
// 接口文档参考：https://developers.google.com/blogger/docs/3.0/reference/comments/delete
//
// 参数：
//
//	BlogId - 博客ID (必填)
//	CommentId - 要删除的评论ID (必填)
//	PostId - 帖子ID (必填)
type BloggerCommentsDeleteReq struct {
	BlogId    string `json:"blogId"`    // 博客ID
	CommentId string `json:"commentId"` // 评论ID
	PostId    string `json:"postId"`    // 帖子ID
}

// BloggerCommentsDeleteRes 表示删除评论的响应
// 成功时返回空响应
type BloggerCommentsDeleteRes struct{}
