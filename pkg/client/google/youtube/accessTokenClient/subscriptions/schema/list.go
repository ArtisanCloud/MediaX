package schema

import "github.com/ArtisanCloud/MediaXCore/utils/object"

// YouTubeSubscriptionsListReq 获取订阅列表请求参数
type YouTubeSubscriptionsListReq struct {
	Part                          string `json:"part"`                                    // 指定返回的资源部分（必填，如 snippet,contentDetails 等）
	ChannelId                     string `json:"channelId,omitempty"`                     // 频道ID（可选，与 id、mine、myRecentSubscribers、mySubscribers 互斥）
	Id                            string `json:"id,omitempty"`                            // 订阅ID列表（可选，以逗号分隔）
	Mine                          bool   `json:"mine,omitempty"`                          // 是否获取经过身份验证的用户的订阅（可选）
	MyRecentSubscribers           bool   `json:"myRecentSubscribers,omitempty"`           // 是否按时间倒序获取订阅者的 Feed（可选）
	MySubscribers                 bool   `json:"mySubscribers,omitempty"`                 // 是否获取订阅者的 Feed（可选）
	ForChannelId                  string `json:"forChannelId,omitempty"`                  // 频道ID列表（可选，以逗号分隔）
	MaxResults                    int    `json:"maxResults,omitempty"`                    // 返回的最大结果数（可选，默认5，最大50）
	OnBehalfOfContentOwner        string `json:"onBehalfOfContentOwner,omitempty"`        // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	OnBehalfOfContentOwnerChannel string `json:"onBehalfOfContentOwnerChannel,omitempty"` // 内容所有者频道（可选，仅供 YouTube 内容合作伙伴使用）
	Order                         string `json:"order,omitempty"`                         // 排序方式（可选，如 alphabetical,relevance,unread）
	PageToken                     string `json:"pageToken,omitempty"`                     // 分页令牌（可选）
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type ResourceId struct {
	Kind      string `json:"kind"`
	ChannelId string `json:"channelId"`
}

type Snippet struct {
	PublishedAt  string          `json:"publishedAt"`
	ChannelTitle string          `json:"channelTitle"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	ResourceId   ResourceId      `json:"resourceId"`
	ChannelId    string          `json:"channelId"`
	Thumbnails   *object.HashMap `json:"thumbnails"`
}

type ContentDetails struct {
	TotalItemCount int    `json:"totalItemCount"`
	NewItemCount   int    `json:"newItemCount"`
	ActivityType   string `json:"activityType"`
}

type SubscriberSnippet struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ChannelId   string          `json:"channelId"`
	Thumbnails  *object.HashMap `json:"thumbnails"`
}

type Subscription struct {
	Kind              string            `json:"kind"`
	Etag              string            `json:"etag"`
	Id                string            `json:"id"`
	Snippet           Snippet           `json:"snippet"`
	ContentDetails    ContentDetails    `json:"contentDetails"`
	SubscriberSnippet SubscriberSnippet `json:"subscriberSnippet"`
}

type YouTubeSubscriptionsListRes struct {
	Kind          string         `json:"kind"`
	Etag          string         `json:"etag"`
	NextPageToken string         `json:"nextPageToken"`
	PrevPageToken string         `json:"prevPageToken"`
	PageInfo      PageInfo       `json:"pageInfo"`
	Items         []Subscription `json:"items"`
}
