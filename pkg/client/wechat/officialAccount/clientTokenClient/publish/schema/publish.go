package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"

// PublishSubmitRes 提交发布响应结构体
type PublishSubmitRes struct {
	response.OfficialAccountRes

	PublishId uint64 `json:"publish_id"` // 发布ID
}

// ArticleItem 文章项结构体
type ArticleItem struct {
	Idx        int    `json:"idx"`         // 文章索引
	ArticleUrl string `json:"article_url"` // 文章URL
}

// ArticleDetail 文章详情结构体
type ArticleDetail struct {
	Count int            `json:"count"` // 文章数量
	Item  []*ArticleItem `json:"item"`  // 文章项列表
}

// PublishGetRes 获取发布状态响应结构体
type PublishGetRes struct {
	response.OfficialAccountRes

	PublishId     uint64         `json:"publish_id"`     // 发布ID
	PublishStatus int            `json:"publish_status"` // 发布状态
	ArticleId     string         `json:"article_id"`     // 文章ID
	ArticleDetail *ArticleDetail `json:"article_detail"` // 文章详情
	FailIdx       []int          `json:"fail_idx"`       // 失败的文章索引列表
}

// PublishGetArticleRes 获取已发布文章响应结构体
type PublishGetArticleRes struct {
	response.OfficialAccountRes

	NewsItem []*NewsItem `json:"news_item"` // 图文消息列表
}
