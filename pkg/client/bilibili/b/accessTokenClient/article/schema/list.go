package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleListReq 表示 GET /arcopen/fn/article/list API 的请求参数
type BiliBiliArticleListReq struct {
	PN       *int `json:"pn,omitempty"`       // 页码，默认1
	PS       *int `json:"ps,omitempty"`       // 每页大小，默认10
	Sort     *int `json:"sort,omitempty"`     // 排序类型：1 默认（创建时间倒序）2 喜欢 3 评论 4浏览 5收藏 6硬币
	Group    *int `json:"group,omitempty"`    // 默认0，0-全部（除了草稿和已删除）1-进行中 2-已通过 3-未通过
	Category *int `json:"category,omitempty"` // 分类id
}

// BiliBiliArticleListRes 表示 GET /arcopen/fn/article/list API 的响应
type BiliBiliArticleListRes struct {
	response.BiliBiliRes
	Data *ArticleListData `json:"data"` // 文章列表数据
}

// BiliBiliArticleListData 表示文章列表数据
type ArticleListData struct {
	Articles         []*ArticleData    `json:"articles"`         // 文章列表
	CreationArtsType *CreationArtsType `json:"creationArtsType"` // 文章类型统计
	ArtPage          *ArtPage          `json:"artPage"`          // 分页信息
}

// CreationArtsType 表示文章类型统计
type CreationArtsType struct {
	All       int `json:"all"`       // 所有文章数
	Audit     int `json:"audit"`     // 审核中的文章数
	Passed    int `json:"passed"`    // 通过的文章数
	NotPassed int `json:"notPassed"` // 未通过的文章数
}

// ArtPage 表示分页信息
type ArtPage struct {
	PN    int `json:"pn"`    // 当前页码
	PS    int `json:"ps"`    // 每页大小
	Total int `json:"total"` // 总文章数
}
