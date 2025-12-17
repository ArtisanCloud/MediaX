package schema

// BloggerPagesUpdateReq 表示 PUT /blogger/v3/blogs/{blogId}/pages/{pageId} API 的请求参数
type BloggerPagesUpdateReq struct {
	BlogId       string `json:"blogId"` // 必填，博客ID
	PageId       string `json:"pageId"` // 必填，页面ID
	PageResource        // 必填，要更新的页面资源
}

// BloggerPagesUpdateRes 表示 PUT /blogger/v3/blogs/{blogId}/pages/{pageId} API 的响应
type BloggerPagesUpdateRes struct {
	PageResource
}
