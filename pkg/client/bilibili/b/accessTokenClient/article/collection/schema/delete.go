package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleCollectionDeleteReq 表示 DELETE /arcopen/fn/article/anthology/delete API 的请求参数
type BiliBiliArticleCollectionDeleteReq struct {
	ID int `json:"id"` // 文集ID
}

// BiliBiliArticleCollectionDeleteRes 表示 DELETE /arcopen/fn/article/anthology/delete API 的响应
type BiliBiliArticleCollectionDeleteRes struct {
	response.BiliBiliRes
}
