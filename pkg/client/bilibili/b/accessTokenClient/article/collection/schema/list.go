package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleCollectionListRes 表示 GET /arcopen/fn/article/anthology/list API 的响应
type BiliBiliArticleCollectionListRes struct {
	response.BiliBiliRes
	Data ArticleCollectionData `json:"data"` // 响应数据
}

// ArticleCollectionData 包含文集列表和总数
type ArticleCollectionData struct {
	Lists []ArticleCollectionItem `json:"lists"` // 文集列表
	Total int                     `json:"total"` // 文集总数
}

// ArticleCollectionItem 表示单个文集的信息
type ArticleCollectionItem struct {
	ID          int    `json:"id"`           // 文集ID
	Name        string `json:"name"`         // 文集名称
	ImageURL    string `json:"image_url"`    // 文集封面URL
	UpdateTime  int    `json:"update_time"`  // 最后修改时间戳
	CTime       int    `json:"ctime"`        // 文集创建时间戳
	PublishTime int    `json:"publish_time"` // 文集发布时间戳
	Summary     string `json:"summary"`      // 文集简介
	State       int    `json:"state"`        // 文集状态
	ApplyTime   string `json:"apply_time"`   // 文集提交时间
	CheckTime   string `json:"check_time"`   // 文集回查时间
	Total       int    `json:"total"`        // 文集下的文章数
	Words       int    `json:"words"`        // 文集下文章字数总和
	Reason      string `json:"reason"`       // 文集审核不通过原因
}
