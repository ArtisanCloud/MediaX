package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// PurchaseInfo 购买信息
type PurchaseInfo struct {
	ServiceId          string `json:"service_id"`
	ServiceName        string `json:"service_name"`
	ServiceStatus      int    `json:"service_status"`
	SpecificationType  int    `json:"specification_type"`
	SpecificationTitle string `json:"specification_title"`
	ServiceModeId      string `json:"service_mode_id"`
	RemainServiceTimes int    `json:"remain_service_times,omitempty"`
	EffectiveTime      int64  `json:"effective_time,omitempty"`
	ExpireTime         int64  `json:"expire_time,omitempty"`
}

// DouYinMarketServiceGetUserServicePurchaseListRes 获取用户服务购买列表
type DouYinMarketServiceGetUserServicePurchaseListRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		// PurchaseInfoList 购买信息列表
		PurchaseInfoList []PurchaseInfo `json:"purchase_info_list"`
	} `json:"data"`
}
