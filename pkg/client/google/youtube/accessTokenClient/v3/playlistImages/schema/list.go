package schema

// YoutubePlaylistImagesListReq 获取播放列表图片请求参数
type YoutubePlaylistImagesListReq struct {
	Part                          string `json:"part"`                                    // 指定返回的资源部分（必填，如 snippet 等）
	Id                            string `json:"id,omitempty"`                            // 播放列表图片ID列表（可选，以逗号分隔）
	PlaylistId                    string `json:"playlistId,omitempty"`                    // 播放列表ID（可选，与 id 互斥）
	MaxResults                    int    `json:"maxResults,omitempty"`                    // 返回的最大结果数（可选，默认5，最大50）
	OnBehalfOfContentOwner        string `json:"onBehalfOfContentOwner,omitempty"`        // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	OnBehalfOfContentOwnerChannel string `json:"onBehalfOfContentOwnerChannel,omitempty"` // 内容所有者频道（可选，仅供 YouTube 内容合作伙伴使用）
	PageToken                     string `json:"pageToken,omitempty"`                     // 分页令牌（可选）
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type Snippet struct {
	PlaylistId string `json:"playlistId"`
	Type       string `json:"type"`
	Width      string `json:"width"`
	Height     string `json:"height"`
}

type PlaylistImages struct {
	Kind    string  `json:"kind"`
	Id      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type YoutubePlaylistImagesListRes struct {
	Kind          string           `json:"kind"`
	NextPageToken string           `json:"nextPageToken"`
	PrevPageToken string           `json:"prevPageToken"`
	PageInfo      PageInfo         `json:"pageInfo"`
	Items         []PlaylistImages `json:"items"`
}
