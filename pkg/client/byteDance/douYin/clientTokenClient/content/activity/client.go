package activity

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/activity/schema"
)

type DouYinContentActivityClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinContentActivityClient {
	return &DouYinContentActivityClient{
		BaseClient: c,
	}
}

// ## Create 经营任务创建活动
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/business-task/mission-creation
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • StartTime: 活动开始时间，秒级时间戳
//	       • EndTime: 活动结束时间，秒级时间戳
//	       • CreateBusinessTaskInfoList: 创建任务信息列表
//	       • NotBcCheck: 是否不进行BC检查
//	       • ActivityName: 活动名称
//
// 返回值：
//
//	*schema.DouYinContentActivityCreateRes 包含以下字段：
//	  • ActivityId: 活动ID
//	  • BusinessTaskIdList: 业务任务ID列表
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentActivityClient) Create(ctx context.Context, data *schema.DouYinContentActivityCreateReq) (*schema.DouYinContentActivityCreateRes, error) {
	result := &schema.DouYinContentActivityCreateRes{}

	_, err := c.BaseClient.HttpPost(ctx, "/dy_open_api/apps/v3/activity/create/", nil, data, nil, result)

	return result, err

}

// ## Modify 经营任务修改活动
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/business-task/business-task-update
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • StartTime: 活动开始时间，秒级时间戳
//	       • ActivityId: 活动ID
//	       • ModifyBusinessTaskInfoList: 修改的业务任务信息列表
//	       • NotBcCheck: 是否不进行BC检查
//	       • ActivityName: 活动名称
//	       • EndTime: 活动结束时间，秒级时间戳
//	       • AddBusinessTaskInfoList: 新增的业务任务信息列表
//	       • DelBusinessTaskIdList: 删除的业务任务ID列表
//
// 返回值：
//
//	*schema.DouYinContentActivityCreateRes 包含以下字段：
//	  • ModifyBusinessTaskIdList: 修改的业务任务ID列表
//	  • ActivityId: 活动ID
//	  • AddBusinessTaskIdList: 新增的业务任务ID列表
//	  • DelBusinessTaskIdList: 删除的业务任务ID列表
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentActivityClient) Modify(ctx context.Context, data *schema.DouYinContentActivityCreateReq) (*schema.DouYinContentActivityCreateRes, error) {
	result := &schema.DouYinContentActivityCreateRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/dy_open_api/apps/v3/activity/modify/", nil, data, nil, result)
	return result, err
}

// ## QueryInfo 经营任务查询活动信息
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/business-task/business-task-info
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • ActivityId: 活动ID
//	       • TaskIdList: 任务ID列表
//
// 返回值：
//
//	*schema.DouYinContentActivityQueryInfoRes 包含以下字段：
//	  • BusinessTaskInfoList: 业务任务信息列表
//	  • ActivityInfo: 活动信息
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentActivityClient) QueryInfo(ctx context.Context, data *schema.DouYinContentActivityQueryInfoReq) (*schema.DouYinContentActivityQueryInfoRes, error) {
	result := &schema.DouYinContentActivityQueryInfoRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/dy_open_api/apps/v3/activity/query_info/", nil, data, nil, result)
	return result, err
}

// ## QueryUserStatus 经营任务查询用户是否完成
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/business-task/business-task-query
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • ActivityId: 活动ID
//	       • TargetOpenId: 目标用户OpenId
//	       • AlliedId: 联盟ID
//	       • TaskIdList: 任务ID列表
//
// 返回值：
//
//	*schema.DouYinContentActivityQueryUserStatusRes 包含以下字段：
//	  • TaskCompleteStatusMap: 任务完成状态映射
//	  • PostingVideoBindCompleteInfoMap: 发布视频绑定完成信息映射
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentActivityClient) QueryUserStatus(ctx context.Context, data *schema.DouYinContentActivityQueryUserStatusReq) (*schema.DouYinContentActivityQueryUserStatusRes, error) {
	result := &schema.DouYinContentActivityQueryUserStatusRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/dy_open_api/apps/v3/activity/query_activity_user_completion_status/", nil, data, nil, result)
	return result, err
}
