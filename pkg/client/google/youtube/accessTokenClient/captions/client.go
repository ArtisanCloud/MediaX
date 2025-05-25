package captions

import (
	"context"
	"io"
	"net/http"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/captions/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type YoutubeCaptionsClient struct {
	*kernel.BaseClient
}

func NewClient(c *kernel.BaseClient) *YoutubeCaptionsClient {
	return &YoutubeCaptionsClient{
		BaseClient: c,
	}
}

// ## Captions: list - 字幕列表
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/captions/list?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet,id 等）
//	  • videoId: 视频ID（必填）
//	  • id: 字幕ID列表（可选）
//	  • onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//
// 返回值：
//
//	*schema.YoutubeCaptionsListRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Items: 字幕列表
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCaptionsClient) List(ctx context.Context, data *schema.YoutubeCaptionsListReq) (*schema.YoutubeCaptionsListRes, error) {
	result := &schema.YoutubeCaptionsListRes{}
	params, err := object.StructToStringMap(data)
	_, err = c.BaseClient.HttpGet(ctx, "/youtube/v3/captions", params, nil, nil, result)
	return result, err
}

// ## Captions: insert - 上传字幕
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/captions/insert?hl=zh-cn
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定返回的资源部分（必填，如 snippet）
//	  • sync: 是否同步上传（可选）
//	  • onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//	  • file: 字幕文件内容（必填）
//
// 返回值：
//
//	*schema.YoutubeCaptionsInsertRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 字幕ID
//	  • Snippet: 字幕信息
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeCaptionsClient) Insert(ctx context.Context, data *schema.YoutubeCaptionsInsertReq) (*schema.YoutubeCaptionsInsertRes, error) {
	result := &schema.YoutubeCaptionsInsertRes{}
	_, err := c.BaseClient.HttpPost(ctx, "/youtube/v3/captions", &object.StringMap{
		"part":                   data.Part,
		"onBehalfOfContentOwner": data.OnBehalfOfContentOwner,
	}, data, nil, result)
	return result, err
}

// ## Captions: update - 更新字幕
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/captions/update?hl=zh-cn
//
// 功能：
// 更新字幕轨道，可更改轨道的草稿状态或上传新的字幕文件。
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • part: 指定写入操作的属性（必填，如snippet或id）
//	  • onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//	  • file: 更新后的字幕文件内容（可选）
//
// 返回值：
//
//	*schema.YoutubeCaptionsUpdateRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • Id: 字幕ID
//	  • Snippet: 字幕信息
//	error 调用过程中遇到的错误（如有）
//
// 注意事项：
// 1. 上传的文件大小上限为100MB
// 2. 调用此方法会产生450个单元的配额费用
func (c *YoutubeCaptionsClient) Update(ctx context.Context, data *schema.YoutubeCaptionsUpdateReq) (*schema.YoutubeCaptionsUpdateRes, error) {
	result := &schema.YoutubeCaptionsUpdateRes{}
	_, err := c.BaseClient.HttpPut(ctx, "/youtube/v3/captions", &object.StringMap{
		"part":                   data.Part,
		"onBehalfOfContentOwner": data.OnBehalfOfContentOwner,
	}, data, nil, result)
	return result, err
}

// ## Captions: download - 下载字幕
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/captions/download?hl=zh-cn
//
// 功能：
// 下载指定字幕轨道的内容。默认以原始格式和语言返回字幕文件，
// 可通过参数指定返回格式和翻译语言。
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 字幕轨道ID（必填）
//	  • tfmt: 字幕格式（可选，如 sbv, srt, vtt 等）
//	  • tlang: 翻译语言代码（可选，ISO 639-1格式）
//	  • onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//
// 返回值：
//
//	*schema.YoutubeCaptionsDownloadRes 包含以下字段：
//	  • ContentType: 返回内容的MIME类型
//	  • Data: 字幕文件二进制内容
//	error 调用过程中遇到的错误（如有）
//
// 注意事项：
// 1. 调用此方法会产生200个单位的配额费用
// 2. 文件下载大小限制为100MB
func (c *YoutubeCaptionsClient) Download(ctx context.Context, data *schema.YoutubeCaptionsDownloadReq) (*schema.YoutubeCaptionsDownloadRes, error) {
	result := &schema.YoutubeCaptionsDownloadRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	resp, err := c.BaseClient.HttpGet(ctx, "/youtube/v3/captions/id", params, nil, nil, result)
	if err != nil {
		return nil, err
	}
	httpResp := resp.(*http.Response)
	result.ContentType = httpResp.Header.Get("Content-Type")
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	result.Data = body
	return result, nil
}

// ## Captions: delete - 删除字幕
//
// 接口文档参考：
// https://developers.google.cn/youtube/v3/docs/captions/delete?hl=zh-cn
//
// 功能：
// 删除指定的字幕轨道。
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • id: 字幕轨道ID（必填）
//	  • onBehalfOfContentOwner: 代表内容所有者执行操作（可选）
//
// 返回值：
//
//	*schema.YoutubeCaptionsDeleteRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	error 调用过程中遇到的错误（如有）
//
// 注意事项：
// 1. 调用此方法会产生450个单位的配额费用
// 2. 删除操作不可逆，请谨慎操作
func (c *YoutubeCaptionsClient) Delete(ctx context.Context, data *schema.YoutubeCaptionsDeleteReq) (*schema.YoutubeCaptionsDeleteRes, error) {
	result := &schema.YoutubeCaptionsDeleteRes{}
	params, err := object.StructToStringMap(data)
	if err != nil {
		return nil, err
	}
	_, err = c.BaseClient.HttpDelete(ctx, "/youtube/v3/captions", params, nil, nil, result)
	return result, err
}
