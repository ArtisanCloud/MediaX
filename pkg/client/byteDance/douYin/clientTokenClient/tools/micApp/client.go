package micApp

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/tools/micApp/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type DouYinToolMicAppClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinToolMicAppClient {
	return &DouYinToolMicAppClient{
		BaseClient: c,
	}
}

// ## IsLegal 检查小程序合法性
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/tools-ability/mini-app-interface
//
// 参数：
//   ctx      - 请求上下文
//   micappId - 小程序ID
//
// 返回值：
//   *schema.DouYinToolMicAppIsLegalRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//       • IsLegal: 小程序是否合法
//   error 调用过程中遇到的错误（如有）
func (c *DouYinToolMicAppClient) IsLegal(ctx context.Context, micappId string) (*schema.DouYinToolMicAppIsLegalRes, error) {
	result := &schema.DouYinToolMicAppIsLegalRes{}
	params := &object.StringMap{
		"micapp_id": micappId,
	}
	_, err := c.BaseClient.HttpGet(ctx, "/devtool/micapp/is_legal/", params, nil, result)
	return result, err
}
