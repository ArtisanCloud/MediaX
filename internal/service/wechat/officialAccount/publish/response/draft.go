package response

import (
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/core/response"
)

type DraftAddRes struct {
	response.OfficialAccountRes

	MediaID string `json:"media_id"`
}

type NewsItem struct {
	Title              string `json:"title"`
	Author             string `json:"author"`
	Digest             string `json:"digest"`
	Content            string `json:"content"`
	ContentSourceUrl   string `json:"content_source_url"`
	ThumbMediaId       string `json:"thumb_media_id"`
	ThumbUrl           string `json:"thumb_url"`
	ShowCoverPic       int    `json:"show_cover_pic"`
	NeedOpenComment    int    `json:"need_open_comment"`
	OnlyFansCanComment int    `json:"only_fans_can_comment"`
	Url                string `json:"url"`
	IsDeleted          bool   `json:"is_deleted"`
}

type DraftGetRes struct {
	response.OfficialAccountRes

	NewsItem   []*NewsItem `json:"news_item"`
	CreateTime int64       `json:"create_time"`
	UpdateTime int64       `json:"update_time"`
}

type DraftCountRes struct {
	TotalCount int `json:"total_count"`
}

type Content struct {
	NewsItem   []*NewsItem `json:"news_item"`
	CreateTime int64       `json:"create_time"`
	UpdateTime int64       `json:"update_time"`
}

type Item struct {
	MediaId    string   `json:"media_id"`
	ArticleId  string   `json:"article_id"`
	Content    *Content `json:"content"`
	UpdateTime int64    `json:"update_time"`
}

type BatchGetRes struct {
	response.OfficialAccountRes

	TotalCount int     `json:"total_count"`
	ItemCount  int     `json:"item_count"`
	Item       []*Item `json:"item"`
}

type CheckSwitchRes struct {
	response.OfficialAccountRes

	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`
	IsOpen     int `json:"is_open"`
}
