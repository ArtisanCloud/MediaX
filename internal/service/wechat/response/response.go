package response

import "github.com/ArtisanCloud/MediaX/internal/service/wechat/core/response"

type GetCallBackIPRes struct {
	response.OfficialAccountRes

	IPList []string `json:"ip_list"`
}
