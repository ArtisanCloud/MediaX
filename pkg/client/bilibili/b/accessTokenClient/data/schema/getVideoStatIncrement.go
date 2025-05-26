package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type DataGetVideoStatIncrementData struct {
	IncClick int `json:"inc_click"` // 播放数增量
	IncDm    int `json:"inc_dm"`    // 弹幕数增量
	IncReply int `json:"inc_reply"` // 评论数增量
	IncFav   int `json:"inc_fav"`   // 收藏数增量
	IncCoin  int `json:"inc_coin"`  // 投币数增量
	IncShare int `json:"inc_share"` // 分享数增量
	IncLike  int `json:"inc_like"`  // 点赞数增量
	IncElec  int `json:"inc_elec"`  // 充电数增量
}

// BiliBiliDataGetVideoStatIncrementRes 表示 GET /arcopen/fn/data/arc/inc-stats API 的响应
// 接口文档参考：https://member.bilibili.com/arcopen/fn/data/arc/inc-stats
type BiliBiliDataGetVideoStatIncrementRes struct {
	response.BiliBiliRes
	Data DataGetVideoStatIncrementData `json:"data"`
}
