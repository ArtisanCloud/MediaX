package response

// BaseRes 基础响应结构体
type BaseRes struct {
	Code    string `json:"code,omitempty"`    // 响应代码
	Message string `json:"message,omitempty"` // 响应消息
}

// OfficialAccountRes 微信公众号响应结构体
type OfficialAccountRes struct {
	BaseRes

	ErrCode int    `json:"errcode,omitempty"` // 错误代码
	ErrMsg  string `json:"errmsg,omitempty"`  // 错误信息

	ResultCode string `json:"resultcode,omitempty"` // 结果代码
	ResultMsg  string `json:"resultmsg,omitempty"`  // 结果信息
}
