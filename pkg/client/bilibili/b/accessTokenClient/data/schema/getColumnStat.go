package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliDataGetColumnStatReq 表示 GET /arcopen/fn/data/art/stat API 的请求参数
type BiliBiliDataGetColumnStatReq struct {
	IDs string `json:"ids"` // 文章ids；若要查询多个文章数据，多个id用英文逗号分割
}

type ColumnCategory struct {
	ID       int    `json:"id"`        // 分类id
	ParentID int    `json:"parent_id"` // 父级分类id
	Name     string `json:"name"`      // 分类名称
}

type ColumnStats struct {
	View     int `json:"view"`     // 浏览数
	Favorite int `json:"favorite"` // 收藏数
	Like     int `json:"like"`     // 喜欢数
	Dislike  int `json:"dislike"`  // 不喜欢数
	Reply    int `json:"reply"`    // 评论数
	Share    int `json:"share"`    // 分享数
	Coin     int `json:"coin"`     // 硬币数
}

type ColumnList struct {
	ID          int    `json:"id"`           // 文集id
	Name        string `json:"name"`         // 文集名称
	ImageURL    string `json:"image_url"`    // 文集封面图片
	Summary     string `json:"summary"`      // 文集简介
	Words       int    `json:"words"`        // 文集总字数
	PublishTime int    `json:"publish_time"` // 最后发表文章时间戳
	Ctime       int    `json:"ctime"`        // 文集创建时间戳
	UpdateTime  int    `json:"update_time"`  // 最后修改时间戳
}

type ColumnData struct {
	ID           int            `json:"id"`             // 文章id
	Category     ColumnCategory `json:"category"`       // 分类信息
	Title        string         `json:"title"`          // 文章标题
	Summary      string         `json:"summary"`        // 文章简介
	BannerURL    string         `json:"banner_url"`     // 文章头图
	TemplateID   int            `json:"template_id"`    // 模板类型
	State        int            `json:"state"`          // 文章状态
	Reason       string         `json:"reason"`         // 文章打回理由
	ImageURLs    []string       `json:"image_urls"`     // 封面图
	PublishTime  int            `json:"publish_time"`   // 文章发表时间
	Ctime        int            `json:"ctime"`          // 文章创建时间戳
	Stats        ColumnStats    `json:"stats"`          // 统计数据
	Words        int            `json:"words"`          // 字数
	List         *ColumnList    `json:"list"`           // 文集信息
	TopVideoBvid string         `json:"top_video_bvid"` // 头部视频bvid
}

// BiliBiliDataGetColumnStatRes 表示 GET /arcopen/fn/data/art/stat API 的响应
// 接口文档参考：https://member.bilibili.com/arcopen/fn/data/art/stat
type BiliBiliDataGetColumnStatRes struct {
	response.BiliBiliRes
	Data map[string]ColumnData `json:"data"` // 文章数据，key为文章ID
}
