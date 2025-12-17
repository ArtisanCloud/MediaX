package schema

// YoutubeCommentsSetModerationStatusReq 表示设置评论审核状态的请求参数
type YoutubeCommentsSetModerationStatusReq struct {
	// ID 评论ID列表，多个ID用逗号分隔
	ID string `json:"id"`

	// ModerationStatus 审核状态
	// 可选值：heldForReview, published, rejected
	ModerationStatus string `json:"moderationStatus"`

	// BanAuthor 是否禁止评论作者
	// 仅当ModerationStatus为rejected时有效
	BanAuthor *bool `json:"banAuthor,omitempty"`
}

// YoutubeCommentsSetModerationStatusRes 表示设置评论审核状态的响应
// YoutubeCommentsSetModerationStatusRes 表示设置评论审核状态的响应
// 成功时返回HTTP 204 (No Content)状态码
type YoutubeCommentsSetModerationStatusRes struct{}
