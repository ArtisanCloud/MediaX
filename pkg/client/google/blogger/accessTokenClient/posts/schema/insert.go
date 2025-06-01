package schema

// BloggerPostsInsertReq 表示 POST /blogger/v3/blogs/{blogId}/posts API 的请求参数
type BloggerPostsInsertReq struct {
	BlogId       string `json:"blogId"`  // 必填，要向其添加博文的博客的 ID
	IsDraft      *bool  `json:"isDraft"` // 可选，是否以草稿形式创建博文
	PostResource        // 嵌入 Post 结构体作为请求正文
}

// BloggerPostsInsertRes 表示 POST /blogger/v3/blogs/{blogId}/posts API 的响应
type BloggerPostsInsertRes struct {
	PostResource // 嵌入 Post 结构体作为响应正文
}
