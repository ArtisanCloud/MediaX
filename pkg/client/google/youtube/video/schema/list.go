package schema

type PageInfo struct {
	TotalResults   int         `json:"totalResults"`
	ResultsPerPage interface{} `json:"resultsPerPage"`
}

// YouTubeVideosRequest 表示 GET /youtube/v3/videos API 的请求参数
type ListReq struct {
	Part                   string  `url:"part"`                             // 必填，指定 API 响应包含的 video 资源属性
	Chart                  *string `url:"chart,omitempty"`                  // 可选，指定要检索的热门视频图表
	ID                     *string `url:"id,omitempty"`                     // 可选，指定要检索的视频 ID 列表
	MyRating               *string `url:"myRating,omitempty"`               // 可选，仅限授权用户，筛选喜欢或不喜欢的视频
	HL                     *string `url:"hl,omitempty"`                     // 可选，指定本地化语言代码
	MaxHeight              *int    `url:"maxHeight,omitempty"`              // 可选，指定嵌入式播放器的最大高度
	MaxResults             *int    `url:"maxResults,omitempty"`             // 可选，指定最大返回项数 (1-50, 默认5)
	MaxWidth               *int    `url:"maxWidth,omitempty"`               // 可选，指定嵌入式播放器的最大宽度
	OnBehalfOfContentOwner *string `url:"onBehalfOfContentOwner,omitempty"` // 可选，代表内容所有者执行操作
	PageToken              *string `url:"pageToken,omitempty"`              // 可选，分页标记
	RegionCode             *string `url:"regionCode,omitempty"`             // 可选，指定区域代码 (ISO 3166-1 alpha-2)
	VideoCategoryID        *string `url:"videoCategoryId,omitempty"`        // 可选，指定要检索的类别 (默认 0)
}

type ListRes struct {
	Kind          string   `json:"kind"`
	Etag          string   `json:"etag"`
	NextPageToken string   `json:"nextPageToken"`
	PrevPageToken string   `json:"prevPageToken"`
	PageInfo      PageInfo `json:"pageInfo"`
	Items         []Video  `json:"items"`
}
