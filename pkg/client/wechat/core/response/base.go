package response

type BaseRes struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type OfficialAccountRes struct {
	BaseRes

	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`

	ResultCode string `json:"resultcode,omitempty"`
	ResultMsg  string `json:"resultmsg,omitempty"`
}
