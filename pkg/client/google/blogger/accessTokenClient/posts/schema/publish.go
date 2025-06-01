package schema

import "time"

// BloggerPostsPublishReq 表示 POST /blogger/v3/blogs/{blogId}/posts/{postId}/publish API 的请求参数
// 参考文档：https://developers.google.com/blogger/docs/3.0/reference/posts/publish?hl=zh-cn
type BloggerPostsPublishReq struct {
	BlogId      string     `json:"blogId"`                // (必需) 博客ID
	PostId      string     `json:"postId"`                // (必需) 帖子ID
	PublishDate *time.Time `json:"publishDate,omitempty"` // (可选) 安排发布的日期时间
}

// BloggerPostsPublishRes 表示 POST /blogger/v3/blogs/{blogId}/posts/{postId}/publish API 的响应
// 响应正文包含一个 Posts 资源
type BloggerPostsPublishRes struct {
	PostResource // 嵌入Post结构体
}
