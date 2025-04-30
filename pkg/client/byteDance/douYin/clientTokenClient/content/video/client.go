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

// ShareResult
//
// ## 获取抖音视频分享结果 ShareResult。
//
// 接口文档参考：
// 请替换为实际的接口文档链接
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • Count: 每页返回数量，必填字段，控制每次查询返回的视频分享结果数量
//	       • Keyword: 查询关键字，例如美食，必填字段，用于筛选视频分享结果
//	       • City: 查询城市，例如上海、北京，可选字段，用于按城市筛选视频分享结果
//	       • Cursor: 分页游标，第一页请求 cursor 是 0，response 中会返回下一页请求用到的 cursor，
//	         同时 response 还会返回 has_more 来表明是否有更多的数据
//
// 返回值：
//
//	*schema.DouYinContentVideoShareResultRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - Pois: 与视频分享结果相关的地点信息列表
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//
//	error 调用过程中遇到的错误（如有）
func (comp *DouYinContentVideoClient) ShareResult(ctx context.Context, data *schema.DouYinContentVideoShareResultReq) (*schema.DouYinContentVideoShareResultRes, error) {
	result := &schema.DouYinContentVideoShareResultRes{}

	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}

	_, err = comp.BaseClient.HttpGet(ctx, "share-id/", params, nil, result)
	return result, err
}

// PoiSearch
//
// ## 查询视频携带的地点（POI）信息 PoiSearch。
//
// 接口文档参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/video-management/douyin/search-video/video-poi
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含：
//	       • Keyword: 搜索关键词
//	       • PageSize: 每页返回数量
//	       • Page: 当前页码
//	       • AccessToken: 授权访问令牌
//
// 返回值：
//
//	*schema.DouYinContentVideoPoiSearchRes 包含以下字段：
//	  • Extra: 通用返回信息（log_id、now、error_code 等）
//	  • Data:
//	      - ShareId: 视频分享 ID
//	      - ErrorCode: 错误码，0 表示成功，其他为失败
//	      - Description: 错误描述或状态说明
//
//	error 调用过程中遇到的错误（如有）的地点（POI）信息。
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
