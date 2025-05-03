package message

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/message/schema"
)

// DouYinIMMessageClient 抖音IM消息客户端
// 提供与抖音IM消息相关的接口封装
type DouYinIMMessageClient struct {
	*kernel.BaseClient
}

// NewClient 初始化并返回一个新的 DouYinIMMessageClient 实例
func NewClient(c *kernel.BaseClient) *DouYinIMMessageClient {
	return &DouYinIMMessageClient{
		BaseClient: c,
	}
}

// ## Send 发送私信消息
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/private-message/send-msg#%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E
//
// 参数：
//   ctx  - 请求上下文
//   accessToken - 访问令牌
//   req - 发送消息请求参数
//
// 返回值：
//   *schema.DouYinIMMessageSendRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//     • MsgId: 消息ID
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMMessageClient) Send(ctx context.Context, accessToken string, req *schema.DouYinIMMessageSendReq) (*schema.DouYinIMMessageSendRes, error) {
	result := &schema.DouYinIMMessageSendRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/api/douyin/v1/im/message/send/", nil, req, nil, result)
	return result, err
}

// ## QueryAuthorizeUserList 查询授权主动私信用户
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/private-message/query-private-auth-user
//
// 参数：
//   ctx  - 请求上下文
//   req - 查询授权用户列表请求参数
//
// 返回值：
//   *schema.DouYinIMMessageQueryAuthorizeUserListRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data:
//         - AuthUserList: 授权用户列表
//         - ErrorCode: 错误码，0 表示成功，其他为失败
//         - Description: 错误描述或状态说明
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMMessageClient) QueryAuthorizeUserList(ctx context.Context, req *schema.DouYinIMMessageQueryAuthorizeUserListReq) (*schema.DouYinIMMessageQueryAuthorizeUserListRes, error) {
	result := &schema.DouYinIMMessageQueryAuthorizeUserListRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/im/authorize/user_list/", nil, req, nil, result)
	return result, err
}

// ## Recall 私信消息撤回
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/private-message/recall-msg
//
// 参数：
//   ctx  - 请求上下文
//   req - 撤回消息请求参数
//
// 返回值：
//   *schema.DouYinIMMessageRecallRes 包含以下字段：
//     • ErrMsg: 错误信息
//     • ErrNo: 错误码，0-成功，非0-失败
//     • LogId: 日志ID，用于问题定位
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMMessageClient) Recall(ctx context.Context, req *schema.DouYinIMMessageRecallReq) (*schema.DouYinIMMessageRecallRes, error) {
	result := &schema.DouYinIMMessageRecallRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/im/recall/msg/", nil, req, nil, result)
	return result, err

}
