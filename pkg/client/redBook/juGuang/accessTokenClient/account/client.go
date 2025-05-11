package account

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient/account/schema"
)

// JuGuangAccountClient 账号管理
type JuGuangAccountClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个聚光管理客户端实例
func NewClient(c *kernel.BaseClient) *JuGuangAccountClient {
	return &JuGuangAccountClient{
		BaseClient: c,
	}
}

// ## GetAccountBalance 获取账号余额
//
// 接口文档参考：
// https://ad.xiaohongshu.com/openApiDoc?uba_pre=18.target_package..1712048419713&uba_ppre=18.aurora_asset_manage..1712048416043&uba_index=11&articleId=3043
//
// 参数：
//   ctx  - 请求上下文
//   data - 请求参数，包含以下字段：
//     • AdvertiserId: 广告主ID
//
// 返回值：
//   *schema.JuGuangAccountGetAccountBalanceRes 包含以下字段：
// 	  • JuGuangRes 基础响应
//    • Data: 业务数据主体，包含具体的业务响应信息
//     		• FreezeBalance: 冻结余额
//     		• AvailableBalance: 可用余额
//     		• TodaySpend: 今日花费
//     		• CompensateReturnBalance: 补偿返还余额
//     		• TotalBalance: 总余额
//     		• CashBalance: 现金余额
//     		• ReturnBalance: 返还余额
func (c *JuGuangAccountClient) GetAccountBalance(ctx context.Context, data *schema.JuGuangAccountGetAccountBalanceReq) (*schema.JuGuangAccountGetAccountBalanceRes, error) {
	res := &schema.JuGuangAccountGetAccountBalanceRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/account/balance/info", nil, data, nil, res)
	return res, err
}
