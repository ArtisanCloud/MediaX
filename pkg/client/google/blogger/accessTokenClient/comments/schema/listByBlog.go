package schema

// BloggerCommentsListByBlogReq 表示 GET /blogs/{blogId}/comments API 的请求参数
type BloggerCommentsListByBlogReq struct {
	BlogId      string  `json:"blogId"`                // 必填，要从中提取评论的博客的ID
	EndDate     *string `json:"endDate,omitempty"`     // 可选，要提取评论的最新日期(RFC 3339格式)
	FetchBodies *bool   `json:"fetchBodies,omitempty"` // 可选，是否包含评论的正文内容
	MaxResults  *uint   `json:"maxResults,omitempty"`  // 可选，返回的最大结果数
	PageToken   *string `json:"pageToken,omitempty"`   // 可选，分页令牌
	StartDate   *string `json:"startDate,omitempty"`   // 可选，要提取评论的最早日期(RFC 3339格式)
}

// BloggerCommentsListByBlogRes 表示 GET /blogs/{blogId}/comments API 的响应
type BloggerCommentsListByBlogRes struct {
	Kind          string            `json:"kind"`          // 资源类型，始终为"blogger#commentList"
	NextPageToken string            `json:"nextPageToken"` // 用于获取下一页的分页令牌
	PrevPageToken string            `json:"prevPageToken"` // 用于获取上一页的分页令牌
	Items         []CommentResource `json:"items"`         // 评论列表
}
