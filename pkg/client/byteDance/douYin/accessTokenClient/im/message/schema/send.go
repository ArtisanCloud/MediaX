package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// 消息类型常量
const (
	DouYinMsgTypeText                   = 1   // 文本消息
	DouYinMsgTypeImage                  = 2   // 图片消息
	DouYinMsgTypeVideo                  = 3   // 视频消息
	DouYinMsgTypeRetainConsultCard      = 8   // 保留咨询卡片
	DouYinMsgTypeGroupInvitation        = 9   // 群邀请
	DouYinMsgTypeAppletCard             = 10  // 小程序卡片
	DouYinMsgTypeAppletCoupon           = 11  // 小程序优惠券
	DouYinMsgTypeAuthPrivateMessageCard = 12  // 授权私信卡片
	DouYinMsgTypeQuestionGuideMsgCard   = 204 // 问题引导消息卡片
	DouYinMsgTypeQuestionGuideJumpCard  = 205 // 问题引导跳转卡片
)

// Text 文本消息内容
type Text struct {
	Text string `json:"text"` // 文本内容
}

// Image 图片消息内容
type Image struct {
	MediaId string `json:"media_id"` // 媒体文件ID
}

// Video 视频消息内容
type Video struct {
	ItemId string `json:"item_id"` // 视频ID
}

// RetainConsultCard 保留咨询卡片
type RetainConsultCard struct {
	CardId string `json:"card_id"` // 卡片ID
}

// GroupInvitation 群邀请
type GroupInvitation struct {
	GroupId    string `json:"group_id"`    // 群ID
	GroupToken string `json:"group_token"` // 群Token
}

// AppletCard 小程序卡片
type AppletCard struct {
	CardTemplateId string `json:"card_template_id"` // 卡片模板ID
	Path           string `json:"path"`             // 小程序路径
	Query          string `json:"query"`            // 小程序查询参数
	AppId          string `json:"app_id"`           // 小程序ID
}

// AppletCoupon 小程序优惠券
type AppletCoupon struct {
	ActivityId   int `json:"activity_id"`    // 活动ID
	CouponMetaId int `json:"coupon_meta_id"` // 优惠券元数据ID
}

// AppInfo 应用信息
type AppInfo struct {
	AppId string `json:"app_id"` // 应用ID
}

// ToUserInfo 接收用户信息
type ToUserInfo struct {
	OpenId string `json:"open_id"` // 用户OpenID
	AppId  string `json:"app_id"`  // 应用ID
}

// AuthPrivateMessageCard 授权私信卡片
type AuthPrivateMessageCard struct {
	AppInfo    AppInfo    `json:"app_info"`     // 应用信息
	ToUserInfo ToUserInfo `json:"to_user_info"` // 接收用户信息
}

// SchemaData 跳转数据
type SchemaData struct {
	AppId string `json:"app_id"` // 应用ID
	Path  string `json:"path"`   // 跳转路径
	Query string `json:"query"`  // 跳转查询参数
}

// Question 问题
type Question struct {
	Text       string     `json:"text"`        // 问题文本
	JumpType   int        `json:"jump_type"`   // 跳转类型
	SchemaData SchemaData `json:"schema_data"` // 跳转数据
}

// QuestionGuideMsgCard 问题引导消息卡片
type QuestionGuideMsgCard struct {
	Title        string     `json:"title"`         // 标题
	QuestionList []Question `json:"question_list"` // 问题列表
}

// QuestionGuideJumpCard 问题引导跳转卡片
type QuestionGuideJumpCard struct {
	Title        string     `json:"title"`         // 标题
	QuestionList []Question `json:"question_list"` // 问题列表
}

// Content 消息内容
type Content struct {
	MsgType                int                    `json:"msg_type"`                  // 消息类型
	Text                   Text                   `json:"text"`                      // 文本消息
	Image                  Image                  `json:"image"`                     // 图片消息
	Video                  Video                  `json:"video"`                     // 视频消息
	RetainConsultCard      RetainConsultCard      `json:"retain_consult_card"`       // 保留咨询卡片
	GroupInvitation        GroupInvitation        `json:"group_invitation"`          // 群邀请
	AppletCard             AppletCard             `json:"applet_card"`               // 小程序卡片
	AppletCoupon           AppletCoupon           `json:"applet_coupon"`             // 小程序优惠券
	AuthPrivateMessageCard AuthPrivateMessageCard `json:"auth_private_message_card"` // 授权私信卡片
	QuestionGuideMsgCard   QuestionGuideMsgCard   `json:"question_guide_msg_card"`   // 问题引导消息卡片
	QuestionGuideJumpCard  QuestionGuideJumpCard  `json:"question_guide_jump_card"`  // 问题引导跳转卡片
}

// DouYinIMMessageSendReq 发送消息请求结构
type DouYinIMMessageSendReq struct {
	Content        Content   `json:"content"`         // 消息内容
	ContentList    []Content `json:"content_list"`    // 消息内容列表
	Scene          string    `json:"scene"`           // 场景
	MsgId          string    `json:"msg_id"`          // 消息ID
	ConversationId string    `json:"conversation_id"` // 会话ID
	ToUserId       string    `json:"to_user_id"`      // 接收用户ID
}

// DouYinIMMessageSendRes 发送消息响应结构
type DouYinIMMessageSendRes struct {
	// AllSendSuccess 是否全部发送成功
	AllSendSuccess bool `json:"all_send_success"`
	// MsgId 消息id
	MsgId string `json:"msg_id"`
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	SendMsgStatusList []response.DouYinRes `json:"send_msg_status_list"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	}
}
