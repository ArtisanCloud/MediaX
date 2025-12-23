package sandbox

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/tools/sandbox/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type DouYinToolSandboxClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinToolSandboxClient {
	return &DouYinToolSandboxClient{
		BaseClient: c,
	}
}

// ## MockWebhookEvent 模拟webhook事件
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/tools-ability/sandbox-management/mock-webhook-event
//
// 参数：
//
//	ctx       - 请求上下文
//	eventType - 事件类型
//
// 返回值：
//
//	*schema.DouYinToolSandboxMockWebhookEventRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data: 业务数据主体，包含具体的业务响应信息
//	error 调用过程中遇到的错误（如有）
func (c *DouYinToolSandboxClient) MockWebhookEvent(ctx context.Context, eventType string) (*schema.DouYinToolSandboxMockWebhookEventRes, error) {
	result := &schema.DouYinToolSandboxMockWebhookEventRes{}
	params := &object.HashMap{
		"event_type": eventType,
	}
	_, err := c.BaseClient.HttpPost(ctx, "/sandbox/webhook/event/send/", nil, params, nil, result)
	return result, err
}
