package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
)

type DouYinContentVideoDataReq struct {
}

type Data struct {
	Title       string     `json:"title"`
	CreateTime  int        `json:"create_time"`
	VideoStatus int        `json:"video_status"`
	ShareUrl    string     `json:"share_url"`
	Cover       string     `json:"cover"`
	IsTop       bool       `json:"is_top"`
	Statistics  Statistics `json:"statistics"`
	ItemId      string     `json:"item_id"`
	IsReviewed  bool       `json:"is_reviewed"`
	MediaType   int        `json:"media_type"`
}

type DouYinContentVideoDataRes struct {
	Extra response.DouYinRes `json:"extra,omitempty"`
	Data  struct {
		response.DouYinRes
		List []Data `json:"list"`
	}
}
