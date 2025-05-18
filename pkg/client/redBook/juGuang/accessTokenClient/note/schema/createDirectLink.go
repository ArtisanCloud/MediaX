package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteCreateDirectLinkReq 表示创建直达链接的请求参数
type JuGuangNoteCreateDirectLinkReq struct {
	AdvertiserID   int64        `json:"advertiser_id"`    // 广告主id
	DirectLinkList []DirectLink `json:"direct_link_list"` // 直达链接信息列表
}

// JuGuangNoteCreateDirectLinkRes 表示创建直达链接的响应
type JuGuangNoteCreateDirectLinkRes struct {
	response.JuGuangRes
	Data      []DirectLink `json:"data"`       // 返回的直达链接信息
	RequestID string       `json:"request_id"` // 请求id
}
