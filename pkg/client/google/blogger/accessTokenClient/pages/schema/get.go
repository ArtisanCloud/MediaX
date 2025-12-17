package schema

// BloggerPagesGetReq 表示获取单个页面详情的请求参数
//
// 接口文档参考：
// https://developers.google.com/blogger/docs/3.0/reference/pages/get?hl=zh-cn
type BloggerPagesGetReq struct {
	BlogId string `json:"blogId"`         // 包含相应网页的博客ID(必填)
	PageId string `json:"pageId"`         // 要获取的网页ID(必填)
	View   string `json:"view,omitempty"` // 可选，指定返回视图级别(ADMIN/AUTHOR/READER)
}

// BloggerPagesGetRes 表示获取单个页面详情的响应结构
type BloggerPagesGetRes struct {
	PageResource
}
