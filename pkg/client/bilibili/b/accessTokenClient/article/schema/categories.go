package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleCategoriesReq 表示 GET /arcopen/fn/article/categories API 的请求参数
type BiliBiliArticleCategoriesReq struct{}

// BiliBiliArticleCategoriesRes 表示 GET /arcopen/fn/article/categories API 的响应
type BiliBiliArticleCategoriesRes struct {
	response.BiliBiliRes
	// ArticleCategoriesData 文章分类列表
	Data []ArticleCategoriesData `json:"data"`
}

// ArticleCategoriesData 表示文章分类信息
type ArticleCategoriesData struct {
	ID       int                     `json:"id"`        // 分类ID
	ParentID int                     `json:"parent_id"` // 父分类ID
	Name     string                  `json:"name"`      // 分类名称
	Children []ArticleCategoriesData `json:"children"`  // 子分类列表
}
