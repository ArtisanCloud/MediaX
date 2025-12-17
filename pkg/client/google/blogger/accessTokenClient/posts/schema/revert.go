package schema

import "time"

// BloggerPostsRevertReq 表示 POST /blogger/v3/blogs/{blogId}/posts/{postId}/revert API 的请求参数
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/posts/revert?hl=zh-cn
type BloggerPostsRevertReq struct {
	BlogId     string     `json:"blogId"`               // (必需) 博客ID
	PostId     string     `json:"postId"`               // (必需) 帖子ID
	RevertDate *time.Time `json:"RevertDate,omitempty"` // (可选) 安排发布的日期时间
}

// BloggerPostsRevertRes 表示 POST /blogger/v3/blogs/{blogId}/posts/{postId}/revert API 的响应
// 响应正文包含一个 Posts 资源
type BloggerPostsRevertRes struct {
	PostResource // 嵌入Post结构体
}
