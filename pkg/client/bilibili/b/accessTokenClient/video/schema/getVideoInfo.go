package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type BiliBiliVideoGetVideoInfoReq struct {
	ResourceID string `json:"resource_id"` // 稿件唯一ID
}

type VideoInfo struct {
	CId       int    `json:"cid"`        // 视频ID
	Filename  string `json:"filename"`   // 视频文件名
	Duration  int    `json:"duration"`   // 视频时长
	ShareURL  string `json:"share_url"`  // 播放页链接
	IframeURL string `json:"iframe_url"` // 内嵌iframe链接
}

type AdditInfo struct {
	State        int    `json:"state"`         // 审核状态
	StateDesc    string `json:"state_desc"`    // 审核状态描述
	RejectReason string `json:"reject_reason"` // 审核打回理由
}

type GetVideoInfoData struct {
	ResourceID string    `json:"resource_id"` // 稿件ID
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

type BiliBiliVideoGetVideoInfoRes struct {
	response.BiliBiliRes
	Data GetVideoInfoData `json:"data"`
}
