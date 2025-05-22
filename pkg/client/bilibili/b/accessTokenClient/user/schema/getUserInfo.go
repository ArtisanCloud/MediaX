package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type BiliBiliUserGetUserInfoRes struct {
	response.BiliBiliRes
	Data *BiliBiliUserInfoData `json:"data"` // 用户信息数据
}

type BiliBiliUserInfoData struct {
	Face   string `json:"face"`   // 用户头像
	Name   string `json:"name"`   // 用户昵称
	Openid string `json:"openid"` // 用户openid
}
