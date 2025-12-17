package schema

// BloggerPagesPatchReq 表示 PATCH /blogger/v3/blogs/{blogId}/pages/{pageId} API 的请求参数
type BloggerPagesPatchReq struct {
	BlogId       string `json:"blogId"` // 必填，博客ID
	PageId       string `json:"pageId"` // 必填，页面ID
	PageResource        // 必填，要更新的页面资源
}

// BloggerPagesPatchRes 表示 PATCH /blogger/v3/blogs/{blogId}/pages/{pageId} API 的响应
type BloggerPagesPatchRes struct {
	PageResource
}
