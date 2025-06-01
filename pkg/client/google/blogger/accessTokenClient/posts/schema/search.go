package schema

// BloggerPostsSearchReq 表示 GET /blogger/v3/blogs/{blogId}/posts/search API 的请求参数
type BloggerPostsSearchReq struct {
	BlogId      string  `json:"blogId"`      // 必填，要在其中进行搜索的博客的 ID
	Q           string  `json:"q"`           // 必填，要搜索的查询字词
	FetchBodies *bool   `json:"fetchBodies"` // 可选，是否包含帖子的正文内容（默认值：true）
	OrderBy     *string `json:"orderBy"`     // 可选，应用于搜索结果的排序顺序（可接受的值："published"、"updated"）
}

// BloggerPostsSearchRes 表示 GET /blogger/v3/blogs/{blogId}/posts/search API 的响应
type BloggerPostsSearchRes struct {
	Kind          string         `json:"kind"`          // 此实体的种类，始终为 blogger#postList
	NextPageToken string         `json:"nextPageToken"` // 用于提取下一页的分页令牌（如果存在）
	Items         []PostResource `json:"items"`         // 此博客的博文列表
}
