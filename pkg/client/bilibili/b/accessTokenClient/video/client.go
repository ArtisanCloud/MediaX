package video

import (
	"context"
	"fmt"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/video/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// BiliBiliVideoClient 是BiliBili视频API的客户端
type BiliBiliVideoClient struct {
	*kernel.BaseClient
}

// NewClient 创建一个新的BiliBiliVideoClient实例
func NewClient(c *kernel.BaseClient) *BiliBiliVideoClient {
	return &BiliBiliVideoClient{
		BaseClient: c,
	}
}

// ## Edit 视频稿件编辑
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/ddfbe54d-b4ab-340e-c698-933b9a3c619c#h1-u89C6u9891u7A3Fu4EF6u7F16u8F91
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • resource_id: 稿件唯一ID（必填）
//	  • title: 稿件标题（必填）
//	  • cover: 封面地址（可选）
//	  • tid: 分区ID（必填）
//	  • no_reprint: 是否禁止转载（可选，默认0）
//	  • desc: 视频描述（可选）
//
// 返回值：
//
//	*schema.BiliBiliVideoEditRes 包含以下字段：
//	  • resource_id: 稿件唯一ID
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 仅允许编辑稿件中的描述信息部分，如title、desc等
// 2. 若需要修改视频本身，请重新投稿
// 3. 稿件编辑后需重新进行审核
func (c *BiliBiliVideoClient) Edit(ctx context.Context, data *schema.BiliBiliVideoEditReq) (*schema.BiliBiliVideoEditRes, error) {
	result := &schema.BiliBiliVideoEditRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/archive/edit", nil, data, nil, result)
	return result, err
}

// ## Delete 视频稿件删除
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/23d78390-4119-1e5f-2bbe-b45bd5cecdb0#h1-u89C6u9891u7A3Fu4EF6u5220u9664
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • resource_id: 稿件唯一ID（必填）
//
// 返回值：
//
//	*schema.BiliBiliVideoDeleteRes 包含以下字段：
//	  • code: 返回码
//	  • message: 返回信息
//	  • request_id: 请求ID
//	  • ttl: 缓存时间
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 该接口需要申请ARC_BASE权限
// 2. 删除操作不可逆，请谨慎使用
// 3. 删除后稿件将无法恢复
func (c *BiliBiliVideoClient) Delete(ctx context.Context, data *schema.BiliBiliVideoDeleteReq) (*schema.BiliBiliVideoDeleteRes, error) {
	result := &schema.BiliBiliVideoDeleteRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/archive/delete", nil, data, nil, result)
	return result, err
}

// ## PreUpload 文件上传预处理
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/0c532c6a-e6fb-0aff-8021-905ae2409095#h2-u6CE8u610Fu4E8Bu9879
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • name: 文件名（必填，需携带正确扩展名，如test.mp4）
//	  • utype: 上传类型（可选，0-多分片，1-单个小文件，默认0）
//
// 返回值：
//
//	*schema.BiliBiliVideoPreUploadRes 包含以下字段：
//	  • upload_token: 上传令牌（用于后续分片上传/合片/投稿接口）
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 文件大小限制：单个文件不超过100M时使用单个小视频上传流程；超过100M时使用分片上传流程
// 2. 最大文件大小：4GB
// 3. 视频时长限制：小于5小时
// 4. 需要ARC_BASE权限
func (c *BiliBiliVideoClient) PreUpload(ctx context.Context, data *schema.BiliBiliVideoPreUploadReq) (*schema.BiliBiliVideoPreUploadRes, error) {
	result := &schema.BiliBiliVideoPreUploadRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/archive/video/init", nil, data, nil, result)
	return result, err
}

// ## UploadPart 文件分片上传
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/733a520a-c50f-7bb4-17cb-35338ba20500#h1-u89C6u9891u7A3Fu4EF6u5220u9664
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • PartNumber: 分片编号（从1开始）
//
// 返回值：
//
//	*schema.BiliBiliVideoUploadPartRes 包含以下字段：
//	  • code: 返回码
//	  • message: 返回信息
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 视频分片大小建议固定为8MB，最后一个分片除外
// 2. 需要ARC_BASE权限
// 3. 分片上传前需先调用PreUpload接口获取upload_token
func (c *BiliBiliVideoClient) UploadPart(ctx context.Context, data *schema.BiliBiliVideoUploadPartReq) (*schema.BiliBiliVideoUploadPartRes, error) {
	result := &schema.BiliBiliVideoUploadPartRes{}
	params := &object.StringMap{
		"part_number": fmt.Sprintf("%d", data.PartNumber),
	}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/archive/video/part", params, data, nil, result)
	return result, err
}

// ## UploadComplete 文件分片合片
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/0828e499-38d8-9e58-2a70-a7eaebf9dd64#h1-u6587u4EF6u5206u7247u5408u7247
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • Files: 上传文件列表（必填）
//	  • Form: 表单数据（可选）
//
// 返回值：
//
//	*schema.BiliBiliVideoUploadCompleteRes 包含以下字段：
//	  • code: 返回码
//	  • message: 返回信息
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 该接口需要申请ARC_BASE权限
// 2. 调用该接口前需先完成所有分片的上传
// 3. 合片成功后即可进行稿件投递
func (c *BiliBiliVideoClient) UploadComplete(ctx context.Context, data *schema.BiliBiliVideoUploadCompleteReq) (*schema.BiliBiliVideoUploadCompleteRes, error) {
	result := &schema.BiliBiliVideoUploadCompleteRes{}
	_, err := c.BaseClient.HttpUpload(ctx, "/arcopen/fn/archive/video/complete", data.Files, data.Form, nil, nil, nil, result)
	return result, err
}

// ## CoverUpload 视频封面上传
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/8243399e-50e3-4058-7f01-1ebe4c632cf8#h1-u5C01u9762u4E0Au4F20
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • Files: 图片文件（必填，格式为jpeg/png，文件大小≤5MB）
//	  • Form: 表单数据（可选）
//
// 返回值：
//
//	*schema.BiliBiliVideoCoverUploadRes 包含以下字段：
//	  • url: 封面地址
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 图片格式支持jpeg/png，文件大小≤5MB
// 2. 建议封面尺寸≥1146x717，最低尺寸≥960x600
// 3. 需要ARC_BASE权限
// 4. 该接口使用multipart/form-data格式上传文件
func (c *BiliBiliVideoClient) CoverUpload(ctx context.Context, data *schema.BiliBiliVideoCoverUploadReq) (*schema.BiliBiliVideoCoverUploadRes, error) {
	result := &schema.BiliBiliVideoCoverUploadRes{}
	_, err := c.BaseClient.HttpUpload(ctx, "/arcopen/fn/archive/cover/upload", data.Files, data.Form, nil, nil, nil, result)
	return result, err
}

// ## SubmitVideo 视频稿件提交
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/f7fc57dd-55a1-5cb1-cba4-61fb2994bf0f#h2-u63A5u53E3u8BF4u660E
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • UploadToken: 上传令牌（必填，通过PreUpload接口获取）
//	  • Title: 稿件标题（必填，长度小于80）
//	  • Cover: 封面地址（可选，必须由上传封面接口得到）
//	  • Tid: 分区ID（必填，由获取分区信息接口得到）
//	  • Tag: 视频标签（必填，多个标签用英文逗号分隔，总长度小于200）
//	  • Desc: 视频描述（可选，长度小于250）
//	  • Copyright: 版权类型（必填，1-原创，2-转载）
//	  • NoReprint: 是否允许转载（可选，0-允许，1-不允许，默认0）
//
// 返回值：
//
//	*schema.BiliBiliVideoSubmitVideoRes 包含以下字段：
//	  • ResourceID: 稿件唯一ID
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 该接口需要申请ARC_BASE权限
// 2. 调用前需确保视频/封面上传已完成
// 3. 非正式会员单日最多投递5个稿件
func (c *BiliBiliVideoClient) SubmitVideo(ctx context.Context, data *schema.BiliBiliVideoSubmitVideoReq) (*schema.BiliBiliVideoSubmitVideoRes, error) {
	result := &schema.BiliBiliVideoSubmitVideoRes{}
	params := &object.StringMap{
		"upload_token": data.UploadToken,
	}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/archive/add-by-utoken", params, data, nil, result)
	return result, err
}

// ## UploadSmallVideo 单个小视频文件上传
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/f22a0eee-c92d-0f1d-f69c-be170cf533c7#h1-u5355u4E2Au5C0Fu89C6u9891u6587u4EF6u4E0Au4F20
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • UploadToken: 上传令牌（必填，通过PreUpload接口获取）
//	  • Files: 文件二进制数据（必填）
//	  • Form: 表单数据（可选）
//
// 返回值：
//
//	*schema.BiliBiliVideoUploadSmallVideoRes 包含以下字段：
//	  • code: 返回码
//	  • message: 返回信息
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 该接口用于上传单个小视频文件，文件大小不超过100MB
// 2. 需要ARC_BASE权限
// 3. 上传成功后无需进行合片操作
func (c *BiliBiliVideoClient) UploadSmallVideo(ctx context.Context, data *schema.BiliBiliVideoUploadSmallVideoReq) (*schema.BiliBiliVideoUploadSmallVideoRes, error) {
	result := &schema.BiliBiliVideoUploadSmallVideoRes{}
	params := &object.StringMap{
		"upload_token": data.UploadToken,
	}
	_, err := c.BaseClient.HttpUpload(ctx, "https://openupos.bilivideo.com/video/v2/part/upload", data.Files, data.Form, params, nil, nil, result)
	return result, err
}

// ## AddShareUrl 获取用于投稿的连接
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/8b825dee-e49e-5843-7e80-f289b878076a#h1-u83B7u53D6u7528u4E8Eu6295u7A3Fu7684u8FDEu63A5
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • SceneCode: 场景码，用于标识唤起的客户端投稿页面
//	  • BizCode: 业务码，需与运营沟通确认
//	  • CommonMsg: 预填内容，格式为JSON字符串
//
// 返回值：
//
//	*schema.BiliBiliVideoAddShareUrlRes 包含以下字段：
//	  • LinkUrl: 唤起客户端投稿页面的URL
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 该接口需要申请ARC_BASE权限
// 2. 获取的LinkUrl可通过移动端系统浏览器打开，已安装客户端时会自动唤起投稿页面
// 3. 未安装客户端时会引导用户进行客户端下载安装
func (c *BiliBiliVideoClient) AddShareUrl(ctx context.Context, data *schema.BiliBiliVideoAddShareUrlReq) (*schema.BiliBiliVideoAddShareUrlRes, error) {
	result := &schema.BiliBiliVideoAddShareUrlRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/arcopen/fn/resource/add_share", nil, data, nil, result)
	return result, err
}

// ## GetVideoInfo 查询单一视频稿件详情
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/d9554788-dcef-f139-6217-b487d41c3826#h1-u67E5u8BE2u5355u4E00u89C6u9891u7A3Fu4EF6u8BE6u60C5
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • resource_id: 稿件唯一ID（必填）
//
// 返回值：
//
//	*schema.BiliBiliVideoGetVideoInfoRes 包含以下字段：
//	  • code: 返回码
//	  • message: 返回信息
//	  • request_id: 请求ID
//	  • ttl: 缓存时间
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 该接口需要申请ARC_BASE权限
// 2. 仅能查询已通过审核的稿件信息
func (c *BiliBiliVideoClient) GetVideoInfo(ctx context.Context, data *schema.BiliBiliVideoGetVideoInfoReq) (*schema.BiliBiliVideoGetVideoInfoRes, error) {
	result := &schema.BiliBiliVideoGetVideoInfoRes{}
	params := &object.StringMap{
		"resource_id": data.ResourceID,
	}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/archive/view", params, nil, nil, result)
	return result, err
}

// ## GetVideoList 查询用户视频稿件列表
//
// 接口文档参考：
// https://open.bilibili.com/doc/4/a24030b7-6b8f-b36c-32d8-a4aae67fcc35#h1-u67E5u8BE2u7528u6237u89C6u9891u7A3Fu4EF6u5217u8868
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • pn: 页码（必填，最小1）
//	  • ps: 每页数量（必填，最小1，最大50）
//
// 返回值：
//
//	*schema.BiliBiliVideoGetVideoListRes 包含以下字段：
//	  • list: 视频稿件列表
//	  • page: 分页信息
//	  • code: 返回码
//	  • message: 返回信息
//	  • request_id: 请求ID
//	  • ttl: 缓存时间
//	error 调用过程中遇到的错误（如有）
//
// 注意：
// 1. 该接口需要申请ARC_BASE权限
// 2. 仅能查询已授权用户的视频稿件列表
// 3. 非正式会员单日最多查询5次
func (c *BiliBiliVideoClient) GetVideoList(ctx context.Context, data *schema.BiliBiliVideoGetVideoListReq) (*schema.BiliBiliVideoGetVideoListRes, error) {
	result := &schema.BiliBiliVideoGetVideoListRes{}
	params := &object.StringMap{
		"pn": fmt.Sprintf("%d", data.Pn),
		"ps": fmt.Sprintf("%d", data.Ps),
	}
	_, err := c.BaseClient.HttpGet(ctx, "/arcopen/fn/archive/viewlist", params, nil, nil, result)
	return result, err
}
