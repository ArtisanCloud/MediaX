package response

import "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"

type GetCallBackIPRes struct {
	response.OfficialAccountRes

	IPList []string `json:"ip_list"`
}
