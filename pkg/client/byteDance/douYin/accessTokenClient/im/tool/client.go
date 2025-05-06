package tool

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/tool/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type DouYinIMToolClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *DouYinIMToolClient {
	return &DouYinIMToolClient{
		BaseClient: c,
	}
}

// ## UploadImage 图片上传
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/business-tool/image-upload
//
// 参数：
//   ctx  - 请求上下文
//   path - 图片路径
//   form - 表单数据
//
// 返回值：
//   *schema.DouYinIMToolUploadImageRes 包含以下字段：
//     • ImageId: 图片ID
//     • Width: 图片宽度
//     • Height: 图片高度
//     • Size: 图片大小
//     • Md5: 图片的MD5值
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMToolClient) UploadImage(ctx context.Context, path string, form *object.HashMap) (*schema.DouYinIMToolUploadImageRes, error) {
	result := &schema.DouYinIMToolUploadImageRes{}

	// 上传图片
	_, err := c.BaseClient.UploadMedia(ctx, "/tool/imagex/client_upload/", path, form, nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ## GetMessageResources 获取消息中的多媒体资源
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/search-management/business-tool/get-message-resources
//
// 参数：
//   ctx  - 请求上下文
//   data - 获取消息中的多媒体资源请求参数
//     • MessageId: 消息ID
//     • ResourceType: 资源类型
//
// 返回值：
//   *schema.DouYinIMToolGetMessageResourcesRes 包含以下字段：
//     • Extra: 通用返回信息（log_id、now、error_code 等）
//     • Data: 业务数据主体，包含具体的业务响应信息
//       • MediaType: 媒体类型
//       • Url: 媒体URL
//   error 调用过程中遇到的错误（如有）
func (c *DouYinIMToolClient) GetMessageResources(ctx context.Context, data *schema.DouYinIMToolGetMessageResourcesReq) (*schema.DouYinIMToolGetMessageResourcesRes, error) {
	result := &schema.DouYinIMToolGetMessageResourcesRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpGet(ctx, "/api/im/message/resources/", params, nil, result)
	return result, err
}
