package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentVideoBasicInfoReq 视频基本信息请求参数
type DouYinContentVideoBasicInfoReq struct {
	// ItemIds item_id数组，仅能查询access_token对应用户上传的视频
	ItemIds []string `json:"item_ids"`

	// VideoIds 明文item_id数组，如果和item_ids同时传入，优先处理video_ids
	VideoIds []string `json:"video_ids"`
}

type Video struct {
	CreateTime  int    `json:"create_time"`
	ItemId      string `json:"item_id"`
	MediaType   int    `json:"media_type"`
	Title       string `json:"title"`
	VideoId     string `json:"video_id"`
	VideoStatus int    `json:"video_status"`
	Cover       string `json:"cover"`
}

type DouYinContentVideoBasicInfoRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// List 视频信息列表，包含多个视频的详细信息
		List []Video `json:"list"`
	} `json:"data"`
}
