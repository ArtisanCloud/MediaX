package response

import "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"

type PublishSubmitRes struct {
	response.OfficialAccountRes

	PublishId uint64 `json:"publish_id"`
}

type ArticleItem struct {
	Idx        int    `json:"idx"`
	ArticleUrl string `json:"article_url"`
}

type ArticleDetail struct {
	Count int            `json:"count"`
	Item  []*ArticleItem `json:"item"`
}

type PublishGetRes struct {
	response.OfficialAccountRes

	PublishId     uint64         `json:"publish_id"`
	PublishStatus int            `json:"publish_status"`
	ArticleId     string         `json:"article_id"`
	ArticleDetail *ArticleDetail `json:"article_detail"`
	FailIdx       []int          `json:"fail_idx"`
}

type PublishGetArticleRes struct {
	response.OfficialAccountRes

	NewsItem []*NewsItem `json:"news_item"`
}
