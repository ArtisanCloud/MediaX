package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMMessageQueryAuthorizeUserListReq 查询授权用户列表请求参数
type DouYinIMMessageQueryAuthorizeUserListReq struct {
	Cursor   string `json:"cursor"`    // 游标，初次查询传 0 即可，后续查询使用 resp 中的 next_cursor 参数
	Limit    int64  `json:"limit"`     // 每次查询数量，不能大于 100
	PageNum  int64  `json:"page_num"`  // 页码，page_size * page_num 需小于 50000
	PageSize int64  `json:"page_size"` // 每页数量，不能大于 50
}

// AuthUser 授权用户信息
type AuthUser struct {
	AuthUserSourceAppId string `json:"auth_user_source_app_id"`
	TargetOpenId        string `json:"target_open_id"`
	DataImExtra         string `json:"data_im_extra"`
	Query               string `json:"query"`
	Path                string `json:"path"`
}

// DouYinIMMessageQueryAuthorizeUserListRes 查询授权用户列表响应参数
type DouYinIMMessageQueryAuthorizeUserListRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// AuthUserList 授权用户列表
		AuthUserList []AuthUser `json:"auth_user_list"`
	} `json:"data"`
}
