package schema

import "github.com/ArtisanCloud/MediaXCore/utils/object"

// YouTubeSearchListReq 搜索列表请求参数
type YouTubeSearchListReq struct {
	Part                      string `json:"part"`                                // 指定返回的资源部分（必填，如 snippet）
	ForContentOwner           bool   `json:"forContentOwner,omitempty"`           // 是否限制为内容所有者拥有的视频（可选，仅供 YouTube 内容合作伙伴使用）
	ForDeveloper              bool   `json:"forDeveloper,omitempty"`              // 是否限制为开发者上传的视频（可选）
	ForMine                   bool   `json:"forMine,omitempty"`                   // 是否限制为经过身份验证的用户拥有的视频（可选）
	ChannelId                 string `json:"channelId,omitempty"`                 // 频道ID（可选）
	ChannelType               string `json:"channelType,omitempty"`               // 频道类型（可选，如 any,show）
	EventType                 string `json:"eventType,omitempty"`                 // 事件类型（可选，如 completed,live,upcoming）
	Location                  string `json:"location,omitempty"`                  // 地理位置坐标（可选，如 (37.42307,-122.08427)）
	LocationRadius            string `json:"locationRadius,omitempty"`            // 地理位置半径（可选，如 1500m,5km）
	MaxResults                int    `json:"maxResults,omitempty"`                // 返回的最大结果数（可选，默认5，最大50）
	OnBehalfOfContentOwner    string `json:"onBehalfOfContentOwner,omitempty"`    // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	Q                         string `json:"q,omitempty"`                         // 搜索查询字词（可选）
	RegionCode                string `json:"regionCode,omitempty"`                // 区域代码（可选，ISO 3166-1 alpha-2 国家/地区代码）
	RelevanceLanguage         string `json:"relevanceLanguage,omitempty"`         // 相关性语言（可选，如 zh-Hans,zh-Hant）
	SafeSearch                string `json:"safeSearch,omitempty"`                // 安全搜索（可选，如 moderate,none,strict）
	TopicId                   string `json:"topicId,omitempty"`                   // 主题ID（可选）
	Type                      string `json:"type,omitempty"`                      // 资源类型（可选，如 video,channel,playlist）
	VideoCaption              string `json:"videoCaption,omitempty"`              // 视频字幕（可选，如 any,closedCaption,none）
	VideoCategoryId           string `json:"videoCategoryId,omitempty"`           // 视频分类ID（可选）
	VideoDefinition           string `json:"videoDefinition,omitempty"`           // 视频清晰度（可选，如 any,high,standard）
	VideoDimension            string `json:"videoDimension,omitempty"`            // 视频维度（可选，如 2d,3d,any）
	VideoDuration             string `json:"videoDuration,omitempty"`             // 视频时长（可选，如 any,long,medium,short）
	VideoEmbeddable           string `json:"videoEmbeddable,omitempty"`           // 视频是否可嵌入（可选，如 any,true）
	VideoLicense              string `json:"videoLicense,omitempty"`              // 视频许可类型（可选，如 any,creativeCommon,youtube）
	VideoPaidProductPlacement string `json:"videoPaidProductPlacement,omitempty"` // 视频是否包含付费宣传内容（可选，如 any,true）
	VideoSyndicated           string `json:"videoSyndicated,omitempty"`           // 视频是否已同步（可选，如 any,true）
	VideoType                 string `json:"videoType,omitempty"`                 // 视频类型（可选，如 any,episode,movie）
	PageToken                 string `json:"pageToken,omitempty"`                 // 分页令牌（可选）
}

// Id 资源ID
type Id struct {
	Kind       string `json:"kind"`       // 资源类型
	VideoId    string `json:"videoId"`    // 视频ID
	ChannelId  string `json:"channelId"`  // 频道ID
	PlaylistId string `json:"playlistId"` // 播放列表ID
}

// Snippet 资源片段信息
type Snippet struct {
	PublishedAt          string          `json:"publishedAt"`          // 发布时间
	ChannelId            string          `json:"channelId"`            // 频道ID
	Title                string          `json:"title"`                // 标题
	Description          string          `json:"description"`          // 描述
	Thumbnails           *object.HashMap `json:"thumbnails"`           // 缩略图
	ChannelTitle         string          `json:"channelTitle"`         // 频道标题
	LiveBroadcastContent string          `json:"liveBroadcastContent"` // 直播状态
}

// SearchResource 搜索结果资源
type SearchResource struct {
	Kind    string      `json:"kind"`    // 资源类型
	Etag    interface{} `json:"etag"`    // 资源的 ETag
	Id      Id          `json:"id"`      // 资源ID
	Snippet Snippet     `json:"snippet"` // 资源片段信息
}

// PageInfo 分页信息
type PageInfo struct {
	TotalResults   int `json:"totalResults"`   // 总结果数
	ResultsPerPage int `json:"resultsPerPage"` // 每页结果数
}

// YouTubeSearchListRes 搜索列表返回结果
type YouTubeSearchListRes struct {
	Kind          string           `json:"kind"`          // 资源类型
	Etag          string           `json:"etag"`          // 资源的 ETag
	NextPageToken string           `json:"nextPageToken"` // 下一页令牌
	PrevPageToken string           `json:"prevPageToken"` // 上一页令牌
	RegionCode    string           `json:"regionCode"`    // 区域代码
	PageInfo      PageInfo         `json:"pageInfo"`      // 分页信息
	Items         []SearchResource `json:"items"`         // 搜索结果列表
}
