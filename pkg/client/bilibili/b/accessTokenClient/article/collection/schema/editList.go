package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleCollectionEditListReq 表示编辑文集下文章列表的请求参数
type BiliBiliArticleCollectionEditListReq struct {
	ListID     int    `json:"list_id"`     // 文集ID
	ArticleIDs string `json:"article_ids"` // 文章ID列表，多个ID用逗号分隔
}

// BiliBiliArticleCollectionEditListRes 表示编辑文集下文章列表的响应结果
type BiliBiliArticleCollectionEditListRes struct {
	response.BiliBiliRes
}
