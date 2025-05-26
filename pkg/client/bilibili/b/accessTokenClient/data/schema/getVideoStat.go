package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type DataGetVideoStatData struct {
	Title    string `json:"title"`    // 标题
	Ptime    int    `json:"ptime"`    // 发布时间
	View     int    `json:"view"`     // 播放数
	Danmaku  int    `json:"danmaku"`  // 弹幕数
	Reply    int    `json:"reply"`    // 评论数
	Favorite int    `json:"favorite"` // 收藏数
	Coin     int    `json:"coin"`     // 投币数
	Share    int    `json:"share"`    // 分享数
	Like     int    `json:"like"`     // 点赞数
}

// BiliBiliDataGetVideoStatRes 表示 GET /arcopen/fn/data/arc/stat API 的响应
// 接口文档参考：https://member.bilibili.com/arcopen/fn/data/arc/stat
type BiliBiliDataGetVideoStatRes struct {
	response.BiliBiliRes
	Data DataGetVideoStatData `json:"data"`
}
