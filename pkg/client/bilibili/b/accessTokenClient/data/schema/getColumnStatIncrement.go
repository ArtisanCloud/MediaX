package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type GetColumnStatIncrementData struct {
	IncReply int `json:"inc_reply"` // 评论数增量
	IncRead  int `json:"inc_read"`  // 阅读数增量
	IncFav   int `json:"inc_fav"`   // 收藏数增量
	IncLikes int `json:"inc_likes"` // 点赞数增量
	IncShare int `json:"inc_share"` // 分享数增量
	IncCoin  int `json:"inc_coin"`  // 投币数增量
}

// BiliBiliDataGetColumnStatIncrementRes 获取专栏增量数据响应
type BiliBiliDataGetColumnStatIncrementRes struct {
	response.BiliBiliRes
	Data GetColumnStatIncrementData `json:"data"`
}
