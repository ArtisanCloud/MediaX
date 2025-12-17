package schema

// BloggerUsersGetReq 表示 GET /blogger/v3/users/{userId} API 的请求参数
type BloggerUsersGetReq struct {
	// UserId 当前用户的数字ID，或者为"self"表示当前认证用户
	UserId string `json:"userId"`
}

type Locale struct {
	// Language 用户的语言设置
	Language string `json:"language"`
	// Country 用户的国家/地区设置
	Country string `json:"country"`
	// Variant 用户的语言变体设置
	Variant string `json:"variant"`
}

type Blogs struct {
	// SelfLink 此用户的博客网址
	SelfLink string `json:"selfLink"`
}

type UserResource struct {
	// Kind 此实体的种类，始终为"blogger#user"
	Kind string `json:"kind"`
	// Id 此用户的ID
	Id string `json:"id"`
	// Created 此个人资料的创建时间戳（以秒计，自纪元开始）
	Created string `json:"created"`
	// Url 用户的个人资料页面
	Url string `json:"url"`
	// SelfLink 要从中提取此资源的API REST网址
	SelfLink string `json:"selfLink"`
	// Blogs 此用户博客的容器
	Blogs Blogs `json:"blogs"`
	// DisplayName 用户的显示名称
	DisplayName string `json:"displayName"`
	// About 个人资料摘要信息
	About string `json:"about"`
	// Locale 此用户的语言区域
	Locale Locale `json:"locale"`
}

// BloggerUsersGetRes 表示 GET /blogger/v3/users/{userId} API 的响应
type BloggerUsersGetRes struct {
	UserResource
}
