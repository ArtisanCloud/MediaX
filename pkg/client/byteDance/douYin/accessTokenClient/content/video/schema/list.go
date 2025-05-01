package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
)

// DouYinContentVideoListReq 表示获取抖音视频列表的请求结构体。
// 包含分页、用户唯一标识等字段。
type DouYinContentVideoListReq struct {
	// Count 每页数量，必填字段，每次查询小于等于20。
	Count string `json:"count" validate:"required,max=20"`
	// OpenID 用户唯一标志，通过 /oauth/access_token/ 获取，必填字段。
	OpenID string `json:"open_id" validate:"required"`
	// Cursor 分页游标，第一页请求 cursor 是 0，可选字段。
	Cursor string `json:"cursor,omitempty"`
}

// DouYinContentVideoListRes 表示获取抖音视频列表的响应结构体。
// 包含通用返回信息、视频列表、分页游标等字段。
type DouYinContentVideoListRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	// Cursor 分页游标，用于获取下一页数据。
	// List 视频列表，包含多个视频信息。
	// HasMore 表示是否有更多数据，true 表示还有更多数据，false 表示没有更多数据。
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Extra response.DouYinRes `json:"extra,omitempty"`
	Data  struct {
		// Data通用返回信息
		response.DouYinRes
		HasMore bool    `json:"has_more"`
		List    []Video `json:"list"`
		Cursor  int     `json:"cursor"`
	} `json:"data"`
}
