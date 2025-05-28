package schema

// BloggerBlogsGetReq 表示 GET /blogger/v3/blogs/{blogId} API 的请求参数
type BloggerBlogsGetReq struct {
	BlogID   string `json:"blogId"`   // 必填，要获取的博客的ID
	MaxPosts uint   `json:"maxPosts"` // 可选，随博客一起检索的博文数量上限，不指定则不返回任何博文
}

// Posts 表示博客文章列表信息
type Posts struct {
	TotalItems uint   `json:"totalItems"` // 博文总数
	SelfLink   string `json:"selfLink"`   // 博文集合的API链接
	Items      []Post `json:"items"`      // 博文列表，仅当指定maxPosts时返回
}

// Pages 表示博客页面列表信息
type Pages struct {
	TotalItems uint   `json:"totalItems"` // 页面总数
	SelfLink   string `json:"selfLink"`   // 页面集合的API链接
}

// Locale 表示博客的语言区域设置
type Locale struct {
	Language string `json:"language"` // 博客语言，如"en"表示英语
	Country  string `json:"country"`  // 语言的国家/地区变体，如"US"表示美式英语
	Variant  string `json:"variant"`  // 语言变体
}

// Post 表示博文信息
type Post struct {
	ID        string `json:"id"`        // 博文ID
	Title     string `json:"title"`     // 博文标题
	Content   string `json:"content"`   // 博文内容
	Published string `json:"published"` // 发布时间
	Updated   string `json:"updated"`   // 更新时间
	URL       string `json:"url"`       // 博文URL
	SelfLink  string `json:"selfLink"`  // API资源链接
}

type BloggerBlog struct {
	Kind        string `json:"kind"`        // 资源类型，固定为"blogger#blog"
	ID          string `json:"id"`          // 博客ID
	Name        string `json:"name"`        // 博客名称，可包含HTML
	Description string `json:"description"` // 博客描述，可包含HTML
	Published   string `json:"published"`   // 博客发布时间，RFC 3339格式
	Updated     string `json:"updated"`     // 博客更新时间，RFC 3339格式
	URL         string `json:"url"`         // 博客URL
	SelfLink    string `json:"selfLink"`    // API资源链接
	Posts       Posts  `json:"posts"`       // 博文列表，仅当指定maxPosts时返回
	Pages       Pages  `json:"pages"`       // 页面列表
	Locale      Locale `json:"locale"`      // 博客语言设置
}

// BloggerBlogsGetRes 表示 GET /blogger/v3/blogs/{blogId} API 的响应
type BloggerBlogsGetRes struct {
	BloggerBlog
}
