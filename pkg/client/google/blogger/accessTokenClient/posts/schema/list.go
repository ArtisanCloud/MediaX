package schema

// BloggerPostsListReq 表示 GET /blogger/v3/blogs/{blogId}/posts API 的请求参数
type BloggerPostsListReq struct {
	BlogId      string `json:"blogId"`      // 必填，要从中提取帖子的博客的ID
	EndDate     string `json:"endDate"`     // 可选，要提取的最新发布日期(RFC 3339格式)
	FetchBodies *bool  `json:"fetchBodies"` // 可选，是否包含帖子的正文内容(默认true)
	FetchImages *bool  `json:"fetchImages"` // 可选，是否包含每个帖子的图片网址元数据
	Labels      string `json:"labels"`      // 可选，要搜索的标签的逗号分隔列表
	MaxResults  *int   `json:"maxResults"`  // 可选，要提取的帖子数量上限
	OrderBy     string `json:"orderBy"`     // 可选，排序顺序("published"或"updated")
	SortOption  string `json:"sortOption"`  // 可选，排序方向("descending"或"ascending")
	PageToken   string `json:"pageToken"`   // 可选，分页令牌
	StartDate   string `json:"startDate"`   // 可选，要提取的最早发布日期(RFC 3339格式)
	Status      string `json:"status"`      // 可选，帖子状态("draft"、"live"或"scheduled")
	View        string `json:"view"`        // 可选，查看级别("ADMIN"、"AUTHOR"或"READER")
}

// PostImage 表示博客文章中的图片
type PostImage struct {
	URL string `json:"url"` // 图片的URL
}

// PostAuthorImage 表示博客作者的头像
type PostAuthorImage struct {
	URL string `json:"url"` // 作者头像网址
}

// PostAuthor 表示博客作者信息
type PostAuthor struct {
	ID          string          `json:"id"`          // 作者ID
	DisplayName string          `json:"displayName"` // 作者显示名称
	URL         string          `json:"url"`         // 作者个人资料页面网址
	Image       PostAuthorImage `json:"image"`       // 作者头像信息
}

// PostReplies 表示博客评论信息
type PostReplies struct {
	TotalItems int64         `json:"totalItems"` // 评论总数
	SelfLink   string        `json:"selfLink"`   // 用于检索评论的API网址
	Items      []interface{} `json:"items"`      // 评论列表
}

// PostLocation 表示博客地理位置信息
type PostLocation struct {
	Name string  `json:"name"` // 地点名称
	Lat  float64 `json:"lat"`  // 纬度
	Lng  float64 `json:"lng"`  // 经度
	Span string  `json:"span"` // 视口跨度
}

// BlogInfo 表示博客基本信息
type BlogInfo struct {
	ID string `json:"id"` // 包含此帖子的博客的ID
}

// Post 表示博客文章资源
type PostResource struct {
	Kind      string       `json:"kind"`      // 资源类型，固定值：blogger#post
	ID        string       `json:"id"`        // 帖子的ID
	Blog      BlogInfo     `json:"blog"`      // 包含此帖子的博客的相关数据
	Published string       `json:"published"` // RFC 3339格式的发布时间
	Updated   string       `json:"updated"`   // RFC 3339格式的最后更新时间
	URL       string       `json:"url"`       // 显示此帖子的网址
	SelfLink  string       `json:"selfLink"`  // 用于提取此资源的Blogger API网址
	Title     string       `json:"title"`     // 帖子的标题
	TitleLink string       `json:"titleLink"` // 标题链接网址
	Content   string       `json:"content"`   // 帖子的内容(可包含HTML标记)
	Images    []PostImage  `json:"images"`    // 帖子中的图片列表
	Author    PostAuthor   `json:"author"`    // 帖子作者信息
	Replies   PostReplies  `json:"replies"`   // 帖子评论信息
	Labels    []string     `json:"labels"`    // 帖子标签列表
	Location  PostLocation `json:"location"`  // 地理位置信息
	Status    string       `json:"status"`    // 帖子状态(仅管理员级请求)
}

// BloggerPostsListRes 表示 GET /blogger/v3/blogs/{blogId}/posts API 的响应
type BloggerPostsListRes struct {
	Kind          string         `json:"kind"`          // 资源类型，固定值：blogger#postList
	NextPageToken string         `json:"nextPageToken"` // 用于提取下一页的分页令牌
	Items         []PostResource `json:"items"`         // 此博客的博文列表
}
