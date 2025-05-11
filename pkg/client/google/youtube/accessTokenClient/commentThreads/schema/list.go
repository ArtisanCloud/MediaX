package schema

// YoutubeCommentThreadsListReq 表示 GET /youtube/v3/commentThreads API 的请求参数
type YoutubeCommentThreadsListReq struct {
	Part                         string  `json:"part"`                                   // 必填，指定 API 响应包含的 commentThread 资源属性
	AllThreadsRelatedToChannelId *string `json:"allThreadsRelatedToChannelId,omitempty"` // 可选，返回与指定频道关联的所有评论会话
	Id                           *string `json:"id,omitempty"`                           // 可选，指定评论会话 ID 列表（以英文逗号分隔）
	VideoId                      *string `json:"videoId,omitempty"`                      // 可选，返回与指定视频 ID 相关联的评论会话
	MaxResults                   *int    `json:"maxResults,omitempty"`                   // 可选，指定结果集中应返回的最大项数 (1-100, 默认20)
	ModerationStatus             *string `json:"moderationStatus,omitempty"`             // 可选，将返回的评论会话串限制为特定的审核状态
	Order                        *string `json:"order,omitempty"`                        // 可选，指定 API 响应列出评论线程的顺序
	PageToken                    *string `json:"pageToken,omitempty"`                    // 可选，标识结果集中应返回的特定网页
	SearchTerms                  *string `json:"searchTerms,omitempty"`                  // 可选，将 API 响应限制为仅包含指定搜索字词的评论
	TextFormat                   *string `json:"textFormat,omitempty"`                   // 可选，指定返回评论的格式 (html 或 plainText)
}

// YoutubeCommentThreadsListRes 表示 GET /youtube/v3/commentThreads API 的响应
type CommentThread struct {
	Kind    string               `json:"kind"`    // 资源类型
	Etag    string               `json:"etag"`    // 资源的 ETag
	Id      string               `json:"id"`      // 评论会话ID
	Snippet CommentThreadSnippet `json:"snippet"` // 评论会话基本信息
}

type CommentThreadSnippet struct {
	ChannelId       string  `json:"channelId"`       // 频道ID
	VideoId         string  `json:"videoId"`         // 视频ID
	TopLevelComment Comment `json:"topLevelComment"` // 顶级评论
	CanReply        bool    `json:"canReply"`        // 是否可以回复
	TotalReplyCount int     `json:"totalReplyCount"` // 总回复数
	IsPublic        bool    `json:"isPublic"`        // 是否公开
}

type Comment struct {
	Kind    string         `json:"kind"`    // 资源类型
	Etag    string         `json:"etag"`    // 资源的 ETag
	Id      string         `json:"id"`      // 评论ID
	Snippet CommentSnippet `json:"snippet"` // 评论基本信息
}

type CommentSnippet struct {
	AuthorDisplayName     string `json:"authorDisplayName"`     // 作者显示名称
	AuthorProfileImageUrl string `json:"authorProfileImageUrl"` // 作者头像URL
	AuthorChannelUrl      string `json:"authorChannelUrl"`      // 作者频道URL
	ChannelId             string `json:"channelId"`             // 频道ID
	VideoId               string `json:"videoId"`               // 视频ID
	TextDisplay           string `json:"textDisplay"`           // 评论内容
	TextOriginal          string `json:"textOriginal"`          // 评论原文
	ParentId              string `json:"parentId"`              // 父评论ID
	CanRate               bool   `json:"canRate"`               // 是否可以评分
	ViewerRating          string `json:"viewerRating"`          // 查看者评分
	LikeCount             int    `json:"likeCount"`             // 点赞数
	ModerationStatus      string `json:"moderationStatus"`      // 审核状态
	PublishedAt           string `json:"publishedAt"`           // 发布时间
	UpdatedAt             string `json:"updatedAt"`             // 更新时间
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`   // 总结果数
	ResultsPerPage int `json:"resultsPerPage"` // 每页结果数
}

type YoutubeCommentThreadsListRes struct {
	Kind          string          `json:"kind"`          // 资源类型
	Etag          string          `json:"etag"`          // 资源的 ETag
	NextPageToken string          `json:"nextPageToken"` // 下一页令牌
	PageInfo      PageInfo        `json:"pageInfo"`      // 分页信息
	Items         []CommentThread `json:"items"`         // 评论会话列表
}
