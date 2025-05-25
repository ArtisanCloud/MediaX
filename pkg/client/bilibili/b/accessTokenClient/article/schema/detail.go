package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleDetailReq 表示 GET /arcopen/fn/article/detail API 的请求参数
type BiliBiliArticleDetailReq struct {
	ID int `json:"id"` // 文章ID
}

// BiliBiliArticleDetailRes 表示 GET /arcopen/fn/article/detail API 的响应
type BiliBiliArticleDetailRes struct {
	response.BiliBiliRes
	Data *ArticleData `json:"data"` // 文章数据
}

// BiliBiliArticleData 表示文章详情数据
type ArticleData struct {
	Content      string           `json:"content"`        // 文章内容
	ID           int              `json:"id"`             // 文章ID
	Category     *ArticleCategory `json:"category"`       // 文章分类
	Title        string           `json:"title"`          // 文章标题
	Summary      string           `json:"summary"`        // 文章简介
	BannerURL    string           `json:"banner_url"`     // 文章封面URL
	TemplateID   int              `json:"template_id"`    // 模板ID
	State        int              `json:"state"`          // 文章状态
	ImageURLs    []string         `json:"image_urls"`     // 文章图片URL列表
	PublishTime  int64            `json:"publish_time"`   // 发布时间
	CTime        int64            `json:"ctime"`          // 创建时间
	Stats        *ArticleStats    `json:"stats"`          // 文章统计数据
	Tags         []*ArticleTag    `json:"tags"`           // 文章标签列表
	Reason       string           `json:"reason"`         // 审核原因
	Words        int              `json:"words"`          // 字数
	List         *ArticleList     `json:"list"`           // 文集信息
	Original     int              `json:"original"`       // 是否原创
	TopVideoBVID string           `json:"top_video_bvid"` // 置顶视频BVID
	Type         int              `json:"type"`           // 文章类型
}

// ArticleCategory 表示文章分类
type ArticleCategory struct {
	ID       int    `json:"id"`        // 分类ID
	ParentID int    `json:"parent_id"` // 父分类ID
	Name     string `json:"name"`      // 分类名称
}

// ArticleStats 表示文章统计数据
type ArticleStats struct {
	View     int `json:"view"`     // 浏览量
	Favorite int `json:"favorite"` // 收藏数
	Like     int `json:"like"`     // 点赞数
	Dislike  int `json:"dislike"`  // 点踩数
	Reply    int `json:"reply"`    // 评论数
	Share    int `json:"share"`    // 分享数
	Coin     int `json:"coin"`     // 投币数
}

// ArticleTag 表示文章标签
type ArticleTag struct {
	TID  int    `json:"tid"`  // 标签ID
	Name string `json:"name"` // 标签名称
}

// ArticleList 表示文集信息
type ArticleList struct {
	ID          int    `json:"id"`           // 文集ID
	Name        string `json:"name"`         // 文集名称
	ImageURL    string `json:"image_url"`    // 文集封面URL
	UpdateTime  int64  `json:"update_time"`  // 更新时间
	CTime       int64  `json:"ctime"`        // 创建时间
	PublishTime int64  `json:"publish_time"` // 发布时间
	Summary     string `json:"summary"`      // 文集简介
	Words       int    `json:"words"`        // 字数
}
