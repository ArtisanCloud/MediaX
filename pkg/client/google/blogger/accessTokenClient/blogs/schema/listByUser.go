package schema

// BloggerBlogsListByUserReq 表示 GET /blogger/v3/users/{userId}/blogs API 的请求参数
type BloggerBlogsListByUserReq struct {
	UserID        string `json:"userId"`        // 必填，要提取其博客的用户的ID，可以是'self'或用户的个人资料ID
	FetchUserInfo bool   `json:"fetchUserInfo"` // 可选，是否在响应中包含用户信息
	View          string `json:"view"`          // 可选，详细信息级别："ADMIN"(管理员级别)、"AUTHOR"(作者级别)、"READER"(读者级别)
}

// BlogUserInfo 表示每个用户的博客信息
type BlogUserInfo struct {
	Kind           string `json:"kind"`           // 资源类型，固定为"blogger#blogPerUserInfo"
	UserID         string `json:"userId"`         // 用户ID
	BlogID         string `json:"blogId"`         // 博客ID
	PhotosAlbumKey string `json:"photosAlbumKey"` // 照片相册密钥
	HasAdminAccess bool   `json:"hasAdminAccess"` // 是否具有管理员访问权限
}

// BlogUserInfos 表示博客用户信息资源
type BlogUserInfos struct {
	Kind         string       `json:"kind"`           // 资源类型，固定为"blogger#blogUserInfo"
	Blog         BloggerBlog  `json:"blog"`           // 博客信息
	BlogUserInfo BlogUserInfo `json:"blog_user_info"` // 博客用户详细信息
}

// Photo 表示用户照片资源
type Photo struct {
	URL string `json:"url"` // 照片URL
}

// BloggerBlogsListByUserRes 表示 GET /blogger/v3/users/{userId}/blogs API 的响应
type BloggerBlogsListByUserRes struct {
	Kind          string          `json:"kind"`          // 资源类型，固定为"blogger#blogList"
	Items         []BloggerBlog   `json:"items"`         // 博客列表
	BlogUserInfos []BlogUserInfos `json:"blogUserInfos"` // 博客用户信息列表，仅当fetchUserInfo=true时返回
}
