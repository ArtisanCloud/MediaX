package task

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/content/task/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// DouYinContentTaskClient 抖音内容任务客户端
// 提供与抖音内容任务相关的接口封装
type DouYinContentTaskClient struct {
	*kernel.BaseClient
}

// NewClient 初始化并返回一个新的 DouYinContentTaskClient 实例
// 参数：
//   c - 基础客户端实例
// 返回值：
//   *DouYinContentTaskClient - 新的抖音内容任务客户端实例
func NewClient(c *kernel.BaseClient) *DouYinContentTaskClient {
	return &DouYinContentTaskClient{
		BaseClient: c,
	}
}

// BindVideo 绑定视频到任务
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/posting-task/bind-video
//
// 参数：
//   ctx  - 请求上下文
//   data - 请求参数，包含：
//          • TaskId: 任务ID，创建任务之后获取的任务ID，必填
//          • VideoId: 视频ID，必填
// 返回值：
//   *schema.DouYinContentBindVideoRes - 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error - 调用过程中遇到的错误（如有）
func (c *DouYinContentTaskClient) BindVideo(ctx context.Context, data *schema.DouYinContentBindVideoReq) (*schema.DouYinContentBindVideoRes, error) {
	result := &schema.DouYinContentBindVideoRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpPost(ctx, "/task/posting/bind_video/", nil, params, nil, result)
	return result, err
}
