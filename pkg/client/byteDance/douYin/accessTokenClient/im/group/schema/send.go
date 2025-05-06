package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/message/schema"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
)

// DouYinIMGroupSendReq 发送群消息请求参数
type DouYinIMGroupSendReq struct {
	GroupId    string         `json:"group_id"`    // 群ID，来源于群消息webhook事件或查询群信息接口
	Content    schema.Content `json:"content"`     // 消息体，支持文本、图片、视频、留资卡片、小程序引导卡片等
	GroupToken string         `json:"group_token"` // 群Token
}

// DouYinIMGroupSendRes 发送群消息响应参数
type DouYinIMGroupSendRes struct {
	// MsgId 消息ID，发送成功后返回的消息ID。
	MsgId string `json:"msg_id"`
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
