package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleCardReq 表示 GET /arcopen/fn/article/cards API 的请求参数
type BiliBiliArticleCardReq struct {
	ResourceID string `json:"resource_id"` // 资源id（视频BV号/文章cv号）：BV开头为视频，cv开头为文章
}

type ArticleCardData struct {
	Snippet string `json:"snippet"` // 卡片片段
}

// BiliBiliArticleCardRes 表示 GET /arcopen/fn/article/cards API 的响应
type BiliBiliArticleCardRes struct {
	response.BiliBiliRes
	Data ArticleCardData `json:"data"`
}
