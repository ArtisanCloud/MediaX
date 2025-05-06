package material

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	request2 "github.com/ArtisanCloud/MediaX/internal/kernel/request"
	response2 "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/material/schema"
	schema2 "github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/media/schema"

	"github.com/ArtisanCloud/MediaXCore/utils/object"
	"net/http"
	"os"
	"path/filepath"
)

// OfficialAccountMaterialClient 微信公众号素材管理客户端
type OfficialAccountMaterialClient struct {
	*kernel.BaseClient

	AllowTypes []string
}

// NewClient 创建新的微信公众号素材管理客户端实例
func NewClient(c *kernel.BaseClient) *OfficialAccountMaterialClient {

	return &OfficialAccountMaterialClient{
		BaseClient: c,
		AllowTypes: []string{"image", "voice", "video", "thumb", "news_image"},
	}
}

// ## UploadImage 上传永久图片素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
//
// 参数：
//
//	ctx  - 请求上下文
//	path - 图片文件路径
//
// 返回值：
//
//	*schema.MaterialAddMaterialRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • URL: 素材URL
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadImage(ctx context.Context, path string) (*schema.MaterialAddMaterialRes, error) {
	result := &schema.MaterialAddMaterialRes{}
	_, err := client.Upload(ctx, "image", path, &object.StringMap{}, result)
	return result, err
}

// ## UploadImageByData 通过二进制数据上传永久图片素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 图片二进制数据
//
// 返回值：
//
//	*schema.MaterialAddMaterialRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • URL: 素材URL
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadImageByData(ctx context.Context, data []byte) (*schema.MaterialAddMaterialRes, error) {
	result := &schema.MaterialAddMaterialRes{}
	_, err := client.UploadByData(ctx, "image", "image", data, &object.StringMap{}, result)
	return result, err
}

// ## UploadVoice 上传永久语音素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
//
// 参数：
//
//	ctx  - 请求上下文
//	path - 语音文件路径
//
// 返回值：
//
//	*schema.MaterialAddMaterialRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • URL: 素材URL
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadVoice(ctx context.Context, path string) (*schema.MaterialAddMaterialRes, error) {
	result := &schema.MaterialAddMaterialRes{}
	_, err := client.Upload(ctx, "voice", path, &object.StringMap{}, result)
	return result, err
}

// ## UploadVoiceByData 通过二进制数据上传永久语音素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 语音二进制数据
//
// 返回值：
//
//	*schema.MaterialAddMaterialRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • URL: 素材URL
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadVoiceByData(ctx context.Context, data []byte) (*schema.MaterialAddMaterialRes, error) {
	result := &schema.MaterialAddMaterialRes{}
	_, err := client.UploadByData(ctx, "voice", "voice", data, &object.StringMap{}, result)
	return result, err
}

// ## UploadThumb 上传永久缩略图素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
//
// 参数：
//
//	ctx  - 请求上下文
//	path - 缩略图文件路径
//
// 返回值：
//
//	*schema.MaterialAddMaterialRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • URL: 素材URL
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadThumb(ctx context.Context, path string) (*schema.MaterialAddMaterialRes, error) {
	result := &schema.MaterialAddMaterialRes{}
	_, err := client.Upload(ctx, "thumb", path, &object.StringMap{}, result)
	return result, err
}

// 上传永久缩略图素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *OfficialAccountMaterialClient) UploadThumbByData(ctx context.Context, data []byte) (*schema.MaterialAddMaterialRes, error) {
	result := &schema.MaterialAddMaterialRes{}
	_, err := client.UploadByData(ctx, "thumb", "thumb", data, &object.StringMap{}, result)
	return result, err
}

// ## UploadThumbByData 通过二进制数据上传永久缩略图素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 缩略图二进制数据
//
// 返回值：
//
//	*schema.MaterialAddMaterialRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • URL: 素材URL
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadVideo(ctx context.Context, path string, title string, description string) (*schema.MaterialAddMaterialRes, error) {

	result := &schema.MaterialAddMaterialRes{}

	jsonDescription, err := object.JsonEncode(&object.HashMap{
		"title":        title,
		"introduction": description,
	})
	if err != nil {
		return nil, err
	}

	params := &object.StringMap{
		"Description": jsonDescription,
	}

	_, err = client.Upload(ctx, "video", path, params, result)
	return result, err
}

// ## UploadVideoByData 通过二进制数据上传永久视频素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
//
// 参数：
//
//	ctx         - 请求上下文
//	data        - 视频二进制数据
//	title       - 视频标题
//	description - 视频描述
//
// 返回值：
//
//	*schema.MaterialAddMaterialRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • URL: 素材URL
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadVideoByData(ctx context.Context, data []byte, title string, description string) (*schema.MaterialAddMaterialRes, error) {

	result := &schema.MaterialAddMaterialRes{}

	jsonDescription, err := object.JsonEncode(&object.HashMap{
		"title":        title,
		"introduction": description,
	})
	if err != nil {
		return nil, err
	}

	params := &object.StringMap{
		"Description": jsonDescription,
	}

	_, err = client.UploadByData(ctx, "video", "video", data, params, result)
	return result, err
}

