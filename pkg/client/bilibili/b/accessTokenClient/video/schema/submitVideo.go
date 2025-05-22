package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliVideoSubmitVideoReq 表示视频稿件提交请求
type BiliBiliVideoSubmitVideoReq struct {
	UploadToken string `json:"upload_token"`         // 上传令牌
	Title       string `json:"title"`                // 稿件标题
	Cover       string `json:"cover,omitempty"`      // 封面地址
	Tid         int    `json:"tid"`                  // 分区ID
	NoReprint   int    `json:"no_reprint,omitempty"` // 是否允许转载
	Desc        string `json:"desc,omitempty"`       // 视频描述
	Tag         string `json:"tag"`                  // 视频标签
	Copyright   int    `json:"copyright"`            // 版权类型
	Source      string `json:"source,omitempty"`     // 转载来源
	TopicID     int    `json:"topic_id,omitempty"`   // 话题ID
}

// SubmitVideoData 表示提交视频响应数据
type SubmitVideoData struct {
	ResourceID string `json:"resource_id"` // 稿件唯一ID
}

// BiliBiliVideoSubmitVideoRes 表示视频稿件提交响应
type BiliBiliVideoSubmitVideoRes struct {
	response.BiliBiliRes
	Data SubmitVideoData `json:"data"`
}
