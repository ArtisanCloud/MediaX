package task

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/task/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type DouYinContentTaskClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinContentTaskClient {
	return &DouYinContentTaskClient{
		BaseClient: c,
	}
}

// ## CreatePost 创建任务
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/posting-task/create-posting-task
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • EndTime: 任务结束时间，秒级时间戳，必填
//	       • StartTime: 任务开始时间，秒级时间戳，必填
//	       • TaskCondition: 任务条件，必填
//	       • TaskName: 任务名称，长度不超过50个字符，必填
//
// 返回值：
//
//	*schema.DouYinContentCreatePostRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - TaskId: 任务ID
//	      - TaskStatus: 任务状态
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//
//	error 调用过程中遇到的错误（如有）
func (c *DouYinContentTaskClient) CreatePost(ctx context.Context, data *schema.DouYinContentCreatePostReq) (*schema.DouYinContentCreatePostRes, error) {
	result := &schema.DouYinContentCreatePostRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpPost(ctx, "/task/posting/create/", nil, params, nil, result)
	return result, err
}
