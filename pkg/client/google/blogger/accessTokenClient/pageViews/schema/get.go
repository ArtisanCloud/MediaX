package schema

// BloggerPageViewsGetReq 表示 GET /blogger/v3/blogs/{blogId}/pageviews API 的请求参数
type BloggerPageViewsGetReq struct {
	BlogId string `json:"blogId"` // 必填，要获取的博客的ID
	Range  string `json:"range"`  // 可选，指定要检索的时间范围，可选值："30DAYS"(过去30天)、"7DAYS"(过去7天)、"all"(所有页面)
}

type PageViewCount struct {
	TimeRange string `json:"timeRange"` // 指定计数适用的时间范围
	Count     int64  `json:"count"`     // 给定时间范围内的网页浏览量
}

type PageViewResource struct {
	Kind   string          `json:"kind"`   // 资源类型，固定值：blogger#page_views
	BlogId int64           `json:"blogId"` // 博客ID
	Counts []PageViewCount `json:"counts"` // 此博客中帖子的容器
}

// BloggerPageViewsGetRes 表示 GET /blogger/v3/blogs/{blogId}/pageviews API 的响应
type BloggerPageViewsGetRes struct {
	PageViewResource
}
