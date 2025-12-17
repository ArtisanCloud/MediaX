package schema

// BloggerPagesDeleteReq 表示 DELETE /blogger/v3/blogs/{blogId}/pages/{pageId} API 的请求参数
type BloggerPagesDeleteReq struct {
	BlogId   string `json:"blogId"`             // 必填，博客的ID
	PageId   string `json:"pageId"`             // 必填，页面的ID
	UseTrash bool   `json:"useTrash,omitempty"` // 可选，是否将页面移至回收站
}

// BloggerPagesDeleteRes 表示 DELETE /blogger/v3/blogs/{blogId}/pages/{pageId} API 的响应
type BloggerPagesDeleteRes struct {
	// 成功时返回空响应体
}
