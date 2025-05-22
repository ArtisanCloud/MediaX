package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type BiliBiliVideoGetVideoListReq struct {
	Pn     int    `json:"pn"`     // 页码
	Ps     int    `json:"ps"`     // 每页数量
	Status string `json:"status"` // 稿件状态
}

type VideoListItem struct {
	ResourceID string    `json:"resource_id"` // 稿件ID
	Title      string    `json:"title"`       // 稿件标题
	Cover      string    `json:"cover"`       // 封面地址
	Tid        int       `json:"tid"`         // 分区ID
	NoReprint  int       `json:"no_reprint"`  // 是否禁止转载
	Desc       string    `json:"desc"`        // 视频描述
	Tag        string    `json:"tag"`         // 标签
	Copyright  int       `json:"copyright"`   // 版权类型
	VideoInfo  VideoInfo `json:"video_info"`
	AdditInfo  AdditInfo `json:"addit_info"`
	Ctime      int       `json:"ctime"` // 创建时间
	Ptime      int       `json:"ptime"` // 发布时间
}

type PageInfo struct {
	Pn    int `json:"pn"`    // 当前页码
	Ps    int `json:"ps"`    // 每页数量
	Total int `json:"total"` // 总条数
}

type GetVideoListData struct {
	List []VideoListItem `json:"list"` // 稿件列表
	Page PageInfo        `json:"page"` // 分页信息
}

type BiliBiliVideoGetVideoListRes struct {
	response.BiliBiliRes
	Data GetVideoListData `json:"data"`
}
