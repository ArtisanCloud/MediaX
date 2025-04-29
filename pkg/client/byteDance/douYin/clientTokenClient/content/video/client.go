// Package video 提供抖音内容视频相关接口的客户端功能封装。
package video

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/video/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// DouYinContentVideoClient 是抖音内容视频接口的客户端。
type DouYinContentVideoClient struct {
	*kernel.BaseClient
}

// NewClient 初始化并返回一个 DouYinContentVideoClient 实例。
func NewClient(c *kernel.BaseClient) *DouYinContentVideoClient {
	return &DouYinContentVideoClient{
		BaseClient: c,
	}
}

// ShareResult 查询指定分享 ID 的视频发布结果。
//
// 文档参考：https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/douyin/search-video/video-share-result
//
// 参数:
// * ctx: 上下文 Context，用于控制请求的生命周期和传递请求范围的数据
// * data: 请求参数，封装为 DouYinContentVideoShareResultReq 结构体，包含以下字段:
//   - ShareID: 分享 ID，用于标识要查询发布结果的视频分享操作，类型为 string
//   - AccessToken: 访问令牌，用于身份验证，确保请求合法，类型为 string
//
// 返回值：
//   - *schema.DouYinContentVideoShareResultRes：查询结果指针
//   - error：调用过程中遇到的错误（如果有）
func (comp *DouYinContentVideoClient) ShareResult(ctx context.Context, data *schema.DouYinContentVideoShareResultReq) (*schema.DouYinContentVideoShareResultRes, error) {
	result := &schema.DouYinContentVideoShareResultRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = comp.BaseClient.HttpGet(ctx, "share-id/", params, nil, result)
	return result, err
}

// PoiSearch 查询视频携带的地点（POI）信息。
//
// 文档参考：https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/douyin/search-video/video-poi
//
// 参数:
// * ctx: 上下文 Context，用于控制请求的生命周期和传递请求范围的数据
// * data: 请求参数，封装为 DouYinContentVideoPoiSearchReq 结构体，包含以下字段:
//   - Keyword: 用于搜索的关键词，指定要查询的地点信息关键词，类型为 string
//   - PageSize: 每页显示的结果数量，控制每次查询返回的 POI 信息数量，类型为 int
//   - Page: 当前页码，指定要获取的结果页码，类型为 int
//   - AccessToken: 访问令牌，用于身份验证，确保请求合法，类型为 string
//
// 返回值：
// * *schema.DouYinContentVideoPoiSearchRes: 查询结果指针，包含以下字段：
//   - Extra: 包含通用的响应扩展字段，如 log_id、now、error_code 等
//   - Data: 业务数据主体，包含以下字段：
//   - ShareId: 视频分享 ID
//   - ErrorCode: 错误码，0 表示成功，其它表示失败
//   - Description: 错误描述或状态说明
//
// * error: 调用过程中遇到的错误，若请求正常则为 nil
func (comp *DouYinContentVideoClient) PoiSearch(ctx context.Context, data *schema.DouYinContentVideoPoiSearchReq) (*schema.DouYinContentVideoPoiSearchRes, error) {
	result := &schema.DouYinContentVideoPoiSearchRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = comp.BaseClient.HttpGet(ctx, "poi/search/keyword/", params, nil, result)
	return result, err
}
