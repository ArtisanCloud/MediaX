package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangDataReportRealtimeAccountLevelReq 表示获取账户层级实时报表数据的请求结构体
type JuGuangDataReportRealtimeAccountLevelReq struct {
	AdvertiserID   int64  `json:"advertiser_id"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	NeedHourlyData bool   `json:"need_hourly_data,omitempty"`
}

// JuGuangDataReportRealtimeAccountLevelRes 表示获取账户层级实时报表数据的响应结构体
type JuGuangDataReportRealtimeAccountLevelRes struct {
	response.RedBookAccessTokenRes
	Data       DataReportDTO   `json:"data"`
	HourlyData []DataReportDTO `json:"hourly_data,omitempty"`
}
