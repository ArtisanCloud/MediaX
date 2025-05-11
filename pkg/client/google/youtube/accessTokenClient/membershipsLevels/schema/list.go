package schema

// YoutubeMembershipsLevelsListReq 表示 GET /youtube/v3/membershipsLevels API 的请求参数
type YoutubeMembershipsLevelsListReq struct {
	Part string `json:"part"` // 必填，指定 API 响应包含的 membershipsLevel 资源属性（如 id,snippet）
}

type LevelDetails struct {
	DisplayName string `json:"displayName"` // 会员等级显示名称
}

type Snippet struct {
	CreatorChannelId string `json:"creatorChannelId"` // 创建者频道 ID
	LevelDetails    LevelDetails `json:"levelDetails"` // 等级详情
}

type MembershipsLevels struct {
	Kind    string `json:"kind"`    // 资源类型
	Etag    string `json:"etag"`    // 资源的 ETag
	Id      string `json:"id"`      // 会员等级 ID
	Snippet Snippet `json:"snippet"` // 会员等级的基本信息
}

type PageInfo  struct {
	TotalResults   int `json:"totalResults"`   // 结果总数
	ResultsPerPage int `json:"resultsPerPage"` // 每页结果数
}

// YoutubeMembershipsLevelsListRes 表示 GET /youtube/v3/membershipsLevels API 的响应
type YoutubeMembershipsLevelsListRes struct {
	Kind     string `json:"kind"`     // 资源类型
	Etag     string `json:"etag"`     // 资源的 ETag
	PageInfo PageInfo `json:"pageInfo"` // 分页信息
	Items [] MembershipsLevels `json:"items"` // 会员等级列表
}

