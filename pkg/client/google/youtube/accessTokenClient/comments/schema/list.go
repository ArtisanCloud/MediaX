package schema

// YoutubeCommentsListReq 表示 GET /youtube/v3/comments API 的请求参数
type YoutubeCommentsListReq struct {
	Part       string  `json:"part"`                 // 必填，指定 API 响应包含的 comment 资源属性
	ID         *string `json:"id,omitempty"`         // 可选，指定要检索的评论 ID 列表
	ParentID   *string `json:"parentId,omitempty"`   // 可选，指定应检索其回复的评论 ID
	MaxResults *int    `json:"maxResults,omitempty"` // 可选，指定结果集中应返回的最大项数 (1-100, 默认20)
	PageToken  *string `json:"pageToken,omitempty"`  // 可选，标识结果集中应返回的特定网页
	TextFormat *string `json:"textFormat,omitempty"` // 可选，指定返回评论的格式 (html 或 plainText)
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

// YoutubeCommentsListRes 表示 GET /youtube/v3/comments API 的响应
type YoutubeCommentsListRes struct {
	Kind          string    `json:"kind"`
	Etag          string    `json:"etag"`
	NextPageToken string    `json:"nextPageToken"`
	PageInfo      PageInfo  `json:"pageInfo"`
	Items         []Comment `json:"items"`
}

type Comment struct {
	Kind    string         `json:"kind"`
	Etag    string         `json:"etag"`
	Id      string         `json:"id"`
	Snippet CommentSnippet `json:"snippet"`
}

type AuthorChannelId struct {
	Value string `json:"value"`
}

type CommentSnippet struct {
	AuthorDisplayName     string          `json:"authorDisplayName"`
	AuthorProfileImageUrl string          `json:"authorProfileImageUrl"`
	AuthorChannelUrl      string          `json:"authorChannelUrl"`
	AuthorChannelId       AuthorChannelId `json:"authorChannelId"`
	ChannelId             string          `json:"channelId"`
	TextDisplay           string          `json:"textDisplay"`
	TextOriginal          string          `json:"textOriginal"`
	ParentId              string          `json:"parentId"`
	CanRate               bool            `json:"canRate"`
	ViewerRating          string          `json:"viewerRating"`
	LikeCount             int             `json:"likeCount"`
	ModerationStatus      string          `json:"moderationStatus"`
	PublishedAt           string          `json:"publishedAt"`
	UpdatedAt             string          `json:"updatedAt"`
}
