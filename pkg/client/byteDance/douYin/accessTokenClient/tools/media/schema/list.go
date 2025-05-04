package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

type DouYinToolMediaGetListRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// Medias 视频列表，包含视频ID、URL列表等信息。
		Medias []Media `json:"medias"`
	} `json:"data"`
}
