package schema

// YoutubeChannelSectionsListReq 表示 GET /youtube/v3/channelSections API 的请求参数
type YoutubeChannelSectionsListReq struct {
	Part                   string  `json:"part"`                             // 必填，指定 API 响应包含的 channelSection 资源属性
	ChannelId              *string `json:"channelId,omitempty"`              // 可选，指定 YouTube 频道 ID
	Id                     *string `json:"id,omitempty"`                     // 可选，指定 channelSection ID 列表
	Mine                   *bool   `json:"mine,omitempty"`                   // 可选，仅限授权用户，检索与经过身份验证的用户相关联的频道版块
	Hl                     *string `json:"hl,omitempty"`                     // 可选，已弃用，指定本地化语言代码
	OnBehalfOfContentOwner *string `json:"onBehalfOfContentOwner,omitempty"` // 可选，代表内容所有者执行操作
}

// YoutubeChannelSectionsListRes 表示 GET /youtube/v3/channelSections API 的响应
// ChannelSection 表示YouTube频道的一个版块
type ChannelSection struct {
	Kind           string                 `json:"kind"`           // API资源类型
	Etag           string                 `json:"etag"`           // 资源的Etag标识
	Id             string                 `json:"id"`             // YouTube用于唯一标识频道版块的ID
	Snippet        ChannelSectionsSnippet `json:"snippet"`        // 频道版块的基本信息
	ContentDetails ContentDetails         `json:"contentDetails"` // 频道版块内容的详细信息
}

// ChannelSectionsSnippet 包含频道版块的基本信息
type ChannelSectionsSnippet struct {
	Type      string `json:"type"`      // 频道版块中的内容类型
	ChannelId string `json:"channelId"` // 此版块所属的频道ID
	Title     string `json:"title"`     // 频道版块的标题
	Position  uint   `json:"position"`  // 此版块在频道中的位置
}

// ContentDetails 包含频道版块的内容信息
type ContentDetails struct {
	Playlists []string `json:"playlists"` // 此版块中包含的播放列表ID列表
	Channels  []string `json:"channels"`  // 此版块中包含的频道ID列表
}

type PageInfo struct {
	TotalResults   uint `json:"totalResults"`   // 总结果数
	ResultsPerPage uint `json:"resultsPerPage"` // 每页结果数
}

type YoutubeChannelSectionsListRes struct {
	Kind          string           `json:"kind"`          // 资源类型
	Etag          string           `json:"etag"`          // 资源的 ETag
	NextPageToken string           `json:"nextPageToken"` // 下一页令牌
	PrevPageToken string           `json:"prevPageToken"` // 上一页令牌
	PageInfo      PageInfo         `json:"pageInfo"`      // 分页信息
	Items         []ChannelSection `json:"items"`         // 频道版块列表
}
