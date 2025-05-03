package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// Source 粉丝来源
type Source struct {
	Source  string `json:"source"`
	Percent string `json:"percent"`
}

// DouYinConnectionFanProfileFansSourceRes 获取用户粉丝来源接口响应结构
type DouYinConnectionFanProfileFansSourceRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// List 粉丝来源列表
		List []Source `json:"list"`
	} `json:"data"`
}
