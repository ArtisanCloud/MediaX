package schema

import "time"

type BloggerCommentsListReq struct {
	BlogId      string    `json:"blogId"`      // 必需参数，要从中提取评论的博客ID
	PostId      string    `json:"postId"`      // 必需参数，要从中提取评论的帖子ID
	EndDate     time.Time `json:"endDate"`     // 可选参数，要提取评论的最新日期
	FetchBodies bool      `json:"fetchBodies"` // 可选参数，是否包含评论的正文内容
	MaxResults  uint      `json:"maxResults"`  // 可选参数，结果中可包含的评论数量上限
	PageToken   string    `json:"pageToken"`   // 可选参数，如果请求已分页，则为延续令牌
	StartDate   time.Time `json:"startDate"`   // 可选参数，要提取评论的最早日期
	Status      string    `json:"status"`      // 可选参数，评论状态
	View        string    `json:"view"`        // 可选参数，视图级别
}

type InReplyTo struct {
	Id string `json:"id"` // 回复的评论ID
}

type Post struct {
	Id string `json:"id"` // 所属帖子ID
}

type Blog struct {
	Id string `json:"id"` // 所属博客ID
}

type Image struct {
	Url string `json:"url"` // 作者头像URL
}

type Author struct {
	Id          string `json:"id"`          // 作者ID
	DisplayName string `json:"displayName"` // 作者显示名称
	Url         string `json:"url"`         // 作者主页URL
	Image       Image  `json:"image"`
}

type CommentResource struct {
	Kind      string    `json:"kind"`   // 资源类型，始终为blogger#comment
	Status    string    `json:"status"` // 评论状态
	Id        string    `json:"id"`     // 评论ID
	InReplyTo InReplyTo `json:"inReplyTo"`
	Post      Post      `json:"post"`
	Blog      Blog      `json:"blog"`
	Published time.Time `json:"published"` // 发布时间
	Updated   time.Time `json:"updated"`   // 更新时间
	SelfLink  string    `json:"selfLink"`  // 资源链接
	Content   string    `json:"content"`   // 评论内容
	Author    Author    `json:"author"`
}

type BloggerCommentsListRes struct {
	Kind          string            `json:"kind"`          // 此条目的类型，始终为blogger#commentList
	NextPageToken string            `json:"nextPageToken"` // 用于提取下一页的分页令牌
	PrevPageToken string            `json:"prevPageToken"` // 用于获取上一页的分页令牌
	Items         []CommentResource `json:"items"`         // 评论资源列表
}
