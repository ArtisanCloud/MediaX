package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinMarketServiceDeletePurchaseInfoReq 删除购买信息请求参数
type DouYinMarketServiceDeletePurchaseInfoReq struct {
	// OpenID 查询订阅信息目标用户的唯一标识
	OpenID string `json:"open_id"`
	// OutTradeNo 外部开发者单号id，用于唯一标识同一服务同一规格下要删除的订购数据
	OutTradeNo string `json:"out_trade_no"`
	// ServiceID 服务id，服务的唯一标识
	ServiceID string `json:"service_id"`
	// ServiceModeID 用户订购的服务规格id
	ServiceModeID string `json:"service_mode_id"`
}

// DouYinMarketServiceDeletePurchaseInfoRes 删除购买信息响应结果
type DouYinMarketServiceDeletePurchaseInfoRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
