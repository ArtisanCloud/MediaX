package schema

// BloggerPagesInsertReq 表示 POST /blogger/v3/blogs/{blogId}/pages API 的请求参数
type BloggerPagesInsertReq struct {
	BlogId string `json:"blogId"` // 必填，要将页面添加到的博客ID
	PageResource
}

// BloggerPagesInsertRes 表示 POST /blogger/v3/blogs/{blogId}/pages API 的响应
type BloggerPagesInsertRes struct {
	PageResource
}
