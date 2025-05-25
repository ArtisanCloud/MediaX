package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleCollectionEditReq 表示编辑文集的请求参数
type BiliBiliArticleCollectionEditReq struct {
	ListID   int    `json:"list_id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url,omitempty"`
	Summary  string `json:"summary,omitempty"`
}

// ArticleCollectionEditData 表示编辑文集的响应数据
type ArticleCollectionEditData struct {
	ID          int    `json:"id"`           // 文集ID
	Name        string `json:"name"`         // 文集名称
	Summary     string `json:"summary"`      // 文集简介
	UpdateTime  int64  `json:"update_time"`  // 更新时间
	Ctime       int64  `json:"ctime"`        // 创建时间
	PublishTime int64  `json:"publish_time"` // 发布时间
	Words       int    `json:"words"`        // 字数
	State       int    `json:"state"`        // 状态
}

// BiliBiliArticleCollectionEditRes 表示编辑文集的响应结果
type BiliBiliArticleCollectionEditRes struct {
	response.BiliBiliRes
	Data ArticleCollectionEditData `json:"data"`
}
