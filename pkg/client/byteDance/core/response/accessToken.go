package response

import "github.com/ArtisanCloud/MediaX/internal/kernel/response"

type ByteDanceAccessTokenRes struct {
	response.AccessTokenRes
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}
