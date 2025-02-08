package response

import "github.com/ArtisanCloud/MediaX/v1/internal/kernel/response"

type WeChatAccessTokenRes struct {
	response.AccessTokenRes
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}
