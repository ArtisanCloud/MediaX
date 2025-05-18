package schema

// YoutubeCaptionsListReq 表示 GET /youtube/v3/captions API 的请求参数
type YoutubeCaptionsListReq struct {
	Part                   string  `json:"part"`                             // 必填，指定 API 响应包含的 caption 资源部分
	VideoId                *string `json:"videoId,omitempty"`                // 可选，指定视频的 YouTube 视频 ID
	Id                     *string `json:"id,omitempty"`                     // 可选，指定以英文逗号分隔的 ID 列表
	OnBehalfOfContentOwner *string `json:"onBehalfOfContentOwner,omitempty"` // 可选，代表内容所有者执行操作
}

// YoutubeCaptionsListRes 表示 GET /youtube/v3/captions API 的响应
type YoutubeCaptionsListRes struct {
	Kind  string    `json:"kind"`  // 资源类型
	Etag  string    `json:"etag"`  // 资源的 ETag
	Items []Caption `json:"items"` // 字幕资源列表
}

// Caption 表示字幕资源
type Caption struct {
	Kind    string  `json:"kind"`    // 资源类型
	Etag    string  `json:"etag"`    // 资源的 ETag
	Id      string  `json:"id"`      // 字幕ID
	Snippet Snippet `json:"snippet"` // 字幕片段信息
}

// Snippet 表示字幕片段信息
type Snippet struct {
	VideoId        string `json:"videoId"`        // 视频ID
	LastUpdated    string `json:"lastUpdated"`    // 最后更新时间
	TrackKind      string `json:"trackKind"`      // 轨道类型
	Language       string `json:"language"`       // 语言代码
	Name           string `json:"name"`           // 字幕名称
	AudioTrackType string `json:"audioTrackType"` // 音轨类型
	IsCC           bool   `json:"isCC"`           // 是否为隐藏式字幕
	IsLarge        bool   `json:"isLarge"`        // 是否为大字幕
	IsEasyReader   bool   `json:"isEasyReader"`   // 是否为易读字幕
	IsDraft        bool   `json:"isDraft"`        // 是否为草稿
	IsAutoSynced   bool   `json:"isAutoSynced"`   // 是否自动同步
	Status         string `json:"status"`         // 字幕状态
	FailureReason  string `json:"failureReason"`  // 失败原因
}
