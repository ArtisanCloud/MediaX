package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v4/video/schema"

// YouTubeVideoAbuseReportReasonsListReq 获取视频举报原因列表请求参数
type YouTubeVideoAbuseReportReasonsListReq struct {
	Part string `json:"part"`         // 指定返回的资源部分（必填，如 id,snippet）
	HL   string `json:"hl,omitempty"` // 语言代码（可选，默认 en_US）
}

type YouTubeVideoAbuseReportReasonsListRes struct {
	Kind  string        `json:"kind"`
	Etag  string        `json:"etag"`
	Items []schema.Item `json:"items"`
}
