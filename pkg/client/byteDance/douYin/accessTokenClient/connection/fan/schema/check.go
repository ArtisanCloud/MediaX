package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinConnectionFanCheckRes 抖音粉丝关系检查响应。
type DouYinConnectionFanCheckRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// IsFollower 是否为粉丝，true表示是粉丝，false表示不是粉丝。
		IsFollower bool `json:"is_follower"`
		// FollowTime 粉丝关系建立时间，单位：秒。
		FollowTime int `json:"follow_time"`
	} `json:"data"`
}
