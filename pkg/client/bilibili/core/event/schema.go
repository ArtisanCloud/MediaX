package event

const (
	BiliBiliAuthorizeEventName   = "authorize"   // 授权事件
	BiliBiliDeauthorizeEventName = "deauthorize" // 解除授权事件
	BiliBiliVideoOpenEventName   = "video_open"  // 视频稿件开放浏览事件
	BiliBiliVideoFailEventName   = "video_fail"  // 视频投稿失败事件
)

// AuthorizeEvent 用户授权事件
type BaseEvent struct {
	Event     string `json:"event"`     // 事件类型
	Timestamp string `json:"timestamp"` // 时间戳
}

// AuthorizeEvent 授权内容
type AuthorizeEvent struct {
	BaseEvent
	Content struct {
		Openid   string `json:"openid"`    // 用户OpenID
		ClientID string `json:"client_id"` // 客户端ID
		Permits  string `json:"permits"`   // 授权权限
	} `json:"content"`
}

// DeauthorizeEvent 解除授权内容
type DeauthorizeEvent struct {
	BaseEvent
	Content struct {
		Openid   string `json:"openid"`    // 用户OpenID
		ClientID string `json:"client_id"` // 客户端ID
		Permits  string `json:"permits"`   // 授权权限
	} `json:"content"`
}

// VideoOpenEvent 视频稿件开放浏览事件
type VideoOpenEvent struct {
	BaseEvent
	Content struct {
		Openid     string `json:"openid"`      // 用户OpenID
		ClientID   string `json:"client_id"`   // 客户端ID
		ResourceID string `json:"resource_id"` // 视频资源ID
		State      int    `json:"state"`       // 状态
		StateDesc  string `json:"state_desc"`  // 状态描述
	} `json:"content"`
}

// VideoFailEvent 视频投稿失败事件
type VideoFailEvent struct {
	BaseEvent
	Content struct {
		Openid     string `json:"openid"`      // 用户OpenID
		ClientID   string `json:"client_id"`   // 客户端ID
		ResourceID string `json:"resource_id"` // 视频资源ID
		State      int    `json:"state"`       // 状态
		StateDesc  string `json:"state_desc"`  // 状态描述
	} `json:"content"`
}
