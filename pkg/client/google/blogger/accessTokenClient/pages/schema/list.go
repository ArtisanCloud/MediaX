package schema

// BloggerPagesListReq 表示 GET /blogger/v3/blogs/{blogId}/pages API 的请求参数
type BloggerPagesListReq struct {
	BlogId      string  `json:"blogId"`                // 必填，要从中提取网页的博客的ID
	FetchBodies *bool   `json:"fetchBodies,omitempty"` // 可选，是否检索网页正文
	Status      *string `json:"status,omitempty"`      // 可选，页面状态(draft/imported/live)
	View        *string `json:"view,omitempty"`        // 可选，视图级别(ADMIN/AUTHOR/READER)
}

// BloggerPagesListRes 表示 GET /blogger/v3/blogs/{blogId}/pages API 的响应
type BloggerPagesListRes struct {
	Kind  string         `json:"kind"`  // 资源类型，始终为"blogger#pageList"
	Items []PageResource `json:"items"` // 页面资源列表
}

// Page 表示 Blogger API 返回的页面资源
type PageResource struct {
	Kind      string `json:"kind"`      // 资源类型，始终为"blogger#page"
	ID        string `json:"id"`        // 页面唯一标识符
	Status    string `json:"status"`    // 页面状态(draft/imported/live)
	Blog      Blog   `json:"blog"`      // 所属博客信息
	Published string `json:"published"` // 发布时间(RFC 3339格式)
	Updated   string `json:"updated"`   // 最后更新时间(RFC 3339格式)
	URL       string `json:"url"`       // 页面公开URL
	SelfLink  string `json:"selfLink"`  // 页面API资源URL
	Title     string `json:"title"`     // 页面标题
	Content   string `json:"content"`   // 页面内容(HTML格式)
	Author    Author `json:"author"`    // 作者信息
}

// Blog 表示页面所属的博客信息
type Blog struct {
	ID string `json:"id"` // 博客ID
}

// Author 表示页面作者信息
type Author struct {
	ID          string `json:"id"`          // 作者ID
	DisplayName string `json:"displayName"` // 作者显示名称
	URL         string `json:"url"`         // 作者个人资料URL
	Image       Image  `json:"image"`       // 作者头像信息
}

// Image 表示作者头像信息
type Image struct {
	URL string `json:"url"` // 头像图片URL
}
