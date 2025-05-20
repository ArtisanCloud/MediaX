package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteGetPOIListReq 表示获取门店列表的请求参数
type JuGuangNoteGetPOIListReq struct {
	AdvertiserID int64 `json:"advertiser_id"` // 广告主id
}

type GetPOIListData struct {
	POIInfoList []POIInfo `json:"poi_info_list"` // 门店列表
	Page        PageInfo  `json:"page"`          // 分页信息
}

// JuGuangNoteGetPOIListRes 表示获取门店列表的响应
type JuGuangNoteGetPOIListRes struct {
	response.RedBookAccessTokenRes
	Data GetPOIListData `json:"data"`
}

// POIInfo 表示门店信息
type POIInfo struct {
	POIID   string `json:"poi_id"`   // 门店id
	POIName string `json:"poi_name"` // 门店名称
}

// PageInfo 表示分页信息
type PageInfo struct {
	Total int `json:"total"` // 总数
}
