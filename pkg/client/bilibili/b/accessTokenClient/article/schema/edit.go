package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliArticleEditReq 表示 POST /arcopen/fn/article/edit API 的请求参数
type BiliBiliArticleEditReq struct {
	ID         int    `json:"id"`          // 文章ID
	Title      string `json:"title"`       // 文章标题，建议30字以内，最多不超过40字
	Category   int    `json:"category"`    // 文章分类，通过categories接口获取，填写子分类id
	TemplateID int    `json:"template_id"` // 选择稿件模板
	Summary    string `json:"summary"`     // 文章简介，字数不少于200字
	Content    string `json:"content"`     // 正文内容，字数200-40000，或者至少添加三张图片
}

// BiliBiliArticleEditRes 表示 POST /arcopen/fn/article/edit API 的响应
type BiliBiliArticleEditRes struct {
	response.BiliBiliRes
}
