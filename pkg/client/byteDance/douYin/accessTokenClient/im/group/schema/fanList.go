package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// Group 群信息
type Group struct {
	ExistNum    int      `json:"exist_num"`
	GroupName   string   `json:"group_name"`
	MaxNum      int      `json:"max_num"`
	Description string   `json:"description"`
	EntryLimit  []string `json:"entry_limit"`
	GroupId     string   `json:"group_id"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
	AvatarUri   string   `json:"avatar_uri"`
}

// DouYinIMGroupFanListReq 查询用户所在群列表请求参数
type DouYinIMGroupFanListRes struct {
	// GroupList 群信息列表
	GroupList []Group `json:"group_list"`
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
