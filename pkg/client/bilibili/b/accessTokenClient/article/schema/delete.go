package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleDeleteReq 表示 POST /arcopen/fn/article/delete API 的请求参数
type BiliBiliArticleDeleteReq struct {
	ID int `json:"id"` // 文章ID
}

// BiliBiliArticleDeleteRes 表示 POST /arcopen/fn/article/delete API 的响应
type BiliBiliArticleDeleteRes struct {
	response.BiliBiliRes
}
