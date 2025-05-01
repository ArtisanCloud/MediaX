package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinContentVideoShareResultReq 表示获取抖音视频分享结果的请求结构体。
// 目前该请求结构体为空，可根据实际需求添加字段。
type DouYinContentVideoShareResultReq struct {
	// Count 每页数量，必填字段。
	// 示例值: 10
	Count int64 `json:"count" validate:"required"`
	// Keyword 查询关键字，例如美食，必填字段。
	// 示例值: 美食
	Keyword string `json:"keyword" validate:"required"`
	// City 查询城市，例如上海、北京。
	// 示例值: 北京
	City string `json:"city"`
	// Cursor 分页游标, 第一页请求cursor是0, response中会返回下一页请求用到的cursor,
	// 同时response还会返回has_more来表明是否有更多的数据。
	// 示例值: 0
	Cursor int64 `json:"cursor"`
}

// Poi 表示地理位置信息结构体，包含地点的详细信息。
type Poi struct {
	// Address 地点的具体地址。
	Address string `json:"address"`
	// City 地点所在的城市名称。
	City string `json:"city"`
	// CityCode 地点所在城市的编码。
	CityCode string `json:"city_code"`
	// Country 地点所在的国家名称。
	Country string `json:"country"`
	// CountryCode 地点所在国家的编码。
	CountryCode string `json:"country_code"`
	// District 地点所在的行政区名称。
	District string `json:"district"`
	// Location 地点的地理位置坐标。
	Location string `json:"location"`
	// PoiId 地点的唯一标识 ID。
	PoiId string `json:"poi_id"`
	// PoiName 地点的名称。
	PoiName string `json:"poi_name"`
	// Province 地点所在的省份名称。
	// Pois 与视频分享结果相关的地点信息列表。
	// DouYinRes 通用的抖音响应信息，包含错误码、描述等。
	// Data 业务数据主体，包含具体的业务响应信息和地点列表。
	// Extra 包含通用的响应扩展字段，如日志 ID、当前时间戳、错误码等。
	// DouYinContentVideoShareResultRes 表示获取抖音视频分享结果的响应结构体。
	Province string `json:"province"`
}

type DouYinContentVideoShareResultRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`

	// Data 业务数据主体。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		Pois []Poi `json:"pois"`
	} `json:"data"`
}
