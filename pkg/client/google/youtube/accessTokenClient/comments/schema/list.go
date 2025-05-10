package schema

// YoutubeCommentsListReq 表示 GET /youtube/v3/comments API 的请求参数
// YoutubeCommentsListReq 表示 GET /youtube/v3/comments API 的请求参数
// 该结构体用于定义获取YouTube评论列表时所需的请求参数
type YoutubeCommentsListReq struct {
	Part       string  `json:"part"`                 // 必填，指定API响应包含的comment资源属性，如snippet,id等                 // 必填，指定 API 响应包含的 comment 资源属性
	ID         *string `json:"id,omitempty"`         // 可选，指定要检索的评论ID列表，多个ID用逗号分隔         // 可选，指定要检索的评论 ID 列表
	ParentID   *string `json:"parentId,omitempty"`   // 可选，指定应检索其回复的评论ID   // 可选，指定应检索其回复的评论 ID
	MaxResults *int    `json:"maxResults,omitempty"` // 可选，指定结果集中应返回的最大项数，范围1-100，默认20 // 可选，指定结果集中应返回的最大项数 (1-100, 默认20)
	PageToken  *string `json:"pageToken,omitempty"`  // 可选，用于分页，标识结果集中应返回的特定页面  // 可选，标识结果集中应返回的特定网页
	TextFormat *string `json:"textFormat,omitempty"` // 可选，指定返回评论的格式，可选值：html或plainText // 可选，指定返回评论的格式 (html 或 plainText)
}

// PageInfo 表示分页信息
// 该结构体包含分页查询的统计信息
type PageInfo struct {
	TotalResults   int `json:"totalResults"`   // 总结果数
	ResultsPerPage int `json:"resultsPerPage"` // 每页结果数
}

// YoutubeCommentsListRes 表示 GET /youtube/v3/comments API 的响应
// 该结构体包含YouTube评论列表的响应数据
type YoutubeCommentsListRes struct {
	Kind          string    `json:"kind"`          // 资源类型，固定为"youtube#commentListResponse"
	Etag          string    `json:"etag"`          // 资源的ETag，用于缓存控制
	NextPageToken string    `json:"nextPageToken"` // 下一页的令牌，用于分页
	PageInfo      PageInfo  `json:"pageInfo"`      // 分页信息，包含总结果数和每页结果数
	Items         []Comment `json:"items"`         // 评论列表，包含多个Comment对象
}

// Comment 表示单个YouTube评论的详细信息
// 该结构体包含评论的元数据和内容
type Comment struct {
	Kind    string         `json:"kind"`    // 资源类型，固定为"youtube#comment"
	Etag    string         `json:"etag"`    // 资源的ETag，用于缓存控制
	Id      string         `json:"id"`      // 评论的唯一标识符
	Snippet CommentSnippet `json:"snippet"` // 评论的基本信息，包含作者、内容等
}

type AuthorChannelId struct {
	Value string `json:"value"`
}

// CommentSnippet 表示YouTube评论的详细信息
// 该结构体包含评论的内容、作者信息等元数据
type CommentSnippet struct {
	AuthorDisplayName     string          `json:"authorDisplayName"`     // 评论作者的显示名称
	AuthorProfileImageUrl string          `json:"authorProfileImageUrl"` // 评论作者的头像URL
	AuthorChannelUrl      string          `json:"authorChannelUrl"`      // 评论作者的频道URL
	AuthorChannelId       AuthorChannelId `json:"authorChannelId"`       // 评论作者的频道ID
	ChannelId             string          `json:"channelId"`             // 评论所在频道的ID
	TextDisplay           string          `json:"textDisplay"`           // 评论内容的HTML格式显示
	TextOriginal          string          `json:"textOriginal"`          // 评论内容的原始文本
	ParentId              string          `json:"parentId"`              // 父评论的ID，用于回复评论
	CanRate               bool            `json:"canRate"`               // 当前用户是否可以对评论进行评分
	ViewerRating          string          `json:"viewerRating"`          // 当前用户对评论的评分（like或none）
	LikeCount             int             `json:"likeCount"`             // 评论的点赞数
	ModerationStatus      string          `json:"moderationStatus"`      // 评论的审核状态
	PublishedAt           string          `json:"publishedAt"`           // 评论的发布时间
	UpdatedAt             string          `json:"updatedAt"`             // 评论的最后更新时间
}
