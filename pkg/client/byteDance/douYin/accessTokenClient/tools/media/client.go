package media

import (
	"context"
	"fmt"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/tools/media/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// DouYinToolMediaClient 抖音工具媒体客户端
type DouYinToolMediaClient struct {
	*kernel.BaseClient
}

// NewClient 创建抖音工具媒体客户端实例
func NewClient(c *kernel.BaseClient) *DouYinToolMediaClient {
	return &DouYinToolMediaClient{
		BaseClient: c,
	}
}

// ## Upload 素材上传
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/tools-ability/material-management/upload-material-interface
//
// 参数：
//   ctx  - 请求上下文
//   path - 素材路径 （本地文件路径, 和form二选一）
//   form - 表单数据 （包含素材类型、素材描述等信息 , 和path二选一）
//
// 返回值：
//   *schema.DouYinToolMediaUploadRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
// 			• Media 视频，包含视频ID、URL列表等信息。
//     			• MediaId: 素材ID
//     			• []MediaUrl: 素材URL列表
//   error 调用过程中遇到的错误（如有）
func (c *DouYinToolMediaClient) Upload(ctx context.Context, path string, form *object.HashMap) (*schema.DouYinToolMediaUploadRes, error) {
	result := &schema.DouYinToolMediaUploadRes{}
	// 上传图片
	_, err := c.BaseClient.UploadMedia(ctx, " /enterprise/media/upload/", path, form, nil, result)
	return result, err
}

// ## UploadTempMedia 上传临时素材
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/tools-ability/material-management/upload-temp-material-interface
//
// 参数：
//   ctx  - 请求上下文
//   path - 素材路径 （本地文件路径, 和form二选一）
//   form - 表单数据 （包含素材类型、素材描述等信息 , 和path二选一）
//
// 返回值：
//   *schema.DouYinToolMediaUploadRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
// 			• Media 视频，包含视频ID、URL列表等信息。
//     			• MediaId: 素材ID
//     			• []MediaUrl: 素材URL列表
//   error 调用过程中遇到的错误（如有）
func (c *DouYinToolMediaClient) UploadTempMedia(ctx context.Context, path string, form *object.HashMap) (*schema.DouYinToolMediaUploadRes, error) {
	result := &schema.DouYinToolMediaUploadRes{}
	_, err := c.BaseClient.UploadMedia(ctx, "/enterprise/media/temp/upload/", path, form, nil, result)
	return result, err
}

// ## List 获取素材列表
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/tools-ability/material-management/material-list-interface
//
// 参数：
//   ctx    - 请求上下文
//   cursor - 分页游标，第一次请求传0
//   count  - 每页数量，最大100
//
// 返回值：
//   *schema.DouYinToolMediaGetListRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
// 			• Medias 视频列表，包含视频ID、URL列表等信息。
//     			• MediaId: 素材ID
//     			• []MediaUrl: 素材URL列表
//   error 调用过程中遇到的错误（如有）
func (c *DouYinToolMediaClient) List(ctx context.Context, cursor int64, count int64) (*schema.DouYinToolMediaGetListRes, error) {
	result := &schema.DouYinToolMediaGetListRes{}
	params := &object.StringMap{
		"cursor": fmt.Sprintf("%d", cursor),
		"count":  fmt.Sprintf("%d", count),
	}
	_, err := c.BaseClient.HttpGet(ctx, "/enterprise/media/list/", params, nil, result)
	return result, err
}

// ## Delete 删除素材
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/tools-ability/material-management/delete-material-interface
//
// 参数：
//   ctx     - 请求上下文
//   mediaId - 要删除的素材ID
//
// 返回值：
//   *schema.DouYinToolMediaDeleteRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinToolMediaClient) Delete(ctx context.Context, mediaId string) (*schema.DouYinToolMediaDeleteRes, error) {
	result := &schema.DouYinToolMediaDeleteRes{}
	params := &object.HashMap{
		"media_id": mediaId,
	}
	_, err := c.BaseClient.HttpPost(ctx, "/enterprise/media/delete/", nil, params, nil, result)
	return result, err
}
