package schema

// YouTubeVideoInsertReq 定义上传视频到 YouTube 的请求结构
type YouTubeVideoInsertReq struct {
	Part                          string               `json:"part"`                                    // Part 标识写入操作将设置的属性以及 API 响应将包含的属性
	NotifySubscribers             *bool                `json:"notifySubscribers,omitempty"`             // NotifySubscribers 指明 YouTube 是否应向订阅视频频道的用户发送有关新视频的通知
	OnBehalfOfContentOwner        string               `json:"onBehalfOfContentOwner,omitempty"`        // OnBehalfOfContentOwner 代表的内容所有者
	OnBehalfOfContentOwnerChannel string               `json:"onBehalfOfContentOwnerChannel,omitempty"` // OnBehalfOfContentOwnerChannel 要将视频添加到的频道的 YouTube 频道 ID
	Body                          YouTubeVideoResource `json:"body"`                                    // Body 视频资源的请求正文
}

// YouTubeVideoResource 定义视频资源的结构
type YouTubeVideoResource struct {
	Snippet          Snippet                 `json:"snippet"`
	Localizations    map[string]Localization `json:"localizations,omitempty"`
	Status           Status                  `json:"status"`
	RecordingDetails RecordingDetails        `json:"recordingDetails,omitempty"`
}

// Snippet 定义视频的元数据
type YouTubeVideoInsertRes struct {
	// Video 视频资源
	Video
}
