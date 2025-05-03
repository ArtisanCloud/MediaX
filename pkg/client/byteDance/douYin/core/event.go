package core

// 事件类型常量
const (
	DouYinEventTypeMessageIMReceiveMsg          = "im_receive_msg"           // 收到消息
	DouYinEventTypeMessageIMEnterDirectMsg      = "im_enter_direct_msg"      // 进入私信
	DouYinEventTypeMessageIMRecallMsg           = "im_recall_msg"            // 消息撤回
	DouYinEventTypeMessageEnterGroupFanEvent    = "group_fans_event"         // 进群事件
	DouYinEventTypeMessageEnterGroupAuditChange = "enter_group_audit_change" // 进群审核状态变更
)

// 事件消息类型
const (
	DouYinIMMessageTypeText              = "text"                // 文本消息
	DouYinIMMessageTypeImage             = "image"               // 图片消息
	DouYinIMMessageTypeVideo             = "video"               // 视频消息
	DouYinIMMessageTypeUserLocalImage    = "user_local_image"    // 用户本地图片
	DouYinIMMessageTypeUserLocalVideo    = "user_local_video"    // 用户本地视频
	DouYinIMMessageTypeRetainConsultCard = "retain_consult_card" // 保留咨询卡片
	DouYinIMMessageTypeOther             = "other"               // 其他类型
)

// CardData 卡片数据
type CardData struct {
	Label string `json:"label"` // 标签
	Value string `json:"value"` // 值
}

// UserInfo 用户信息
type UserInfo struct {
	OpenId   string `json:"open_id"`   // 用户OpenID
	NickName string `json:"nick_name"` // 用户昵称
	Avatar   string `json:"avatar"`    // 用户头像
}

type EventContent struct {
	ConversationShortId string     `json:"conversation_short_id"` // 会话短ID
	ServerMessageId     string     `json:"server_message_id"`     // 服务端消息ID
	ConversationType    int        `json:"conversation_type"`     // 会话类型
	CreateTime          int64      `json:"create_time"`           // 创建时间
	MessageType         string     `json:"message_type"`          // 消息类型
	Source              string     `json:"source"`                // 消息来源
	UserInfos           []UserInfo `json:"user_infos"`            // 用户信息列表

	// text 文本消息
	Text string `json:"text"` // 文本内容

	// emoji 表情消息
	ResourceType   string `json:"resource_type"`   // 资源类型
	ResourceWidth  int    `json:"resource_width"`  // 资源宽度
	ResourceHeight int    `json:"resource_height"` // 资源高度
	ResourceUrl    string `json:"resource_url"`    // 资源URL

	// video 视频消息
	ItemId string `json:"item_id"` // 视频ID

	// card 卡片消息
	CardId     string     `json:"card_id"`     // 卡片ID
	CardStatus int        `json:"card_status"` // 卡片状态
	CardData   []CardData `json:"card_data"`   // 卡片数据
}

// DouYinEventMessage 抖音消息事件
type DouYinEventMessage struct {
	Event      string       `json:"event"`        // 事件类型
	ClientKey  string       `json:"client_key"`   // 客户端Key
	FromUserId string       `json:"from_user_id"` // 发送用户ID
	ToUserId   string       `json:"to_user_id"`   // 接收用户ID
	LogId      string       `json:"log_id"`       // 日志ID
	Content    EventContent `json:"content"`      // 消息内容
}