// ## UploadArticle 新增永久图文素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Comments_management/Image_Comments_Management_Interface.html
//
// 参数：
//
//	ctx      - 请求上下文
//	articles - 图文素材内容，包含以下字段：
//	  • Title: 标题
//	  • ThumbMediaID: 封面图片媒体ID
//	  • Author: 作者
//	  • Digest: 摘要
//	  • ShowCoverPic: 是否显示封面图片
//	  • Content: 正文内容
//	  • ContentSourceURL: 原文链接
//	  • NeedOpenComment: 是否打开评论
//	  • OnlyFansCanComment: 是否仅粉丝可评论
//
// 返回值：
//
//	*schema.MaterialAddNewsRes 包含以下字段：
//	  • MediaID: 媒体ID
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadArticle(ctx context.Context, articles schema.AddArticlesReq) (*schema.MaterialAddNewsRes, error) {

	result := &schema.MaterialAddNewsRes{}

	//params, err := object.StructToHashMapWithTag(articles, "json")
	params, err := object.StructToHashMap(articles)
	if err != nil {
		return nil, err
	}

	_, err = client.HttpPost(ctx, "cgi-bin/material/add_news", nil, params, nil, result)
	return result, err
}

// ## UpdateArticle 更新永久图文素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Comments_management/Image_Comments_Management_Interface.html
//
// 参数：
//
//	ctx      - 请求上下文
//	mediaID  - 媒体ID
//	articles - 图文素材内容，包含以下字段：
//	  • Title: 标题
//	  • ThumbMediaID: 封面图片媒体ID
//	  • Author: 作者
//	  • Digest: 摘要
//	  • ShowCoverPic: 是否显示封面图片
//	  • Content: 正文内容
//	  • ContentSourceURL: 原文链接
//	  • NeedOpenComment: 是否打开评论
//	  • OnlyFansCanComment: 是否仅粉丝可评论
//	index    - 要更新的文章在图文素材中的索引
//
// 返回值：
//
//	schema.MaterialAddNewsRes 包含以下字段：
//	  • MediaID: 媒体ID
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UpdateArticle(ctx context.Context, mediaID string, articles schema.AddArticlesReq, index int) (schema.MaterialAddNewsRes, error) {
	result := schema.MaterialAddNewsRes{}

	params := &object.HashMap{
		"media_id": mediaID,
		"index":    index,
		"articles": articles,
	}

	_, err := client.HttpPost(ctx, "cgi-bin/material/update_news", nil, params, nil, result)
	return result, err
}

// ## UploadArticleImage 上传图文消息内的图片获取URL
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
//
// 参数：
//
//	ctx  - 请求上下文
//	path - 图片文件路径
//
// 返回值：
//
//	*schema.MaterialAddMaterialRes 包含以下字段：
//	  • MediaID: 媒体ID
//	  • URL: 素材URL
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) UploadArticleImage(ctx context.Context, path string) (*schema.MaterialAddMaterialRes, error) {
	result := &schema.MaterialAddMaterialRes{}
	_, err := client.Upload(ctx, "news_image", path, &object.StringMap{}, result)
	return result, err
}

// ## GetMaterial 获取永久素材图片
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Getting_Permanent_Assets.html
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
func (client *OfficialAccountMaterialClient) GetMaterial(ctx context.Context, mediaID string) (*http.Response, error) {

	header := &schema2.HeaderMediaRes{}
	res, err := client.RequestRaw(ctx, "cgi-bin/material/get_material", http.MethodPost, nil, &object.HashMap{
		"form_params": &object.HashMap{
			"media_id": mediaID,
		},
	}, header, nil)

	return res, err
}

// ## GetVideo 获取永久视频消息素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Getting_Permanent_Assets.html
//
// 参数：
//
//	ctx    - 请求上下文
//	mediaID - 媒体ID
//
// 返回值：
//
//	*schema.MaterialGetVideoRes 包含以下字段：
//	  • Title: 视频标题
//	  • Description: 视频描述
//	  • DownURL: 视频下载地址
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) GetVideo(ctx context.Context, mediaID string) (*schema.MaterialGetVideoRes, error) {

	result := &schema.MaterialGetVideoRes{}

	options := &object.HashMap{
		"media_id": mediaID,
	}

	_, err := client.HttpPost(ctx, "cgi-bin/material/get_material", nil, options, nil, result)

	return result, err
}

