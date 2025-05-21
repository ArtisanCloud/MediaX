package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// Page 表示分页信息
//   - PageIndex: 当前页码，从1开始
//   - PageSize: 每页大小，最大100
type Page struct {
	PageIndex *int `json:"page_index,omitempty"`
	PageSize  *int `json:"page_size,omitempty"`
}

// EventList 表示事件列表信息
//   - EventID: 事件ID
//   - EventType: 事件类型
//   - EventStatus: 事件状态
type EventList struct {
	EventID     int64 `json:"event_id"`
	EventType   int   `json:"event_type"`
	EventStatus int   `json:"event_status"`
}

// EventAssetDto 表示事件资产信息
//   - EventAssetID: 事件资产ID
//   - EventAssetName: 事件资产名称
//   - Status: 状态
//   - EventList: 关联的事件列表
type EventAssetDto struct {
	EventAssetID   int64       `json:"event_asset_id"`
	EventAssetName string      `json:"event_asset_name"`
	Status         int         `json:"status"`
	EventList      []EventList `json:"event_list"`
}

// JuGuangNoteGetAssetInfoReq 表示获取事件资产信息的请求结构体
//   - AdvertiserID: 广告主ID
//   - Page: 分页信息
type JuGuangNoteGetAssetInfoReq struct {
	AdvertiserID int64 `json:"advertiser_id"`
	Page         Page  `json:"page"`
}

// GetAssetInfoData 表示事件资产信息数据
//   - EventAssetDtos: 事件资产列表
//   - Page: 分页信息
type GetAssetInfoData struct {
	EventAssetDtos []EventAssetDto `json:"event_assert_dtos"`
	Page           Page            `json:"page"`
}

// JuGuangNoteGetAssetInfoRes 表示获取事件资产信息的响应结构体
//   - RedBookAccessTokenRes: 基础响应结构
//   - Data: 事件资产信息数据
type JuGuangNoteGetAssetInfoRes struct {
	response.RedBookAccessTokenRes
	Data GetAssetInfoData `json:"data"`
}
