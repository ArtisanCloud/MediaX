package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinMarketServiceDecrUserRemainTimesReq 扣除用户剩余次数请求参数
type DouYinMarketServiceDecrUserRemainTimesReq struct {
	// Count 要扣除的次数/条数，要求大于0
	Count int32 `json:"count"`
	// IsTestEnv 是否为测试环境
	IsTestEnv bool `json:"is_test_env"`
	// OpenID 查询订阅信息目标用户的唯一标识，
	// *** 注意，可以根据Token对象，获取到授权的openId ***
	OpenID string `json:"open_id"`
	// OutTradeNo 外部开发者单号id，用于保证当前扣除操作的幂等
	OutTradeNo string `json:"out_trade_no"`
	// ServiceID 服务id，服务的唯一标识
	ServiceID string `json:"service_id"`
	// ServiceModeID 服务规格id，由服务商在创建服务时指定
	ServiceModeID string `json:"service_mode_id"`
}

// DouYinMarketServiceDecrUserRemainTimesRes 扣除用户剩余次数响应参数
type DouYinMarketServiceDecrUserRemainTimesRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
