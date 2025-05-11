package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video/schema"

// VideoCategoryResource 视频分类资源
type VideoCategoryResource struct {
	Kind    string         `json:"kind"`    // 资源类型
	Etag    string         `json:"etag"`    // 资源的 ETag
	Id      string         `json:"id"`      // 视频分类ID
	Snippet schema.Snippet `json:"snippet"` // 视频分类片段信息
}

// YouTubeVideoCategoriesReq 获取视频分类列表请求参数
type YouTubeVideoCategoriesReq struct {
	Part       string `json:"part"`                 // 指定返回的资源部分（必填，如 snippet）
	Id         string `json:"id,omitempty"`         // 视频分类ID（可选）
	RegionCode string `json:"regionCode,omitempty"` // 区域代码（可选）
	HL         string `json:"hl,omitempty"`         // 语言代码（可选）
}

// PageInfo 分页信息
type PageInfo struct {
	TotalResults   int `json:"totalResults"`   // 总结果数
	ResultsPerPage int `json:"resultsPerPage"` // 每页结果数
}

// YouTubeVideoCategoriesRes 获取视频分类列表返回结果
type YouTubeVideoCategoriesRes struct {
	Kind          string                  `json:"kind"`          // 资源类型
	Etag          string                  `json:"etag"`          // 资源的 ETag
	NextPageToken string                  `json:"nextPageToken"` // 下一页令牌
	PrevPageToken string                  `json:"prevPageToken"` // 上一页令牌
	PageInfo      PageInfo                `json:"pageInfo"`      // 分页信息
	Items         []VideoCategoryResource `json:"items"`         // 视频分类列表
}
