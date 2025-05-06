package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
)

// NewsItems 图文消息项结构体
type NewsItems struct {
	Title            string `json:"title"`              // 标题
	ThumbMediaId     string `json:"thumb_media_id"`     // 封面图片素材ID
	ShowCoverPic     int8   `json:"show_cover_pic"`     // 是否显示封面图片
	Author           string `json:"author"`             // 作者
	Digest           string `json:"digest"`             // 摘要
	Content          string `json:"content"`            // 内容
	Url              string `json:"url"`                // 文章链接
	ContentSourceUrl string `json:"content_source_url"` // 原文链接
}

// MaterialGetNewsRes 获取图文消息响应结构体
type MaterialGetNewsRes struct {
	response.OfficialAccountRes

	NewsItem []NewsItems `json:"news_item"` // 图文消息列表
}

// ------------------------------------

// MaterialGetVideoRes 获取视频素材响应结构体
type MaterialGetVideoRes struct {
	Title       string `json:"title"`       // 视频标题
	Description string `json:"description"` // 视频描述
	DownUrl     string `json:"down_url"`    // 视频下载地址
}
