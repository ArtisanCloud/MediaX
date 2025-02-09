package response

import "github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/core/response"

type GetCallBackIPRes struct {
	response.OfficialAccountRes

	IPList []string `json:"ip_list"`
}
