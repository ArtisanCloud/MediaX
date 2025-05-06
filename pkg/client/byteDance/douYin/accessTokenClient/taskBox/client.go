package taskBox

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/taskBox/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// DouYinTaskBoxClient 抖音任务盒子客户端
type DouYinTaskBoxClient struct {
	*kernel.BaseClient
}

// NewClient 创建抖音任务盒子客户端
func NewClient(c *kernel.BaseClient) *DouYinTaskBoxClient {
	return &DouYinTaskBoxClient{
		BaseClient: c,
	}
}

// ## QueryViolateTalentList 中介查询违规达人列表
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/query-violate-talent-list
//
// 参数：
//   ctx context.Context 请求上下文
//
// 返回值：
//   *schema.DouYinTaskBoxQueryViolateTalentListRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - ViolateTalentsURL: 违规达人excel下载链接
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) QueryViolateTalentList(ctx context.Context) (*schema.DouYinTaskBoxQueryViolateTalentListRes, error) {
	result := &schema.DouYinTaskBoxQueryViolateTalentListRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v2/taskbox/query_violate_talent_list/", nil, nil, nil, result)
	return result, err
}

// ## AgencyQueryBillLink 查询推广计划视频数据下载链接
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/agency-query-bill-link
//
// 参数：
//   ctx context.Context 请求上下文
//   billDate string 账单日期，格式为YYYYMMDD
//
// 返回值：
//   *schema.DouYinTaskBoxAgencyQueryBillLinkRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - BillLink: 推广计划视频数据下载链接
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) AgencyQueryBillLink(ctx context.Context, billDate string) (*schema.DouYinTaskBoxAgencyQueryBillLinkRes, error) {
	result := &schema.DouYinTaskBoxAgencyQueryBillLinkRes{}
	params := &object.HashMap{
		"bill_date": billDate,
	}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v2/taskbox/agency_query_bill_link/", nil, params, nil, result)
	return result, err
}

// ## AgencyQueryVideoSum 查询视频任务相关实时汇总数据
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/agency-query-video-sum
//
// 参数：
//   ctx context.Context 请求上下文
//   videoIds []int64 视频ID列表
//
// 返回值：
//   *schema.DouYinTaskBoxAgencyQueryVideoSumRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - Data: 视频汇总数据列表，包含以下字段：
//         • BillingGmvTd: 当日GMV（Gross Merchandise Volume，商品交易总额）
//         • BillingRefundGmvTd: 当日退款GMV
//         • VideoId: 视频ID
//         • TaskId: 任务ID
//         • AgentId: 代理商ID
//         • PublishTime: 发布时间（时间戳）
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) AgencyQueryVideoSum(ctx context.Context, videoIds []int64) (*schema.DouYinTaskBoxAgencyQueryVideoSumRes, error) {
	result := &schema.DouYinTaskBoxAgencyQueryVideoSumRes{}
	params := &object.HashMap{
		"video_ids": videoIds,
	}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v2/taskbox/agency_query_video_sum_data/", nil, params, nil, result)
	return result, err
}

// ## GenAgentLink 获取合作链接
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/gen-agent-link
//
// 参数：
//   ctx context.Context 请求上下文
//   data *schema.DouYinTaskBoxGenAgentLinkReq 生成链接请求参数，包含以下字段：
//     • AgentID: 团长ID
//     • AgencyTalentUID: 开发者透传的达人uid
//     • AppID: 小程序appid
//     • TaskCategory: 任务类别枚举（1=短视频任务，2=直播任务）
//     • TaskID: 任务ID
//
// 返回值：
//   *schema.DouYinTaskBoxGenAgentLinkRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - WebLink: 网页链接
//       - AppLink: 小程序链接
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) GenAgentLink(ctx context.Context, data *schema.DouYinTaskBoxGenAgentLinkReq) (*schema.DouYinTaskBoxGenAgentLinkRes, error) {
	result := &schema.DouYinTaskBoxGenAgentLinkRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v1/taskbox/gen_agent_link/", nil, data, nil, result)
	return result, err
}

// ## ChangeUserBindAgent 换绑团长
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/change-user-bind-agent
//
// 参数：
//   ctx context.Context 请求上下文
//   data *schema.DouYinTaskBoxChangeUserBindAgentReq 换绑请求参数，包含以下字段：
//     • DouyinID: 达人的抖音号
//     • NewAgentID: 达人换绑的团长ID
//     • OldAgentID: 达人当前绑定的团长ID
//
// 返回值：
//   *schema.DouYinTaskBoxChangeUserBindAgentRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) ChangeUserBindAgent(ctx context.Context, data *schema.DouYinTaskBoxChangeUserBindAgentReq) (*schema.DouYinTaskBoxChangeUserBindAgentRes, error) {
	result := &schema.DouYinTaskBoxChangeUserBindAgentRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v1/taskbox/change_user_bind_agent/", nil, data, nil, result)
	return result, err
}

