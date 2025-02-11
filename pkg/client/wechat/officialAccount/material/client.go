package material

import (
	"context"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	request2 "github.com/ArtisanCloud/MediaX/internal/kernel/request"
	response2 "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/material/request"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/material/response"
	response3 "github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/response"

	"github.com/ArtisanCloud/MediaXCore/utils/object"
	"net/http"
	"os"
	"path/filepath"
)

type Client struct {
	*kernel.BaseClient

	AllowTypes []string
}

func NewClient(c *kernel.BaseClient) *Client {

	return &Client{
		BaseClient: c,
		AllowTypes: []string{"image", "voice", "video", "thumb", "news_image"},
	}
}

// 上传永久图片素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadImage(ctx context.Context, path string) (*response.MaterialAddMaterialRes, error) {
	result := &response.MaterialAddMaterialRes{}
	_, err := client.Upload(ctx, "image", path, &object.StringMap{}, result)
	return result, err
}

// 上传永久图片素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadImageByData(ctx context.Context, data []byte) (*response.MaterialAddMaterialRes, error) {
	result := &response.MaterialAddMaterialRes{}
	_, err := client.UploadByData(ctx, "image", "image", data, &object.StringMap{}, result)
	return result, err
}

// 上传永久语音素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadVoice(ctx context.Context, path string) (*response.MaterialAddMaterialRes, error) {
	result := &response.MaterialAddMaterialRes{}
	_, err := client.Upload(ctx, "voice", path, &object.StringMap{}, result)
	return result, err
}

// 上传永久语音素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadVoiceByData(ctx context.Context, data []byte) (*response.MaterialAddMaterialRes, error) {
	result := &response.MaterialAddMaterialRes{}
	_, err := client.UploadByData(ctx, "voice", "voice", data, &object.StringMap{}, result)
	return result, err
}

// 上传永久缩略图素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadThumb(ctx context.Context, path string) (*response.MaterialAddMaterialRes, error) {
	result := &response.MaterialAddMaterialRes{}
	_, err := client.Upload(ctx, "thumb", path, &object.StringMap{}, result)
	return result, err
}

// 上传永久缩略图素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadThumbByData(ctx context.Context, data []byte) (*response.MaterialAddMaterialRes, error) {
	result := &response.MaterialAddMaterialRes{}
	_, err := client.UploadByData(ctx, "thumb", "thumb", data, &object.StringMap{}, result)
	return result, err
}

// 上传永久视频素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadVideo(ctx context.Context, path string, title string, description string) (*response.MaterialAddMaterialRes, error) {

	result := &response.MaterialAddMaterialRes{}

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

// 上传永久视频素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadVideoByData(ctx context.Context, data []byte, title string, description string) (*response.MaterialAddMaterialRes, error) {

	result := &response.MaterialAddMaterialRes{}

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

// 新增永久素材
// https://developers.weixin.qq.com/doc/offiaccount/Comments_management/Image_Comments_Management_Interface.html
func (client *Client) UploadArticle(ctx context.Context, articles request.AddArticlesReq) (*response.MaterialAddNewsRes, error) {

	result := &response.MaterialAddNewsRes{}

	//params, err := object.StructToHashMapWithTag(articles, "json")
	params, err := object.StructToHashMap(articles)
	if err != nil {
		return nil, err
	}

	_, err = client.HttpPost(ctx, "cgi-bin/material/add_news", params, nil, result)
	return result, err
}

// 上传永久素材
// https://developers.weixin.qq.com/doc/offiaccount/Comments_management/Image_Comments_Management_Interface.html
func (client *Client) UpdateArticle(ctx context.Context, mediaID string, articles request.AddArticlesReq, index int) (response.MaterialAddNewsRes, error) {
	result := response.MaterialAddNewsRes{}

	params := &object.HashMap{
		"media_id": mediaID,
		"index":    index,
		"articles": articles,
	}

	_, err := client.HttpPost(ctx, "cgi-bin/material/update_news", params, nil, result)
	return result, err
}

// 上传图文消息内的图片获取URL
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (client *Client) UploadArticleImage(ctx context.Context, path string) (*response.MaterialAddMaterialRes, error) {
	result := &response.MaterialAddMaterialRes{}
	_, err := client.Upload(ctx, "news_image", path, &object.StringMap{}, result)
	return result, err
}

// 获取永久素材图片
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Getting_Permanent_Assets.html
func (client *Client) GetMaterial(ctx context.Context, mediaID string) (*http.Response, error) {

	header := &response3.HeaderMediaRes{}
	res, err := client.RequestRaw(ctx, "cgi-bin/material/get_material", http.MethodPost, &object.HashMap{
		"form_params": &object.HashMap{
			"media_id": mediaID,
		},
	}, header, nil)

	return res, err
}

// 获取永久视频消息素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Getting_Permanent_Assets.html
func (client *Client) GetVideo(ctx context.Context, mediaID string) (*response.MaterialGetVideoRes, error) {

	result := &response.MaterialGetVideoRes{}

	options := &object.HashMap{
		"media_id": mediaID,
	}

	_, err := client.HttpPost(ctx, "cgi-bin/material/get_material", options, nil, result)

	return result, err
}

// 获取永久图文素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Getting_Permanent_Assets.html
func (client *Client) GetNews(ctx context.Context, mediaID string) (*response.MaterialGetNewsRes, error) {

	result := &response.MaterialGetNewsRes{}

	options := &object.HashMap{
		"media_id": mediaID,
	}

	_, err := client.HttpPost(ctx, "cgi-bin/material/get_material", options, nil, result)

	return result, err
}

// 删除永久素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Deleting_Permanent_Assets.html
func (client *Client) Delete(ctx context.Context, mediaID string) (*response2.OfficialAccountRes, error) {

	result := &response2.OfficialAccountRes{}

	options := &object.HashMap{
		"media_id": mediaID,
	}

	_, err := client.HttpPost(ctx, "cgi-bin/material/del_material", options, nil, result)

	return result, err
}

// 获取素材列表
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_materials_list.html
func (client *Client) List(ctx context.Context, options *request.MaterialBatchGetMaterialReq) (*response.MaterialBatchGetMaterialRes, error) {

	result := &response.MaterialBatchGetMaterialRes{}

	_, err := client.HttpPost(ctx, "cgi-bin/material/batchget_material", options, nil, result)

	return result, err
}

// 获取素材总数
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_the_total_of_all_materials.html
func (client *Client) Stats(ctx context.Context) (*response.MaterialGetMaterialCountRes, error) {

	result := &response.MaterialGetMaterialCountRes{}

	_, err := client.HttpPost(ctx, "cgi-bin/material/get_materialcount", nil, nil, result)

	return result, err

}

func (client *Client) Upload(ctx context.Context, Type string, path string, query *object.StringMap, result interface{}) (interface{}, error) {

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

func (client *Client) UploadByData(ctx context.Context, Type string, name string, data []byte, query *object.StringMap, result interface{}) (interface{}, error) {

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

func (client *Client) getApiByType(Type string) string {

	switch Type {
	case "news_image":
		return "cgi-bin/media/uploadimg"
	default:
		return "cgi-bin/material/add_material"
	}

}
