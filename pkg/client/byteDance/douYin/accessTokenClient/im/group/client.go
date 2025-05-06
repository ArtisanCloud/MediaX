package group

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/group/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// DouYinIMGroupClient 抖音IM群组客户端
// 提供与抖音IM群组相关的接口封装
type DouYinIMGroupClient struct {
	*kernel.BaseClient
}

// NewClient 初始化并返回一个新的 DouYinIMGroupClient 实例
func NewClient(c *kernel.BaseClient) *DouYinIMGroupClient {
	return &DouYinIMGroupClient{
		BaseClient: c,
	}
}

// ## Create 创建粉丝群
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/create-fans-group
//
// 参数：
//   ctx  - 请求上下文
//   data - 创建粉丝群请求参数
//     • AvatarUri: 群头像链接
//     • Description: 群简介
//     • GroupName: 群名称
//     • ActiveFans: 进群门槛-粉丝活跃度
//     • AllowInvite: 允许群成员邀请朋友
//     • FansLimit: 进群门槛-粉丝团等级
//     • GroupType: 群类型
//     • ItemAutoSync: 作品同步
//     • LiveAutoSync: 直播同步
//     • OpenAuditSwitch: 开启进群审批
//     • RelationType: 进群门槛-关注条件
//     • ShowAtProfile: 展示到个人主页
//
// 返回值：
//   *schema.DouYinIMGroupCreateRes 包含以下字段：
//     • GroupId: 群组ID，创建成功后返回的群组ID
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) Create(ctx context.Context, data *schema.DouYinIMGroupCreateReq) (*schema.DouYinIMGroupCreateRes, error) {
	result := &schema.DouYinIMGroupCreateRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/im/group/fans/create/", nil, data, nil, result)
	return result, err
}

// ## SetUserEnterAudit 变更用户入群申请状态
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/update-user-audit-status
//
// 参数：
//   ctx  - 请求上下文
//   data - 变更用户入群申请状态请求参数
//     • ApplyId: 入群申请ID，来源于查询群主所在群的用户入群申请状态
//     • Status: 目标申请状态，2-通过，3-拒绝
//
// 返回值：
//   *schema.DouYinIMGroupSetUserEnterAuditRes 包含以下字段：
//     • Success: 操作是否成功
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) SetUserEnterAudit(ctx context.Context, data *schema.DouYinIMGroupSetUserEnterAuditReq) (*schema.DouYinIMGroupSetUserEnterAuditRes, error) {
	result := &schema.DouYinIMGroupSetUserEnterAuditRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/im/group/enter/audit/set/", nil, data, nil, result)
	return result, err
}

// ## FanList 查询粉丝群
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/query-fans-group
//
// 参数：
//   ctx  - 请求上下文
//
// 返回值：
//   *schema.DouYinIMGroupFanListRes 包含以下字段：
//     • GroupList: 粉丝群列表，包含群组ID、群名称、群描述等信息
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) FanList(ctx context.Context) (*schema.DouYinIMGroupFanListRes, error) {
	result := &schema.DouYinIMGroupFanListRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/im/group/fans/list/", nil, nil, result)
	return result, err
}

// ## SetGroupSetting 设置进群问候语&群公告
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/set-group-setting
//
// 参数：
//   ctx  - 请求上下文
//   data - 设置群组配置请求参数
//     • GroupSettingType: 群管理配置类型
//     • MsgList: 消息列表，支持文本、小程序卡片和小程序券
//
// 返回值：
//   *schema.DouYinIMGroupSetGroupSettingRes 包含以下字段：
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码，0-成功，非0-失败
//     • LogId: 日志ID，用于问题定位
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) SetGroupSetting(ctx context.Context, data *schema.DouYinIMGroupSetGroupSettingReq) (*schema.DouYinIMGroupSetGroupSettingRes, error) {
	result := &schema.DouYinIMGroupSetGroupSettingRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/im/group/setting/set/", nil, data, nil, result)
	return result, err
}

