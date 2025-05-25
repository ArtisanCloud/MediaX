package media

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/media/schema"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// OfficialAccountMediaClient 是一个用于操作微信公众号素材的客户端。
type OfficialAccountMediaClient struct {
	*kernel.BaseClient

	AllowTypes []string
}

// NewClient 创建一个新的 OfficialAccountMediaClient 实例。
func NewClient(c *kernel.BaseClient) *OfficialAccountMediaClient {
	return &OfficialAccountMediaClient{
		BaseClient: c,
		AllowTypes: []string{"image", "voice", "video", "thumb", "news_image"},
	}
}

// ## UploadImage 上传临时图片素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
//
// 参数：
//
//	ctx  - 请求上下文
//	path - 图片文件路径
//
// 返回值：
//
//	*schema.UploadMediaRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • CreatedAt: 素材上传时间戳
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMediaClient) UploadImage(ctx context.Context, path string) (*schema.UploadMediaRes, error) {
	return client.Upload(ctx, "image", path)
}

// ## UploadVoice 上传临时语音素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
//
// 参数：
//
//	ctx  - 请求上下文
//	path - 语音文件路径
//
// 返回值：
//
//	*schema.UploadMediaRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • CreatedAt: 素材上传时间戳
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMediaClient) UploadVoice(ctx context.Context, path string) (*schema.UploadMediaRes, error) {
	return client.Upload(ctx, "voice", path)
}

// ## UploadVideo 上传临时视频素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
//
// 参数：
//
//	ctx  - 请求上下文
//	path - 视频文件路径
//
// 返回值：
//
//	*schema.UploadMediaRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • CreatedAt: 素材上传时间戳
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMediaClient) UploadVideo(ctx context.Context, path string) (*schema.UploadMediaRes, error) {
	return client.Upload(ctx, "video", path)
}

// ## UploadThumb 上传临时缩略图素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
//
// 参数：
//
//	ctx  - 请求上下文
//	path - 缩略图文件路径
//
// 返回值：
//
//	*schema.UploadMediaRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • CreatedAt: 素材上传时间戳
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMediaClient) UploadThumb(ctx context.Context, path string) (*schema.UploadMediaRes, error) {
	return client.Upload(ctx, "thumb", path)
}

// ## Upload 上传临时素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
//
// 参数：
//
//	ctx      - 请求上下文
//	mediaType - 素材类型（image/voice/video/thumb）
//	path     - 文件路径
//
// 返回值：
//
//	*schema.UploadMediaRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • CreatedAt: 素材上传时间戳
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMediaClient) Upload(ctx context.Context, mediaType string, path string) (*schema.UploadMediaRes, error) {
	_, err := os.Stat(path)
	if (err != nil && os.IsExist(err)) && (err != nil && os.IsPermission(err)) {
		return nil, errors.New(fmt.Sprintf("File does not exist, or the file is unreadable: \"%s\"", path))
	}

	if !utils.Contains[string](client.AllowTypes, mediaType) {
		return nil, errors.New(fmt.Sprintf("Unsupported media type: '%s'", mediaType))
	}

	outResponse := &schema.UploadMediaRes{}
	var files *object.HashMap
	if path != "" {
		files = &object.HashMap{
			"media": path,
		}
	} else {
		return nil, errors.New("path is empty")
	}

	_, err = client.HttpUpload(ctx, "cgi-bin/media/upload", files, nil, &object.StringMap{
		"type": mediaType,
	}, nil, nil, outResponse)

	return outResponse, err
}

// ## Get 获取临时素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_temporary_materials.html
//
// 参数：
//
//	ctx    - 请求上下文
//	mediaID - 媒体ID
//
// 返回值：
//
//	*http.Response 包含素材的HTTP响应
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMediaClient) Get(ctx context.Context, mediaID string) (*http.Response, error) {
	header := &schema.HeaderMediaRes{}
	res, err := client.RequestRaw(ctx, "cgi-bin/media/get", http.MethodPost, nil, &object.HashMap{
		"query": &object.StringMap{
			"media_id": mediaID,
		},
	}, header, nil)

	return res, err
}

// ## Get 获取临时素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_temporary_materials.html
//
// 参数：
//
//	ctx    - 请求上下文
//	mediaID - 媒体ID
//
// 返回值：
//
//	*http.Response 包含素材的HTTP响应
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMediaClient) GetJSSDK(ctx context.Context, mediaID string) (*http.Response, error) {
	header := &schema.HeaderMediaRes{}
	res, err := client.RequestRaw(ctx, "cgi-bin/media/get/jssdk", http.MethodPost, nil, &object.HashMap{
		"query": &object.StringMap{
			"media_id": mediaID,
		},
	}, header, nil)

	return res, err
}
