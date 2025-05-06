package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/message/schema"

// DouYinIMGroupSetGroupSettingReq 设置群组配置请求参数
type DouYinIMGroupSetGroupSettingReq struct {
	GroupSettingType int64            `json:"group_setting_type"` // 群管理配置类型
	MsgList          []schema.Content `json:"msg_list"`           // 消息列表
}

// DouYinIMGroupSetGroupSettingRes 设置群组配置响应参数
type DouYinIMGroupSetGroupSettingRes struct {
	ErrMsg string `json:"err_msg"` // 错误信息
	ErrNo  int    `json:"err_no"`  // 错误码，0-成功，非0-失败
	LogId  string `json:"log_id"`  // 日志ID，用于问题定位
}
