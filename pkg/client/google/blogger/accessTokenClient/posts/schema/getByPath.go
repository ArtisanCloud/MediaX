package schema

// BloggerPostsGetByPathReq 表示 GET /blogger/v3/blogs/{blogId}/posts/bypath API 的请求参数
type BloggerPostsGetByPathReq struct {
	BlogId      string  `json:"blogId"`                // 必需，要从中提取博文的博客的ID
	Path        string  `json:"path"`                  // 必需，要检索的博文的路径
	MaxComments *uint   `json:"maxComments,omitempty"` // 可选，为帖子检索的评论数量上限
	View        *string `json:"view,omitempty"`        // 可选，视图级别（ADMIN/AUTHOR/READER）
}

// BloggerPostsGetByPathRes 表示 GET /blogger/v3/blogs/{blogId}/posts/bypath API 的响应
type BloggerPostsGetByPathRes struct {
	PostResource
}