// ## GetAgencyUserBindRecord 查询绑定关系
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/get-agency-user-bind-record
//
// 参数：
//   ctx context.Context 请求上下文
//   data *schema.DouYinTaskBoxGetAgencyUserBindRecordReq 查询请求参数，包含以下字段：
//     • PageNo: 页号
//     • PageSize: 页大小
//     • AgencyTalentUID: 机构侧达人uid
//     • AgentID: 团长ID
//     • DouyinID: 达人抖音号
//
// 返回值：
//   *schema.DouYinTaskBoxGetAgencyUserBindRecordRes 包含以下字段：
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - Data: 绑定关系列表，包含以下字段：
//         • AgentID: 团长ID
//         • DouyinID: 达人抖音号
//         • BindTime: 绑定时间
//         • UnbindTime: 解绑时间
//     • LogId: 日志ID，oncall时请携带此参数值
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) GetAgencyUserBindRecord(ctx context.Context, data *schema.DouYinTaskBoxGetAgencyUserBindRecordReq) (*schema.DouYinTaskBoxGetAgencyUserBindRecordRes, error) {
	result := &schema.DouYinTaskBoxGetAgencyUserBindRecordRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v1/taskbox/get_agency_user_bind_record/", nil, data, nil, result)
	return result, err
}

// ## SaveAgent 创建团长
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/save-agent
//
// 参数：
//   ctx context.Context 请求上下文
//   data *schema.DouYinTaskBoxSaveAgentReq 创建团长请求参数，包含以下字段：
//     • AgentID: 团长ID（创建团长时，agent_id不传或传0）
//     • AgentNickname: 团长的昵称，全局唯一
//
// 返回值：
//   *schema.DouYinTaskBoxSaveAgentRes 包含以下字段：
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - AgentID: 团长ID
//     • LogId: 日志ID，oncall时请携带此参数值
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) SaveAgent(ctx context.Context, data *schema.DouYinTaskBoxSaveAgentReq) (*schema.DouYinTaskBoxSaveAgentRes, error) {
	result := &schema.DouYinTaskBoxSaveAgentRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v1/taskbox/save_agent/", nil, data, nil, result)
	return result, err
}

// ## AgencyQueryTaskInfo 查询任务详情
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/agency-query-task-info
//
// 参数：
//   ctx  - 请求上下文
//   data - 查询任务详情请求参数，包含以下字段：
//     • Appid: 小程序appid
//     • PageNo: 分页编号，从1开始
//     • PageSize: 分页大小
//     • QueryParamsContent: 查询参数内容
//     • QueryParamsType: 查询参数类型（1=任务ID，2=任务名称，3=任务页面地址）
//     • TaskCategory: 任务类别（1=短视频任务，2=直播任务）
//
// 返回值：
//   *schema.DouYinTaskBoxAgencyQueryTaskInfoRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - Total: 总记录数
//       - PageCount: 总页数
//       - AppId: 小程序appid
//       - Tasks: 任务列表，包含以下字段：
//         • AnchorTitle: 主播标题
//         • Appid: 小程序appid
//         • PaymentAllocateRatio: 支付分配比例
//         • PlatformAddressApp: 小程序地址
//         • PlatformAddressWeb: 网页地址
//         • ReferMaCaptures: 媒体资源截图
//         • ReferVideoCaptures: 视频资源截图
//         • RejectReason: 拒绝原因
//         • StartPage: 任务页面地址
//         • Status: 任务状态
//         • TaskDesc: 任务描述
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) AgencyQueryTaskInfo(ctx context.Context, data *schema.DouYinTaskBoxAgencyQueryTaskInfoReq) (*schema.DouYinTaskBoxAgencyQueryTaskInfoRes, error) {
	result := &schema.DouYinTaskBoxAgencyQueryTaskInfoRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v1/taskbox/query_task_info/", nil, data, nil, result)
	return result, err
}

