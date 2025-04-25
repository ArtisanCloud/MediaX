package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/core/response"

type DouYinContentVideoListReq struct {
	Count  string `json:"count" validate:"required,max=20"` // 每页数量，必填，每次查询小于等于20
	OpenID string `json:"open_id" validate:"required"`      // 通过/oauth/access_token/获取，用户唯一标志
	Cursor string `json:"cursor,omitempty"`                 // 分页游标, 第一页请求cursor是0
}

type DouYinContentVideoListRes struct {
	Extra response.DouYinRes `json:"extra,omitempty"`
	Data  struct {
		response.DouYinRes
		HasMore bool    `json:"has_more"`
		List    []Video `json:"list"`
		Cursor  int     `json:"cursor"`
	} `json:"data"`
}
