package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNotePostNoteReq 表示聚光平台图文笔记发布接口的请求参数
// 参考：https://open.xiaohongshu.com/document/api?apiId=xxx
// HEADER: Access-Token string token
type JuGuangNotePostNoteReq struct {
	AdvertiserID int64       `json:"advertiser_id"`         // 广告主id，必填
	CoverDatas   []CoverData `json:"cover_datas"`           // 图片信息，必填
	CommonData   *CommonData `json:"common_data,omitempty"` // 基础数据，非必填
	Tags         []string    `json:"tags,omitempty"`        // 话题信息，非必填，不超过45条
}

// CoverData 表示图片信息
type CoverData struct {
	ImageURL   string `json:"image_url"`   // 图片地址，必填
	ImagWidth  string `json:"imag_width"`  // 图片宽度，必填
	ImagHeight string `json:"imag_height"` // 图片高度，必填
}

// CommonData 表示基础数据
type CommonData struct {
	Title string `json:"title,omitempty"` // 图文笔记标题，非必填
	Desc  string `json:"desc,omitempty"`  // 图文笔记描述，非必填
}

type PostNoteData struct {
	TemporaryID string `json:"temporary_id"` // 笔记临时id
}

// JuGuangNotePostNoteRes 表示聚光平台图文笔记发布接口的响应结果
type JuGuangNotePostNoteRes struct {
	response.RedBookAccessTokenRes
	Data PostNoteData `json:"data"`
}
