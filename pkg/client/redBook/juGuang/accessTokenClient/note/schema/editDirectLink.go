package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteEditDirectLinkReq 表示编辑直达链接的请求参数
type JuGuangNoteEditDirectLinkReq struct {
	AdvertiserID int64      `json:"advertiser_id"` // 广告主id
	DirectLink   DirectLink `json:"direct_link"`   // 直达链接信息
}

// JuGuangNoteEditDirectLinkRes 表示编辑直达链接的响应
type JuGuangNoteEditDirectLinkRes struct {
	response.RedBookAccessTokenRes
}
