package schema

// BloggerBlogsGetByUrlReq 表示 GET /blogger/v3/blogs/byurl API 的请求参数
type BloggerBlogsGetByUrlReq struct {
	URL string `json:"url"` // 必填，要检索的博客的网址
}
