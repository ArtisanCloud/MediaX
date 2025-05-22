package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// CommonMsg 通用消息结构体，用于传递视频投稿所需的基本信息
type CommonMsg struct {
	Source           string   `json:"source"`             // 视频来源
	Cover            string   `json:"cover,omitempty"`    // 封面图片URL
	Title            string   `json:"title,omitempty"`    // 视频标题
	TypeID           int      `json:"type_id"`            // 视频分类ID
	TopicID          int      `json:"topic_id"`           // 视频话题ID
	DeliveryMode     int      `json:"delivery_mode"`      // 视频投稿方式
	VideoMaterialUrl []string `json:"video_material_url"` // 视频素材URL列表
}

// BiliBiliVideoAddShareUrlReq 获取唤起哔哩哔哩客户端投稿页面URL的请求结构体
type BiliBiliVideoAddShareUrlReq struct {
	SceneCode string `json:"scene_code"` // 场景码
	BizCode   string `json:"biz_code"`   // 业务码
	CommonMsg string `json:"common_msg"` // 通用消息
}

// AddShareUrlData 获取唤起哔哩哔哩客户端投稿页面URL的返回数据结构体
type AddShareUrlData struct {
	LinkURL string `json:"link_url"` // 链接URL
}

// BiliBiliVideoAddShareUrlRes 获取唤起哔哩哔哩客户端投稿页面URL的响应结构体
type BiliBiliVideoAddShareUrlRes struct {
	response.BiliBiliRes
	Data AddShareUrlData `json:"data"`
}