// ## GetNews 获取永久图文素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Getting_Permanent_Assets.html
//
// 参数：
//
//	ctx    - 请求上下文
//	mediaID - 媒体ID
//
// 返回值：
//
//	*schema.MaterialGetNewsRes 包含以下字段：
//	  • NewsItem: 图文素材列表，每个元素包含以下字段：
//	    • Title: 标题
//	    • ThumbMediaID: 封面图片媒体ID
//	    • ShowCoverPic: 是否显示封面图片
//	    • Author: 作者
//	    • Digest: 摘要
//	    • Content: 正文内容
//	    • URL: 原文链接
//	    • ContentSourceURL: 原文链接
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) GetNews(ctx context.Context, mediaID string) (*schema.MaterialGetNewsRes, error) {

	result := &schema.MaterialGetNewsRes{}

	options := &object.HashMap{
		"media_id": mediaID,
	}

	_, err := client.HttpPost(ctx, "cgi-bin/material/get_material", nil, options, nil, result)

	return result, err
}

// ## Delete 删除永久素材
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Deleting_Permanent_Assets.html
//
// 参数：
//
//	ctx    - 请求上下文
//	mediaID - 媒体ID
//
// 返回值：
//
//	*response2.OfficialAccountRes 包含以下字段：
//	  • ErrCode: 错误码
//	  • ErrMsg: 错误信息
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) Delete(ctx context.Context, mediaID string) (*response2.OfficialAccountRes, error) {

	result := &response2.OfficialAccountRes{}

	options := &object.HashMap{
		"media_id": mediaID,
	}

	_, err := client.HttpPost(ctx, "cgi-bin/material/del_material", nil, options, nil, result)

	return result, err
}

// ## List 获取素材列表
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_materials_list.html
//
// 参数：
//
//	ctx     - 请求上下文
//	options - 查询条件，包含以下字段：
//	  • Type: 素材类型（image/voice/video/news）
//	  • Offset: 从全部素材的该偏移位置开始返回
//	  • Count: 返回素材的数量
//
// 返回值：
//
//	*schema.MaterialBatchGetMaterialRes 包含以下字段：
//	  • TotalCount: 该类型的素材总数
//	  • ItemCount: 本次调用获取的素材数量
//	  • Item: 素材列表，每个元素包含以下字段：
//	    • MediaID: 媒体ID
//	    • Name: 文件名
//	    • UpdateTime: 更新时间
//	    • URL: 素材URL（仅图片、语音、视频素材有）
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) List(ctx context.Context, options *schema.MaterialBatchGetMaterialReq) (*schema.MaterialBatchGetMaterialRes, error) {

	result := &schema.MaterialBatchGetMaterialRes{}

	_, err := client.HttpPost(ctx, "cgi-bin/material/batchget_material", nil, options, nil, result)

	return result, err
}

// ## Stats 获取素材总数
//
// 接口文档参考：
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_the_total_of_all_materials.html
//
// 参数：
//
//	ctx - 请求上下文
//
// 返回值：
//
//	*schema.MaterialGetMaterialCountRes 包含以下字段：
//	  • VoiceCount: 语音素材总数
//	  • VideoCount: 视频素材总数
//	  • ImageCount: 图片素材总数
//	  • NewsCount: 图文素材总数
//
//	error 调用过程中遇到的错误（如有）
func (client *OfficialAccountMaterialClient) Stats(ctx context.Context) (*schema.MaterialGetMaterialCountRes, error) {

	result := &schema.MaterialGetMaterialCountRes{}

	_, err := client.HttpPost(ctx, "cgi-bin/material/get_materialcount", nil, nil, nil, result)

	return result, err

}

func (client *OfficialAccountMaterialClient) Upload(ctx context.Context, Type string, path string, query *object.StringMap, result interface{}) (interface{}, error) {

	_, err := os.Stat(path)
	if (err != nil && os.IsExist(err)) && (err != nil && os.IsPermission(err)) {
		return "", err
	}

	var files *object.HashMap
	if path != "" {
		files = &object.HashMap{
			"media": path,
		}
	}

	(*query)["type"] = Type

	form := &request2.UploadForm{
		FileName: filepath.Base(path),
	}

	return client.HttpUpload(ctx, client.getApiByType(Type), files, form, query, nil, result)
}

func (client *OfficialAccountMaterialClient) UploadByData(ctx context.Context, Type string, name string, data []byte, query *object.StringMap, result interface{}) (interface{}, error) {

	formData := &request2.UploadForm{
		Contents: []*request2.UploadContent{
			&request2.UploadContent{
				Name:  name,
				Value: data,
			},
		},
	}

	return client.HttpUpload(ctx, client.getApiByType(Type), nil, formData, query, nil, result)
}

func (client *OfficialAccountMaterialClient) getApiByType(Type string) string {

	switch Type {
	case "news_image":
		return "cgi-bin/media/uploadimg"
	default:
		return "cgi-bin/material/add_material"
	}

}
