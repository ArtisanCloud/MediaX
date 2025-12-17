package schema

// BloggerPostsGetReq 表示 GET /blogger/v3/blogs/{blogId}/posts/{postId} API 的请求参数
type BloggerPostsGetReq struct {
	BlogId      string  `json:"blogId"`      // 必填，要从中提取博文的博客的ID
	PostId      string  `json:"postId"`      // 必填，帖子的ID
	MaxComments *uint32 `json:"maxComments"` // 可选，作为帖子资源一部分要检索的评论数量上限
	View        *string `json:"view"`        // 可选，详细信息级别：ADMIN（管理员级别）、AUTHOR（作者级别）、READER（读者级别）
}

// BloggerPostsGetRes 表示 GET /blogger/v3/blogs/{blogId}/posts/{postId} API 的响应
type BloggerPostsGetRes struct {
	PostResource // 继承博客文章资源的所有字段
}
