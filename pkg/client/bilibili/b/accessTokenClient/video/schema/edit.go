package schema

// BiliBiliVideoEditReq 表示编辑视频稿件的请求参数
type BiliBiliVideoEditReq struct {
	ResourceID string `json:"resource_id"`          // 稿件唯一ID
	Title      string `json:"title"`                // 稿件标题
	Cover      string `json:"cover,omitempty"`      // 封面地址
	Tid        int    `json:"tid"`                  // 分区ID
	NoReprint  int    `json:"no_reprint,omitempty"` // 是否禁止转载
	Desc       string `json:"desc,omitempty"`       // 视频描述
}

type EditData struct {
	ResourceID string `json:"resource_id"` // 稿件唯一ID
}

// BiliBiliVideoEditRes 表示编辑视频稿件的响应
type BiliBiliVideoEditRes struct {
	Code      int      `json:"code"`       // 返回码
	Message   string   `json:"message"`    // 返回消息
	RequestID string   `json:"request_id"` // 请求ID
	TTL       int      `json:"ttl"`        // 有效期
	Data      EditData `json:"data"`
}