// ## DisableGroupSetting 取消进群问候语&群公告配置
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/cancel-group-setting
//
// 参数：
//   ctx  - 请求上下文
//   data - 取消群组配置请求参数
//     • GroupId: 粉丝群ID
//     • GroupSettingType: 群管理配置类型
//
// 返回值：
//   *schema.DouYinIMGroupDisableGroupSettingRes 包含以下字段：
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码，0-成功，非0-失败
//     • LogId: 日志ID，用于问题定位
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) DisableGroupSetting(ctx context.Context, data *schema.DouYinIMGroupDisableGroupSettingReq) (*schema.DouYinIMGroupDisableGroupSettingRes, error) {
	result := &schema.DouYinIMGroupDisableGroupSettingRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/im/group/setting/disable/", nil, data, nil, result)
	return result, err

}

// ## Count 查询用户剩余建群额度
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/group-fans-nums
//
// 参数：
//   ctx  - 请求上下文
//
// 返回值：
//   *schema.DouYinIMGroupCountRes 包含以下字段：
//     • LeftCount: 剩余建群额度
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) Count(ctx context.Context) (*schema.DouYinIMGroupCountRes, error) {
	result := &schema.DouYinIMGroupCountRes{}
	_, err := c.BaseClient.HttpGet(ctx, "/im/group/fans/count/", nil, nil, result)
	return result, err
}

// ## GetUserEnterAudit 查询群主所在群的用户入群申请状态
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/group-fans-audit-get
//
// 参数：
//   ctx  - 请求上下文
//   data - 查询用户入群申请状态请求参数
//     • Count: 每页的数量，最大不超过50，最小不低于1
//     • OpenId: 用户open_id，通过/oauth/access_token/获取，用户唯一标志
//     • Cursor: 分页游标，第一页请求cursor是0
//
// 返回值：
//   *schema.DouYinIMGroupGetUserEnterAuditRes 包含以下字段：
//     • ApplyList: 入群申请列表，包含申请ID、申请状态、创建时间等信息
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) GetUserEnterAudit(ctx context.Context, data *schema.DouYinIMGroupGetUserEnterAuditReq) (*schema.DouYinIMGroupGetUserEnterAuditRes, error) {
	result := &schema.DouYinIMGroupGetUserEnterAuditRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/im/group/enter/audit/get/", params, nil, result)
	return result, err
}

// ## Send 发送群消息
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/send-group-message
//
// 参数：
//   ctx  - 请求上下文
//   data - 发送群消息请求参数
//     • GroupId: 群ID，来源于群消息webhook事件或查询群信息接口
//     • Content: 消息体，支持文本、图片、视频、留资卡片、小程序引导卡片等
//     • GroupToken: 群Token
//
// 返回值：
//   *schema.DouYinIMGroupSendRes 包含以下字段：
//     • MsgId: 消息ID，发送成功后返回的消息ID
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) Send(ctx context.Context, data *schema.DouYinIMGroupSendReq) (*schema.DouYinIMGroupSendRes, error) {
	result := &schema.DouYinIMGroupSendRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/im/send/msg/group/", nil, data, nil, result)
	return result, err
}

// ## Recall 撤回群消息
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/group-chat/join-group/recall-group-message
//
// 参数：
//   ctx  - 请求上下文
//   data - 撤回群消息请求参数
//     • MsgId: 消息ID，来源于发送群消息接口返回的MsgId
//     • GroupToken: 群Token
//
// 返回值：
//   *schema.DouYinIMGroupRecallRes 包含以下字段：
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码，0-成功，非0-失败
//     • LogId: 日志ID，用于问题定位
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMGroupClient) Recall(ctx context.Context, data *schema.DouYinIMGroupRecallReq) (*schema.DouYinIMGroupRecallRes, error) {
	result := &schema.DouYinIMGroupRecallRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/im/recall/msg/", nil, data, nil, result)
	return result, err
}
