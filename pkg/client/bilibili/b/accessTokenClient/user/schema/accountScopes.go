package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliUserRes 表示B站用户权限查询API的响应
type BiliBiliUserRes struct {
	response.BiliBiliRes
	Data *BiliBiliUserData `json:"data"` // 用户数据
}

// BiliBiliUserData 表示用户数据
type BiliBiliUserData struct {
	OpenId string   `json:"openid"` // 用户OpenID
	Scopes []string `json:"scopes"` // 权限列表
}
