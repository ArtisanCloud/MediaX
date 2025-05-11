package schema

// YoutubeMembersListReq 表示 GET /youtube/v3/members API 的请求参数
type YoutubeMembersListReq struct {
	Part                    string `json:"part"`                              // 必填，指定 API 响应包含的 member 资源属性
	Mode                    string `json:"mode,omitempty"`                    // 可选，指定返回的成员类型（all_current 或 updates）
	MaxResults              int    `json:"maxResults,omitempty"`              // 可选，指定结果集中应返回的商品数量上限（0-1000，默认5）
	PageToken               string `json:"pageToken,omitempty"`               // 可选，用于标识结果集中应返回的特定网页
	HasAccessToLevel        string `json:"hasAccessToLevel,omitempty"`        // 可选，指定结果集中的成员应具备的最低级别 ID
	FilterByMemberChannelId string `json:"filterByMemberChannelId,omitempty"` // 可选，以逗号分隔的频道 ID 列表，用于检查特定用户的成员资格状态
}

type MemberDetails struct {
	ChannelId       string `json:"channelId"`       // 成员频道 ID
	ChannelUrl      string `json:"channelUrl"`      // 成员频道 URL
	DisplayName     string `json:"displayName"`     // 成员显示名称
	ProfileImageUrl string `json:"profileImageUrl"` // 成员头像 URL
}
type MembershipsDetails struct {
	AccessibleLevels       []string `json:"accessibleLevels"`       // 可访问的会员等级 ID 列表
	HighestAccessibleLevel string   `json:"highestAccessibleLevel"` // 最高可访问的会员等级 ID
	MembershipsDuration    string   `json:"membershipsDuration"`    // 会员时长
	MemberSince            string   `json:"memberSince"`            // 成为会员的时间
	Status                 string   `json:"status"`                 // 会员状态
	Type                   string   `json:"type"`                   // 会员类型
}

type MemberSnippet struct {
	CreatorChannelId   string             `json:"creatorChannelId"` // 创建者频道 ID
	MemberDetails      MemberDetails      `json:"memberDetails"`
	MembershipsDetails MembershipsDetails `json:"membershipsDetails"`
}

type Member struct {
	Kind    string        `json:"kind"`    // 资源类型
	Etag    string        `json:"etag"`    // 资源的 ETag
	Snippet MemberSnippet `json:"snippet"` // 成员的基本信息
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`   // 总结果数
	ResultsPerPage int `json:"resultsPerPage"` // 每页结果数
}

// YoutubeMembersListRes 表示 GET /youtube/v3/members API 的响应
type YoutubeMembersListRes struct {
	Kind          string   `json:"kind"`          // 资源类型
	Etag          string   `json:"etag"`          // 资源的 ETag
	NextPageToken string   `json:"nextPageToken"` // 下一页令牌
	PageInfo      PageInfo `json:"pageInfo"`      // 分页信息
	Items         []Member `json:"items"`         // 成员列表
}
