package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type DataGetUserStatData struct {
	Following      int `json:"following"`        // 关注数
	Follower       int `json:"follower"`         // 粉丝数
	ArcPassedTotal int `json:"arc_passed_total"` // 视频稿件投稿数（审核通过）
}

// BiliBiliDataGetUserStatRes 表示 GET /arcopen/fn/data/user/stat API 的响应
// 接口文档参考：https://member.bilibili.com/arcopen/fn/data/user/stat
type BiliBiliDataGetUserStatRes struct {
	response.BiliBiliRes
	Data DataGetUserStatData `json:"data"`
}
