package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleAddReq 表示 POST /arcopen/fn/article/add API 的请求参数
type BiliBiliArticleAddReq struct {
	Title         string `json:"title"`           // 文章标题，建议30字以内，最多不超过40字
	Category      int    `json:"category"`        // 文章分类，通过categories接口获取，填写子分类id
	TemplateID    int    `json:"template_id"`     // 选择稿件模板：3-封面三图；4-封面单图或无封面图且banner_url/top_video_bvid不为空；5-封面无图且banner_url/top_video_bvid都为空，平台生成默认封面
	Summary       string `json:"summary"`         // 文章简介，字数不少于200字
	Content       string `json:"content"`         // 正文内容，字数200-40000，或者至少添加三张图片
	BannerURL     string `json:"banner_url"`      // 头图url，与top_video_bvid可都不填写，或者只能二选一
	Original      int    `json:"original"`        // 是否原创：0-非原创 1-原创 默认0
	ImageURLs     string `json:"image_urls"`      // 封面默认使用正文文字图片，多个图片用英文逗号分割
	Tags          string `json:"tags"`            // 用户自定义标签，多标签情况下用英文逗号分割
	ListID        int    `json:"list_id"`         // 文集id
	UpClosedReply int    `json:"up_closed_reply"` // 是否up主关闭评论区：0-否 1-是 默认0
	TopVideoBVID  string `json:"top_video_bvid"`  // 头部视频bvid，与banner_url可都不填写，或者只能二选一
}

type ArticleAddData struct {
	ID int `json:"id"` // 新建的文章id
}

// BiliBiliArticleAddRes 表示 POST /arcopen/fn/article/add API 的响应
type BiliBiliArticleAddRes struct {
	response.BiliBiliRes
	Data ArticleAddData `json:"data"` // 响应数据
}
