package schema

import "time"

// BloggerPostUserInfosListReq 表示 GET /blogger/v3/users/{userId}/blogs/{blogId}/posts API 的请求参数
type BloggerPostUserInfosListReq struct {
	BlogId      string     `json:"blogId"`                // 必需，博客的ID
	UserId      string     `json:"userId"`                // 必需，用户ID（"self"或用户个人资料标识符）
	EndDate     *time.Time `json:"endDate,omitempty"`     // 可选，要提取的最新帖子日期
	FetchBodies *bool      `json:"fetchBodies,omitempty"` // 可选，是否包含帖子正文
	Labels      *string    `json:"labels,omitempty"`      // 可选，要搜索的标签列表
	MaxResults  *uint      `json:"maxResults,omitempty"`  // 可选，提取的帖子数量上限
	OrderBy     *string    `json:"orderBy,omitempty"`     // 可选，排序方式（published/updated）
	PageToken   *string    `json:"pageToken,omitempty"`   // 可选，分页令牌
	StartDate   *time.Time `json:"startDate,omitempty"`   // 可选，最早发布日期
	Status      *string    `json:"status,omitempty"`      // 可选，帖子状态（draft/live/scheduled）
	View        *string    `json:"view,omitempty"`        // 可选，视图级别（ADMIN/AUTHOR/READER）
}

// BloggerPostUserInfosListRes 表示 GET /blogger/v3/users/{userId}/blogs/{blogId}/posts API 的响应
type BloggerPostUserInfosListRes struct {
	Kind          string                  `json:"kind"`          // 资源类型，固定为blogger#postUserInfosList
	NextPageToken string                  `json:"nextPageToken"` // 下一页令牌
	Items         []PoserUserInfoResource `json:"items"`         // 帖子用户信息列表
}
