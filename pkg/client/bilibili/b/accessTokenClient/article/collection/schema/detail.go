package schema

// BiliBiliArticleCollectionDetailReq 表示获取文集详情的请求参数
type BiliBiliArticleCollectionDetailReq struct {
	ID int `json:"id"` // 文集ID
}

// BiliBiliArticleCollectionDetailRes 表示获取文集详情的响应参数
type BiliBiliArticleCollectionDetailRes struct {
	Code    int                          `json:"code"`    // 响应状态码
	Message string                       `json:"message"` // 响应消息
	TTL     int                          `json:"ttl"`     // 响应TTL
	Data    *ArticleCollectionDetailData `json:"data"`    // 文集详情数据
}

// ArticleCollectionDetailData 表示文集详情的具体数据
type ArticleCollectionDetailData struct {
	List     *ListDetail      `json:"list"`     // 文集信息
	Articles []*ArticleDetail `json:"articles"` // 文章列表
	Total    int              `json:"total"`    // 文章总数
}

// ListDetail 表示文集的详细信息
type ListDetail struct {
	ID          int    `json:"id"`           // 文集ID
	Name        string `json:"name"`         // 文集名称
	ImageURL    string `json:"image_url"`    // 文集封面URL
	UpdateTime  int    `json:"update_time"`  // 文集更新时间戳
	CTime       int    `json:"ctime"`        // 文集创建时间戳
	PublishTime int    `json:"publish_time"` // 文集发布时间戳
	Summary     string `json:"summary"`      // 文集简介
	Words       int    `json:"words"`        // 文集下文章字数总和
	Read        int    `json:"read"`         // 文集下文章阅读数总和
	State       int    `json:"state"`        // 文集状态
}

// ArticleDetail 表示文章的详细信息
type ArticleDetail struct {
	ID          int    `json:"id"`           // 文章ID
	Title       string `json:"title"`        // 文章标题
	State       int    `json:"state"`        // 文章状态
	PublishTime int    `json:"publish_time"` // 文章发布时间戳
}
