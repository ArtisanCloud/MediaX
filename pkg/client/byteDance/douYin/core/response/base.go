package response

type DouYinRes struct {
	ErrorCode      int    `json:"error_code,omitempty"`
	Description    string `json:"description,omitempty"`
	SubErrorCode   int    `json:"sub_error_code,omitempty"`
	SubDescription string `json:"sub_description,omitempty"`
	LogId          string `json:"logid,omitempty"`
	Now            int64  `json:"now,omitempty"`
	Cursor         int    `json:"cursor,omitempty"`
	HasMore        bool   `json:"has_more,omitempty"`
}

type DouYinDataRes struct {
	Data DouYinRes `json:"data,omitempty"`
}
