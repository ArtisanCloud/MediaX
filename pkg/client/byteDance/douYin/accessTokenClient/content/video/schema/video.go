package schema

// Statistics 表示视频的统计数据信息。
type Statistics struct {
	// ForwardCount 视频转发次数。
	ForwardCount int `json:"forward_count"`
	// CommentCount 视频评论次数。
	CommentCount int `json:"comment_count"`
	// DiggCount 视频点赞次数。
	DiggCount int `json:"digg_count"`
	// DownloadCount 视频下载次数。
	DownloadCount int `json:"download_count"`
	// PlayCount 视频播放次数。
	PlayCount int `json:"play_count"`
	// ShareCount 视频分享次数。
	ShareCount int `json:"share_count"`
}

// Video 表示视频的详细信息。
type Video struct {
	// Title 视频标题。
	Title string `json:"title"`
	// IsTop 视频是否置顶。
	IsTop bool `json:"is_top"`
	// Cover 视频封面图片 URL。
	Cover string `json:"cover"`
	// MediaType 视频媒体类型。
	MediaType int `json:"media_type"`
	// ItemId 视频唯一标识 ID。
	ItemId string `json:"item_id"`
	// ShareUrl 视频分享链接。
	ShareUrl string `json:"share_url"`
	// VideoStatus 视频状态。
	VideoStatus int `json:"video_status"`
	// IsReviewed 视频是否已审核。
	IsReviewed bool `json:"is_reviewed"`
	// CreateTime 视频创建时间，Unix 时间戳格式。
	CreateTime int `json:"create_time"`
	// Statistics 视频统计数据。
	Statistics Statistics `json:"statistics"`
}
