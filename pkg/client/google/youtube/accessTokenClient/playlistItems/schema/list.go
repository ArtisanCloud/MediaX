package schema

import "github.com/ArtisanCloud/MediaXCore/utils/object"

// YouTubePlaylistItemsListReq 获取播放列表项请求参数
type YouTubePlaylistItemsListReq struct {
	Part                   string `json:"part"`                             // 指定返回的资源部分（必填，如 snippet,contentDetails 等）
	Id                     string `json:"id,omitempty"`                     // 播放列表项ID列表（可选，以逗号分隔）
	PlaylistId             string `json:"playlistId,omitempty"`             // 播放列表ID（可选，与 id 互斥）
	MaxResults             int    `json:"maxResults,omitempty"`             // 返回的最大结果数（可选，默认5，最大50）
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	PageToken              string `json:"pageToken,omitempty"`              // 分页令牌（可选）
	VideoId                string `json:"videoId,omitempty"`                // 视频ID（可选）
}

// PageInfo 分页信息
type PageInfo struct {
	TotalResults   int `json:"totalResults"`   // 总结果数
	ResultsPerPage int `json:"resultsPerPage"` // 每页结果数
}

// ResourceId 资源ID
type ResourceId struct {
	Kind    string `json:"kind"`    // 资源类型
	VideoId string `json:"videoId"` // 视频ID
}

// Snippet 播放列表项片段信息
type Snippet struct {
	PublishedAt            string          `json:"publishedAt"`            // 发布时间
	ChannelId              string          `json:"channelId"`              // 频道ID
	Title                  string          `json:"title"`                  // 标题
	Description            string          `json:"description"`            // 描述
	Thumbnails             *object.HashMap `json:"thumbnails"`             // 缩略图
	ChannelTitle           string          `json:"channelTitle"`           // 频道标题
	VideoOwnerChannelTitle string          `json:"videoOwnerChannelTitle"` // 视频所有者频道标题
	VideoOwnerChannelId    string          `json:"videoOwnerChannelId"`    // 视频所有者频道ID
	PlaylistId             string          `json:"playlistId"`             // 播放列表ID
	Position               int             `json:"position"`               // 位置
	ResourceId             ResourceId      `json:"resourceId"`             // 资源ID
}

// ContentDetails 播放列表项内容详情
type ContentDetails struct {
	VideoId          string `json:"videoId"`          // 视频ID
	StartAt          string `json:"startAt"`          // 开始时间
	EndAt            string `json:"endAt"`            // 结束时间
	Note             string `json:"note"`             // 备注
	VideoPublishedAt string `json:"videoPublishedAt"` // 视频发布时间
}

// Status 播放列表项状态
type Status struct {
	PrivacyStatus string `json:"privacyStatus"` // 隐私状态
}

// PlayListItem 播放列表项
type PlayListItem struct {
	Kind           string         `json:"kind"`           // 资源类型
	Etag           string         `json:"etag"`           // 资源的 ETag
	Id             string         `json:"id"`             // 播放列表项ID
	Snippet        Snippet        `json:"snippet"`        // 播放列表项片段信息
	ContentDetails ContentDetails `json:"contentDetails"` // 播放列表项内容详情
	Status         Status         `json:"status"`         // 播放列表项状态
}

// YouTubePlaylistItemsListRes 播放列表项返回结果
type YouTubePlaylistItemsListRes struct {
	Kind          string         `json:"kind"`          // 资源类型
	Etag          string         `json:"etag"`          // 资源的 ETag
	NextPageToken string         `json:"nextPageToken"` // 下一页令牌
	PrevPageToken string         `json:"prevPageToken"` // 上一页令牌
	PageInfo      PageInfo       `json:"pageInfo"`      // 分页信息
	Items         []PlayListItem `json:"items"`         // 播放列表项列表
}
