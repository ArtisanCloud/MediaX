package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMGroupGetUserEnterAuditReq 查询用户入群申请状态请求参数
type DouYinIMGroupGetUserEnterAuditReq struct {
	Count int32 `json:"count"` // 每页的数量，最大不超过50，最小不低于1
	//OpenId string `json:"open_id"` // 用户open_id，通过/oauth/access_token/获取，用户唯一标志
	Cursor int64 `json:"cursor"` // 分页游标，第一页请求cursor是0
}

type Apply struct {
	ApplyId     string `json:"apply_id"`
	ApplyStatus int    `json:"apply_status"`
	CreateTime  int64  `json:"create_time"`
	GroupId     string `json:"group_id"`
	UserId      string `json:"user_id"`
}

type ApplyList struct {
	Cursor  int     `json:"cursor"`
	HasMore bool    `json:"has_more"`
	List    []Apply `json:"list"`
}

type DouYinIMGroupGetUserEnterAuditRes struct {
	ApplyList ApplyList `json:"apply_list"`

	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
