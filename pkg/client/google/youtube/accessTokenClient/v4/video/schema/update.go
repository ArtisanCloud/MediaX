package schema

// YouTubeVideoUpdateReq 定义更新 YouTube 视频元数据的请求结构
type YouTubeVideoUpdateReq struct {
	Part                   string                     `json:"part"`                             // 标识写入操作将设置的属性以及 API 响应将包含的属性
	OnBehalfOfContentOwner string                     `json:"onBehalfOfContentOwner,omitempty"` // 代表的内容所有者
	Body                   YouTubeVideoUpdateResource `json:"body"`                             // 视频资源的请求正文
}

// YouTubeVideoUpdateResource 定义更新视频资源的结构
type YouTubeVideoUpdateResource struct {
	Kind             string                  `json:"kind,omitempty"`             // 资源类型，固定为 "youtube#video"
	Etag             string                  `json:"etag,omitempty"`             // 资源的 ETag，用于检查资源是否已更改
	ID               string                  `json:"id,omitempty"`               // 视频的唯一标识符，必须指定
	Snippet          *Snippet                `json:"snippet,omitempty"`          // 视频片段信息，更新 snippet 时部分字段必填
	Localizations    map[string]Localization `json:"localizations,omitempty"`    // 视频本地化信息
	Status           *Status                 `json:"status,omitempty"`           // 视频状态信息
	RecordingDetails *RecordingDetails       `json:"recordingDetails,omitempty"` // 视频录制详情
}

type YouTubeVideoUpdateRes struct {
	Video
}
