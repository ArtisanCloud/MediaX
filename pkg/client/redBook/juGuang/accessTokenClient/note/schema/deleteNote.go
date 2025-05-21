package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteDeleteNoteReq 表示删除笔记接口的请求参数
type JuGuangNoteDeleteNoteReq struct {
	AdvertiserID int64  `json:"advertiser_id"` // 广告主ID，必填
	NoteID       string `json:"note_id"`       // 笔记ID，必填
}

// JuGuangNoteDeleteNoteRes 表示删除笔记接口的响应结果
type JuGuangNoteDeleteNoteRes struct {
	response.RedBookAccessTokenRes
}
