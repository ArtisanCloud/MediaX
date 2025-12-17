package schema

import "github.com/ArtisanCloud/MediaXCore/utils/object"

// YouTubePlaylistsListReq 获取播放列表请求参数
type YouTubePlaylistsListReq struct {
	Part                          string `json:"part"`                                    // 指定返回的资源部分（必填，如 snippet,contentDetails 等）
	ChannelId                     string `json:"channelId,omitempty"`                     // 频道ID（可选，与 id、mine 互斥）
	Id                            string `json:"id,omitempty"`                            // 播放列表ID列表（可选，以逗号分隔）
	Mine                          bool   `json:"mine,omitempty"`                          // 是否获取经过身份验证的用户拥有的播放列表（可选）
	Hl                            string `json:"hl,omitempty"`                            // 语言代码（可选）
	MaxResults                    int    `json:"maxResults,omitempty"`                    // 返回的最大结果数（可选，默认5，最大50）
	OnBehalfOfContentOwner        string `json:"onBehalfOfContentOwner,omitempty"`        // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	OnBehalfOfContentOwnerChannel string `json:"onBehalfOfContentOwnerChannel,omitempty"` // 内容所有者频道（可选，仅供 YouTube 内容合作伙伴使用）
	PageToken                     string `json:"pageToken,omitempty"`                     // 分页令牌（可选）
}

type Localized struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Snippet struct {
	PublishedAt     string          `json:"publishedAt"`
	ChannelId       string          `json:"channelId"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Thumbnails      *object.HashMap `json:"thumbnails"`
	ChannelTitle    string          `json:"channelTitle"`
	DefaultLanguage string          `json:"defaultLanguage"`
	Localized       Localized       `json:"localized"`
}

type Status struct {
	PrivacyStatus string `json:"privacyStatus"`
	PodcastStatus int    `json:"podcastStatus"`
}

type ContentDetails struct {
	ItemCount int `json:"itemCount"`
}

type Player struct {
	EmbedHtml string `json:"embedHtml"`
}

type Playlist struct {
	Kind           string          `json:"kind"`
	Etag           string          `json:"etag"`
	Id             string          `json:"id"`
	Snippet        Snippet         `json:"snippet"`
	Status         Status          `json:"status"`
	ContentDetails ContentDetails  `json:"contentDetails"`
	Player         Player          `json:"player"`
	Localizations  *object.HashMap `json:"localizations"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type YouTubePlaylistsListRes struct {
	Kind          string     `json:"kind"`
	Etag          string     `json:"etag"`
	NextPageToken string     `json:"nextPageToken"`
	PrevPageToken string     `json:"prevPageToken"`
	PageInfo      PageInfo   `json:"pageInfo"`
	Items         []Playlist `json:"items"`
}
