package schema

// YouTubeVideoReportAbuseReq 举报视频请求参数
type YouTubeVideoReportAbuseReq struct {
	VideoId           string `json:"videoId"`           // 视频ID（必填）
	ReasonId          string `json:"reasonId"`          // 举报原因ID（必填）
	SecondaryReasonId string `json:"secondaryReasonId"` // 次要举报原因ID（可选）
	Comments          string `json:"comments"`          // 举报说明（可选）
	Language          string `json:"language"`          // 语言代码（可选）
}

// SecondaryReason 次要举报原因
type SecondaryReason struct {
	Id    string `json:"id"`    // 次要举报原因ID
	Label string `json:"label"` // 次要举报原因标签
}

// videoAbuseReportReasonSnippet 举报原因片段信息
type videoAbuseReportReasonSnippet struct {
	Label            string            `json:"label"`            // 举报原因标签
	SecondaryReasons []SecondaryReason `json:"secondaryReasons"` // 次要举报原因列表
}

// VideoAbuseReportReason 举报原因
type VideoAbuseReportReason struct {
	Kind    string                        `json:"kind"`    // 资源类型
	Etag    string                        `json:"etag"`    // 资源的 ETag
	Id      string                        `json:"id"`      // 举报原因ID
	Snippet videoAbuseReportReasonSnippet `json:"snippet"` // 举报原因片段信息
}

// YouTubeVideoReportAbuseRes 举报视频返回结果
type YouTubeVideoReportAbuseRes struct {
	Kind  string                   `json:"kind"`  // 资源类型
	Etag  string                   `json:"etag"`  // 资源的 ETag
	Items []VideoAbuseReportReason `json:"items"` // 举报原因列表
}