// ## AgencyQueryAppTask 查询小程序任务台任务ID
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/agency-query-app-task
//
// 参数：
//   ctx  - 请求上下文
//   data - 查询小程序任务请求参数，包含以下字段：
//     • Appid: 小程序appid
//     • CreateEndTime: 任务创建终止时间，秒级时间戳
//     • CreateStartTime: 任务创建起始时间，秒级时间戳
//     • TaskCategory: 任务类别枚举（1=短视频任务，2=直播任务）
//
// 返回值：
//   *schema.DouYinTaskBoxAgencyQueryAppTaskRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - TaskIds: 任务ID列表
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) AgencyQueryAppTask(ctx context.Context, data *schema.DouYinTaskBoxAgencyQueryAppTaskReq) (*schema.DouYinTaskBoxAgencyQueryAppTaskRes, error) {
	result := &schema.DouYinTaskBoxAgencyQueryAppTaskRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v1/taskbox/query_app_task_id/", nil, data, nil, result)
	return result, err
}

// ## QueryAgencyVideoSumData 查询视频汇总数据-v2
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/query-agency-video-sum-data-v2
//
// 参数：
//   ctx  - 请求上下文
//   data - 查询视频汇总数据请求参数，包含以下字段：
//     • PageNum: 页号
//     • PageSize: 页大小，不能超过10000
//     • VideoPublishEndTime: 视频发布的结束时间，秒级时间戳
//     • VideoPublishStartTime: 视频发布的起始时间，秒级时间戳
//     • AgentID: 团长ID
//     • AppID: 机构合作小程序appid
//     • DouyinID: 抖音号
//
// 返回值：
//   *schema.DouYinTaskBoxQueryAgencyVideoSumDataRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - VideoSumData: 视频汇总数据列表
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) QueryAgencyVideoSumData(ctx context.Context, data *schema.DouYinTaskBoxQueryAgencyVideoSumDataReq) (*schema.DouYinTaskBoxQueryAgencyVideoSumDataRes, error) {
	result := &schema.DouYinTaskBoxQueryAgencyVideoSumDataRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v2/taskbox/query_agency_video_sum_data/", nil, data, nil, result)
	return result, err
}

// ## QueryAgencyVideoDailyData 查询视频每日明细数据-v2
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/query-agency-video-daily-data-v2
//
// 参数：
//   ctx  - 请求上下文
//   data - 查询视频每日明细数据请求参数，包含以下字段：
//     • BillingDate: 计费日期
//     • PageNum: 页号
//     • PageSize: 页大小，最大不能超过10000
//     • VideoPublishEndTime: 最大的视频发布时间戳
//     • VideoPublishStartTime: 视频发布的起始时间戳
//     • AgentID: 抖音侧团长ID
//     • AppID: 撮合中介合作的小程序appid
//     • DouyinID: 达人抖音号
//
// 返回值：
//   *schema.DouYinTaskBoxQueryAgencyVideoDailyDataRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - VideoDailyData: 视频每日明细数据列表
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) QueryAgencyVideoDailyData(ctx context.Context, data *schema.DouYinTaskBoxQueryAgencyVideoDailyDataReq) (*schema.DouYinTaskBoxQueryAgencyVideoDailyDataRes, error) {
	result := &schema.DouYinTaskBoxQueryAgencyVideoDailyDataRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v2/taskbox/query_agency_video_daily_data/", nil, data, nil, result)
	return result, err
}

// ## QueryTaskVideoStatus 查询视频状态-v2
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/taskbox/query-task-video-status-v2
//
// 参数：
//   ctx  - 请求上下文
//   data - 查询视频状态请求参数，包含以下字段：
//     • PageNum: 页号
//     • PageSize: 页大小，最大不能超过100
//     • VideoPublishEndTime: 最大的视频发布时间戳
//     • VideoPublishStartTime: 视频发布的起始时间戳
//     • AgentID: 抖音侧团长ID
//     • AppID: 撮合中介合作的小程序appid
//     • DouyinID: 达人抖音号
//
// 返回值：
//   *schema.DouYinTaskBoxQueryTaskVideoStatusRes 包含以下字段：
//     • LogId: 日志ID，oncall时请携带此参数值
//     • Data: 业务数据主体，包含具体的业务响应信息
//       - VideoStatusList: 视频状态列表
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码
//   error 调用过程中遇到的错误（如有）
func (c *DouYinTaskBoxClient) QueryTaskVideoStatus(ctx context.Context, data *schema.DouYinTaskBoxQueryTaskVideoStatusReq) (*schema.DouYinTaskBoxQueryTaskVideoStatusRes, error) {
	result := &schema.DouYinTaskBoxQueryTaskVideoStatusRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/match/v2/taskbox/query_task_video_status/", nil, data, nil, result)
	return result, err
}
