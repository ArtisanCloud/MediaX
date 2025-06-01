package schema

// BloggerPostsDeleteReq 表示 DELETE /blogger/v3/blogs/{blogId}/posts/{postId} API 的请求参数
type BloggerPostsDeleteReq struct {
	BlogId   string `json:"blogId"`   // 必填，博客的 ID
	PostId   string `json:"postId"`   // 必填，帖子的 ID
	UseTrash *bool  `json:"useTrash"` // 可选，如有可能，请移至回收站
}

// BloggerPostsDeleteRes 表示 DELETE /blogger/v3/blogs/{blogId}/posts/{postId} API 的响应
type BloggerPostsDeleteRes struct {
	// 空响应体
}
