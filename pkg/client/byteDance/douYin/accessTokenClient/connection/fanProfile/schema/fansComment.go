package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// Comment 粉丝评论
type Comment struct {
	Keyword  string `json:"keyword"`
	HotValue int    `json:"hot_value"`
}

// DouYinConnectionFanProfileFansCommentRes 获取粉丝评论接口响应结构
type DouYinConnectionFanProfileFansCommentRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// 粉丝评论列表
		List []Comment `json:"list"`
	} `json:"data"`
}
