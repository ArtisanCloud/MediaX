package schema

// BloggerPostsPatchReq 表示 PATCH /blogger/v3/blogs/{blogId}/posts/{postId} API 的请求参数
type BloggerPostsPatchReq struct {
	BlogId string       `json:"blogId"` // (必需) 博客的ID
	PostId string       `json:"postId"` // (必需) 帖子的ID
	Post   PostResource `json:"post"`   // (必需) 请求正文需要提供的Posts资源
}

// BloggerPostsPatchRes 表示 PATCH /blogger/v3/blogs/{blogId}/posts/{postId} API 的响应
type BloggerPostsPatchRes struct {
	Post PostResource
}
