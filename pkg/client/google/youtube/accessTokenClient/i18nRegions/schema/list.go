package schema

// YoutubeI18nRegionsListReq 表示 GET /youtube/v3/i18nRegions API 的请求参数
// YoutubeI18nRegionsListReq 表示获取YouTube国际区域列表的请求参数
// 文档参考：https://developers.google.cn/youtube/v3/docs/i18nRegions/list?hl=zh-cn
type YoutubeI18nRegionsListReq struct {
	// Part 指定API响应包含的资源属性，必填
	// 示例值："snippet"，表示返回区域的基本信息
	Part string `json:"part"`

	// HL 指定API响应中文本值应使用的语言，可选
	// 示例值："zh-CN"，表示使用简体中文
	HL string `json:"hl,omitempty"`
}

// I18nRegion represents a single i18n region resource
// I18nRegion 表示一个YouTube国际区域资源
// 包含区域的基本信息和元数据
type I18nRegion struct {
	// Kind 资源类型，固定值："youtube#i18nRegion"
	Kind string `json:"kind"`

	// Etag 资源的ETag，用于缓存控制
	Etag string `json:"etag"`

	// Id 区域ID，ISO 3166-1 alpha-2格式的国家代码
	Id string `json:"id"`

	// Snippet 包含区域的基本信息
	Snippet I18nRegionSnippet `json:"snippet"`
}

// I18nRegionSnippet contains basic information about the region
// I18nRegionSnippet 包含YouTube国际区域的基本信息
type I18nRegionSnippet struct {
	// HL 区域的语言代码
	HL string `json:"hl"`

	// Name 区域的显示名称，根据HL参数指定的语言返回
	Name string `json:"name"`

	// Gl 区域的地理位置代码，ISO 3166-1 alpha-2格式
	Gl string `json:"gl"`
}

// YoutubeI18nRegionsListRes represents the response for GET /youtube/v3/i18nRegions API
// YoutubeI18nRegionsListRes 表示获取YouTube国际区域列表的响应
// 包含区域列表和分页信息
type YoutubeI18nRegionsListRes struct {
	// Kind 资源类型，固定值："youtube#i18nRegionListResponse"
	Kind string `json:"kind"`

	// Etag 资源的ETag，用于缓存控制
	Etag string `json:"etag"`

	// Items 国际区域列表，包含所有支持的区域
	Items []I18nRegion `json:"items"`
}
