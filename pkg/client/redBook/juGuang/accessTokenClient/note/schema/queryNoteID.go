package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteQueryNoteIDReq 请求结构体
// 用于通过临时笔记ID查询正式笔记ID
// temporary_id: 临时笔记ID（必填）
type JuGuangNoteQueryNoteIDReq struct {
	TemporaryID string `json:"temporary_id"` // 临时笔记ID
}

// JuGuangNoteQueryNoteIDRes 响应结构体
// 返回查询结果，包括笔记ID、创建时间等信息
type JuGuangNoteQueryNoteIDRes struct {
	response.RedBookAccessTokenRes
	Data *JuGuangNoteQueryNoteIDResData `json:"data"` // 业务数据
}

// JuGuangNoteQueryNoteIDResData 业务数据结构体
type JuGuangNoteQueryNoteIDResData struct {
	TemporaryID string `json:"temporary_id"` // 临时笔记ID
	NoteID      string `json:"note_id"`      // 笔记ID
	CreateTime  int64  `json:"create_time"`  // 创建时间
}
