package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangAccountGetAccountBalanceReq 获取账号余额接口请求
type JuGuangAccountGetAccountBalanceReq struct {
	AdvertiserId int64 `json:"advertiser_id"`
}

// ## JuGuangAccountGetAccountBalanceRes 获取账号余额返回结果
//   JuGuangRes 基础响应
//   Data: 业务数据主体，包含具体的业务响应信息
//     • FreezeBalance: 冻结余额
//     • AvailableBalance: 可用余额
//     • TodaySpend: 今日花费
//     • CompensateReturnBalance: 补偿返还余额
//     • TotalBalance: 总余额
//     • CashBalance: 现金余额
//     • ReturnBalance: 返还余额
type JuGuangAccountGetAccountBalanceRes struct {
	response.JuGuangRes
	Data struct {
		FreezeBalance           int `json:"freeze_balance"`
		AvailableBalance        int `json:"available_balance"`
		TodaySpend              int `json:"today_spend"`
		CompensateReturnBalance int `json:"compensate_return_balance"`
		TotalBalance            int `json:"total_balance"`
		CashBalance             int `json:"cash_balance"`
		ReturnBalance           int `json:"return_balance"`
	} `json:"data"`
}
