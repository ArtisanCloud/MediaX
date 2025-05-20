package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNotePostVideoNoteReq 表示聚光平台视频笔记发布接口的请求参数
// 参考：https://open.xiaohongshu.com/document/api?apiId=xxx
// HEADER: Access-Token string token
type JuGuangNotePostVideoNoteReq struct {
	AdvertiserID int64       `json:"advertiser_id"`         // 广告主id，必填
	VideoData    VideoData   `json:"video_data"`            // 视频信息，必填
	CoverData    CoverData   `json:"cover_data"`            // 视频封面信息，必填
	CommonData   *CommonData `json:"common_data,omitempty"` // 基础数据，非必填
	Tags         []string    `json:"tags,omitempty"`        // 话题信息，非必填，不超过45条
}

// VideoData 表示视频信息
// 视频地址、宽度和高度均为必填
// video_url 域名需提前加白
// video_width、video_height 为视频尺寸
type VideoData struct {
	VideoURL    string `json:"video_url"`    // 视频地址，必填
	VideoWidth  string `json:"video_width"`  // 视频宽度，必填
	VideoHeight string `json:"video_height"` // 视频高度，必填
}

// JuGuangNotePostVideoNoteRes 表示聚光平台视频笔记发布接口的响应结果
type JuGuangNotePostVideoNoteRes struct {
	response.RedBookAccessTokenRes
	Data *PostVideoNoteData `json:"data"` // 业务数据
}

// PostVideoNoteData 业务数据结构体
type PostVideoNoteData struct {
	TemporaryID string `json:"temporary_id"` // 笔记临时id
}
