package schema

// DouYinIMGroupDisableGroupSettingReq 禁用群组配置请求参数
type DouYinIMGroupDisableGroupSettingReq struct {
	GroupId          string `json:"group_id"`           // 粉丝群ID
	GroupSettingType int64  `json:"group_setting_type"` // 群管理配置类型
}

// DouYinIMGroupDisableGroupSettingRes 禁用群组配置响应参数
type DouYinIMGroupDisableGroupSettingRes struct {
	ErrMsg string `json:"err_msg"` // 错误信息
	ErrNo  int    `json:"err_no"`  // 错误码，0-成功，非0-失败
	LogId  string `json:"log_id"`  // 日志ID，用于问题定位
}
