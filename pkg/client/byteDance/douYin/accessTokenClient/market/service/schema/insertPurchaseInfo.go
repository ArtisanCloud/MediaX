package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinMarketServiceInsertPurchaseInfoReq 插入购买信息请求参数
type DouYinMarketServiceInsertPurchaseInfoReq struct {
	// OpenID 查询订阅信息目标用户的唯一标识
	OpenID string `json:"open_id"`
	// OutTradeNo 外部开发者单号id，用于唯一标识同一服务同一规格当前导入的订购数据
	OutTradeNo string `json:"out_trade_no"`
	// PeriodType 服务规格的周期类型
	PeriodType int `json:"period_type"`
	// PurchaseTime 服务订购时间的毫秒级时间戳
	PurchaseTime int64 `json:"purchase_time"`
	// ServiceID 服务id，服务的唯一标识
	ServiceID string `json:"service_id"`
	// ServiceModeID 用户订购的服务规格id
	ServiceModeID string `json:"service_mode_id"`
	// Duration 订购时长值，适用于时间类型的服务规格
	Duration int64 `json:"duration,omitempty"`
	// UsageTimes 用户在服务商侧购买的服务有效使用次数/条数
	UsageTimes int64 `json:"usage_times,omitempty"`
}

// DouYinMarketServiceInsertPurchaseInfoRes 插入购买信息响应结果
type DouYinMarketServiceInsertPurchaseInfoRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
