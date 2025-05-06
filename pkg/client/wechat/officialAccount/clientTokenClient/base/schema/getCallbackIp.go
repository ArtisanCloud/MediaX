package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"

// GetCallBackIP 获取微信服务器 IP 地址
type GetCallBackIPRes struct {
	response.OfficialAccountRes

	IPList []string `json:"ip_list"`
}
