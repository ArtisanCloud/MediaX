package schema

// JuGuangNoteGetGetPOIListReq 表示获取门店列表的请求参数
type JuGuangNoteGetGetPOIListReq struct {
	AdvertiserID int64 `json:"advertiser_id"` // 广告主id
}

// JuGuangNoteGetGetPOIListRes 表示获取门店列表的响应
type JuGuangNoteGetGetPOIListRes struct {
	Code    int    `json:"code"`    // 返回码
	Msg     string `json:"msg"`     // 返回信息
	Success bool   `json:"success"` // 接口是否成功
	Data    struct {
		POIInfoList []POIInfo `json:"poi_info_list"` // 门店列表
		Page        PageInfo  `json:"page"`          // 分页信息
	} `json:"data"`
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
