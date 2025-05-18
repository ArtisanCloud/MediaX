package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteDeleteDirectLinkReq 表示删除直达链接的请求参数
type JuGuangNoteDeleteDirectLinkReq struct {
	AdvertiserID *int64 `json:"advertiser_id"` // 广告主id
	ID           *int64 `json:"id"`            // 直达链接id
}

// JuGuangNoteDeleteDirectLinkRes 表示删除直达链接的响应
type JuGuangNoteDeleteDirectLinkRes struct {
	response.JuGuangRes
	RequestID string `json:"request_id"` // 请求id
}
