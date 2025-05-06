package response

// DouYinRes 抖音接口通用响应结构
type DouYinRes struct {
	// ErrorCode 错误码，0 表示成功，其他为失败
	ErrorCode int `json:"error_code,omitempty"`

	// Description 错误描述或状态说明
	Description string `json:"description,omitempty"`

	// SubErrorCode 子错误码，用于更详细的错误分类
	SubErrorCode int `json:"sub_error_code,omitempty"`

	// SubDescription 子错误描述，提供更详细的错误信息
	SubDescription string `json:"sub_description,omitempty"`

	// LogId 请求的唯一标识，用于日志追踪
	LogId string `json:"logid,omitempty"`

	// Now 服务器当前时间戳
	Now int64 `json:"now,omitempty"`

	// Cursor 分页游标，用于分页查询
	Cursor int `json:"cursor,omitempty"`

	// HasMore 是否还有更多数据，用于分页查询
	HasMore bool `json:"has_more,omitempty"`
}

// DouYinDataRes 包含通用响应数据的结构
// 用于封装业务数据，通常作为接口响应的最外层结构
type DouYinDataRes struct {
	// Data 包含具体的业务响应信息，通常包含业务数据和通用响应字段
	Data DouYinRes `json:"data,omitempty"`
}
