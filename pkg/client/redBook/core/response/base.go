package response

// BaseRes 基础响应
type BaseRes struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// JuGuangRes 基础响应
type JuGuangRes struct {
	BaseRes
	Success   bool   `json:"success,omitempty"`
	RequestId string `json:"request_id,omitempty"`
}
