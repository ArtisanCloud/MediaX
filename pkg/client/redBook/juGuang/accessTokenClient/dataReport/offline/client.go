package offline

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient/dataReport/offline/schema"
)

// JuGuangDataReportOfflineClient 表示巨量引擎离线数据报表 API 客户端
type JuGuangDataReportOfflineClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的 JuGuangDataReportOfflineClient 实例
func NewClient(c *kernel.BaseClient) *JuGuangDataReportOfflineClient {
	return &JuGuangDataReportOfflineClient{
		BaseClient: c,
	}
}

// ## GroupReprot 获取人群包报表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=4065
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • startDate: 开始时间 YYYY-MM-DD（必填）
//	  • endDate: 结束时间 YYYY-MM-DD（必填）
//	  • timeUnit: 时间维度（必填，DAY：分天，SUMMARY：汇总）
//	  • columns: 请求列（必填，必须包含groupName, groupId）
//	  • pageSize: 页面大小（可选，默认20）
//	  • pageNum: 页面index（可选，默认1）
//	  • filters: 过滤选项（可选）
//	  • sorts: 排序选项（可选）
//	  • needTotal: 是否需要综合数据（可选）
//	  • reportType: 查询维度（可选，默认为USER_GROUP）
//
// 返回值：
//
//	*schema.JuGuangDataReportOfflineGroupReprotRes 包含以下字段：
//	  • page: 分页信息
//	  • aggregationData: 综合数据
//	  • dataList: 分时数据
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangDataReportOfflineClient) GroupReprot(ctx context.Context, data *schema.JuGuangDataReportOfflineGroupReprotReq) (*schema.JuGuangDataReportOfflineGroupReprotRes, error) {
	result := &schema.JuGuangDataReportOfflineGroupReprotRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/idea/group_report", nil, data, nil, result)
	return result, err
}

// ## SpuLevel 获取SPU层级离线报表数据
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3835
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • start_date: 开始时间 YYYY-MM-DD（必填）
//	  • end_date: 结束时间 YYYY-MM-DD（必填）
//	  • time_unit: 时间维度（可选，DAY：分天，SUMMARY：汇总，默认DAY）
//	  • sort_column: 排序字段（可选，见文档附录column字段）
//	  • sort: 排序方式（可选，asc：升序，desc：降序）
//	  • page_num: 页码（可选，默认1）
//	  • page_size: 每页数量（可选，默认20，最大500）
//
// 返回值：
//
//	*schema.JuGuangDataReportOfflineSpuLevelRes 包含以下字段：
//	  • total_count: 总条数
//	  • data_list: 详细数据列表
//	  • aggregation_data: 汇总数据
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangDataReportOfflineClient) SpuLevel(ctx context.Context, data *schema.JuGuangDataReportOfflineSpuLevelReq) (*schema.JuGuangDataReportOfflineSpuLevelRes, error) {
	result := &schema.JuGuangDataReportOfflineSpuLevelRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/data/report/offline/spu", nil, data, nil, result)
	return result, err
}

// ## SeriesLevel 获取系列层级离线数据报表
//
// 接口文档参考：
// https://ad-market.xiaohongshu.com/docs-center?bizType=943&articleId=3822
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • advertiser_id: 广告主ID（必填）
//	  • start_date: 开始时间 YYYY-MM-DD（必填）
//	  • end_date: 结束时间 YYYY-MM-DD（必填）
//	  • time_unit: 时间维度（可选，DAY：分天，SUMMARY：汇总，默认DAY）
//	  • sort_column: 排序字段（可选，见文档附录column字段）
//	  • sort: 排序方式（可选，asc：升序，desc：降序）
//	  • page_num: 页码（可选，默认1）
//	  • page_size: 每页数量（可选，默认20，最大500）
//
// 返回值：
//
//	*schema.JuGuangDataReportOfflineSeriesLevelRes 包含以下字段：
//	  • total_count: 总条数
//	  • data_list: 详细数据列表
//	  • aggregation_data: 汇总数据
//	error 调用过程中遇到的错误（如有）
func (c *JuGuangDataReportOfflineClient) SeriesLevel(ctx context.Context, data *schema.JuGuangDataReportOfflineSeriesLevelReq) (*schema.JuGuangDataReportOfflineSeriesLevelRes, error) {
	result := &schema.JuGuangDataReportOfflineSeriesLevelRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/open/jg/data/report/offline/series", nil, data, nil, result)
	return result, err
}
