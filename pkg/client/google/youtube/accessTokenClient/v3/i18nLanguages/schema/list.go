package schema

// YoutubeI18nLanguagesListReq 表示 GET /youtube/v3/i18nLanguages API 的请求参数
// 文档参考：https://developers.google.cn/youtube/v3/docs/i18nLanguages/list?hl=zh-cn
type YoutubeI18nLanguagesListReq struct {
	// Part 指定API响应包含的资源属性，必填
	// 示例值："snippet"，表示返回语言的基本信息
	Part string `json:"part"`

	// HL 指定API响应中文本值应使用的语言，可选
	// 示例值："zh-CN"，表示使用简体中文
	HL string `json:"hl,omitempty"`
}

// YoutubeI18nLanguagesListRes 表示 GET /youtube/v3/i18nLanguages API 的响应
// 包含支持的语言列表和分页信息
type YoutubeI18nLanguagesListRes struct {
	// Kind 资源类型，固定值："youtube#i18nLanguageListResponse"
	Kind string `json:"kind"`

	// Etag 资源的ETag，用于缓存控制
	Etag string `json:"etag"`

	// Items 支持的语言列表，包含所有支持的语言
	Items []I18nLanguage `json:"items"`
}

// I18nLanguage represents a single i18n language resource
// I18nLanguage 表示一个YouTube支持的语言资源
type I18nLanguage struct {
	// Kind 资源类型，固定值："youtube#i18nLanguage"
	Kind string `json:"kind"`

	// Etag 资源的ETag，用于缓存控制
	Etag string `json:"etag"`

	// Id 语言ID，BCP-47格式的语言代码
	Id string `json:"id"`

	// Snippet 包含语言的基本信息
	Snippet I18nLanguageSnippet `json:"snippet"`
}

// I18nLanguageSnippet contains basic information about the language
// I18nLanguageSnippet 包含YouTube支持语言的基本信息
type I18nLanguageSnippet struct {
	// HL 语言的语言代码
	HL string `json:"hl"`

	// Name 语言的显示名称，根据HL参数指定的语言返回
	Name string `json:"name"`
}
