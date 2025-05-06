package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"

// Article 文章结构体
type Article struct {
	Title              string `json:"title"`                 // 标题
	Author             string `json:"author"`                // 作者
	Digest             string `json:"digest"`                // 摘要
	Content            string `json:"content"`               // 内容
	ContentSourceUrl   string `json:"content_source_url"`    // 原文链接
	ThumbMediaId       string `json:"thumb_media_id"`        // 封面图片素材ID
	NeedOpenComment    int    `json:"need_open_comment"`     // 是否需要打开评论
	OnlyFansCanComment int    `json:"only_fans_can_comment"` // 是否只有粉丝可以评论
}

// DraftAddReq 添加草稿请求结构体
type DraftAddReq struct {
	Articles []*Article `json:"articles"` // 文章列表
}

// DraftUpdateReq 更新草稿请求结构体
type DraftUpdateReq struct {
	MediaId  string   `json:"media_id"` // 媒体ID
	Index    int      `json:"index"`    // 文章索引
	Articles *Article `json:"articles"` // 文章内容
}

// BatchGetReq 批量获取草稿请求结构体
type BatchGetReq struct {
	Offset    int `json:"offset"`     // 偏移量
	Count     int `json:"count"`      // 获取数量
	NoContent int `json:"no_content"` // 是否不返回内容
}

// DraftAddRes 添加草稿响应结构体
type DraftAddRes struct {
	response.OfficialAccountRes

	MediaID string `json:"media_id"` // 媒体ID
}

// NewsItem 图文消息项结构体
type NewsItem struct {
	Title              string `json:"title"`                 // 标题
	Author             string `json:"author"`                // 作者
	Digest             string `json:"digest"`                // 摘要
	Content            string `json:"content"`               // 内容
	ContentSourceUrl   string `json:"content_source_url"`    // 原文链接
	ThumbMediaId       string `json:"thumb_media_id"`        // 封面图片素材ID
	ThumbUrl           string `json:"thumb_url"`             // 封面图片URL
	ShowCoverPic       int    `json:"show_cover_pic"`        // 是否显示封面图片
	NeedOpenComment    int    `json:"need_open_comment"`     // 是否需要打开评论
	OnlyFansCanComment int    `json:"only_fans_can_comment"` // 是否只有粉丝可以评论
	Url                string `json:"url"`                   // 文章链接
	IsDeleted          bool   `json:"is_deleted"`            // 是否已删除
}

// DraftGetRes 获取草稿响应结构体
type DraftGetRes struct {
	response.OfficialAccountRes

	NewsItem   []*NewsItem `json:"news_item"`   // 图文消息列表
	CreateTime int64       `json:"create_time"` // 创建时间
	UpdateTime int64       `json:"update_time"` // 更新时间
}

// DraftCountRes 草稿总数响应结构体
type DraftCountRes struct {
	TotalCount int `json:"total_count"` // 草稿总数
}

// Content 内容结构体
type Content struct {
	NewsItem   []*NewsItem `json:"news_item"`   // 图文消息列表
	CreateTime int64       `json:"create_time"` // 创建时间
	UpdateTime int64       `json:"update_time"` // 更新时间
}

// Item 草稿项结构体
type Item struct {
	MediaId    string   `json:"media_id"`    // 媒体ID
	ArticleId  string   `json:"article_id"`  // 文章ID
	Content    *Content `json:"content"`     // 内容
	UpdateTime int64    `json:"update_time"` // 更新时间
}

// BatchGetRes 批量获取草稿响应结构体
type BatchGetRes struct {
	response.OfficialAccountRes

	TotalCount int     `json:"total_count"` // 草稿总数
	ItemCount  int     `json:"item_count"`  // 本次获取的草稿数量
	Item       []*Item `json:"item"`        // 草稿列表
}

// CheckSwitchRes 检查开关状态响应结构体
type CheckSwitchRes struct {
	response.OfficialAccountRes

	TotalCount int `json:"total_count"` // 总数
	ItemCount  int `json:"item_count"`  // 本次获取的数量
	IsOpen     int `json:"is_open"`     // 是否打开
}
