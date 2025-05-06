package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMGroupCreateReq 创建群组请求参数
type DouYinIMGroupCreateReq struct {
	AvatarUri       string `json:"avatar_uri"`        // 群头像链接
	Description     string `json:"description"`       // 群简介
	GroupName       string `json:"group_name"`        // 群名称
	ActiveFans      int64  `json:"active_fans"`       // 进群门槛-粉丝活跃度
	AllowInvite     int64  `json:"allow_invite"`      // 允许群成员邀请朋友
	FansLimit       int64  `json:"fans_limit"`        // 进群门槛-粉丝团等级
	GroupType       int64  `json:"group_type"`        // 群类型
	ItemAutoSync    int64  `json:"item_auto_sync"`    // 作品同步
	LiveAutoSync    int64  `json:"live_auto_sync"`    // 直播同步
	OpenAuditSwitch int64  `json:"open_audit_switch"` // 开启进群审批
	RelationType    int64  `json:"relation_type"`     // 进群门槛-关注条件
	ShowAtProfile   int64  `json:"show_at_profile"`   // 展示到个人主页
}

// DouYinIMGroupCreateRes 创建群组响应参数
type DouYinIMGroupCreateRes struct {
	// GroupId 群组ID，创建成功后返回的群组ID。
	GroupId string `json:"group_id"`
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
