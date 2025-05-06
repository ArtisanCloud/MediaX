package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// MaterialAddMaterialReq 添加素材请求结构体
type MaterialAddMaterialReq struct {
	Type         string          `json:"type"`         // 素材类型
	Media        *object.HashMap `json:"media"`        // 媒体文件
	Title        string          `json:"title"`        // 标题
	Introduction string          `json:"introduction"` // 简介
}

// ---------------------------------------------------

// Article 文章结构体
type Article struct {
	Title              string `json:"title"`                 // 标题
	ThumbMediaID       string `json:"thumb_media_id"`        // 封面图片素材ID
	Author             string `json:"author"`                // 作者
	Digest             string `json:"digest"`                // 摘要
	ShowCoverPic       string `json:"show_cover_pic"`        // 是否显示封面图片
	Content            string `json:"content"`               // 内容
	ContentSourceUrl   string `json:"content_source_url"`    // 原文链接
	NeedOpenComment    string `json:"need_open_comment"`     // 是否需要打开评论
	OnlyFansCanComment string `json:"only_fans_can_comment"` // 是否只有粉丝可以评论
}

// AddArticlesReq 添加图文消息请求结构体
type AddArticlesReq struct {
	Articles []*Article `json:"articles"` // 文章列表
}

// MaterialAddMaterialRes 添加素材响应结构体
type MaterialAddMaterialRes struct {
	response.OfficialAccountRes
	MediaID string `json:"media_id"` // 媒体ID
	URL     string `json:"url"`      // 素材URL

}
