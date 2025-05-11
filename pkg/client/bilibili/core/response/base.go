package response

type BaseRes struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type BiliBiliRes struct {
	BaseRes

	RequestId string `json:"request_id"` // 请求ID
	TTL       int    `json:"ttl"`        // 有效期
}
