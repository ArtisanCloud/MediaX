package schema

type Statistics struct {
	ForwardCount  int `json:"forward_count"`
	CommentCount  int `json:"comment_count"`
	DiggCount     int `json:"digg_count"`
	DownloadCount int `json:"download_count"`
	PlayCount     int `json:"play_count"`
	ShareCount    int `json:"share_count"`
}

type Video struct {
	Title       string     `json:"title"`
	IsTop       bool       `json:"is_top"`
	CreateTime  int        `json:"create_time"`
	IsReviewed  bool       `json:"is_reviewed"`
	VideoStatus int        `json:"video_status"`
	ShareUrl    string     `json:"share_url"`
	ItemId      string     `json:"item_id"`
	MediaType   int        `json:"media_type"`
	Cover       string     `json:"cover"`
	Statistics  Statistics `json:"statistics"`
}
