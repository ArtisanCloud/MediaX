package schema

type YouTubeSubscriptionsInsertReq struct {
	Part    string  `json:"part"` // 指定返回的资源部分（必填，如 snippet,contentDetails 等）
	Snippet Snippet `json:"snippet"`
}
type YouTubeSubscriptionsInsertRes struct {
	Subscription
}
