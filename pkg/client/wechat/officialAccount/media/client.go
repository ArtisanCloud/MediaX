package media

import (
	"context"
	"errors"
	"fmt"
	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/media/response"
	response2 "github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/response"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
	"net/http"
	"os"
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

// 新增临时素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
func (client *Client) UploadImage(ctx context.Context, path string) (*response.UploadMediaRes, error) {
	return client.Upload(ctx, "image", path)
}

func (client *Client) UploadVoice(ctx context.Context, path string) (*response.UploadMediaRes, error) {
	return client.Upload(ctx, "voice", path)
}

func (client *Client) UploadVideo(ctx context.Context, path string) (*response.UploadMediaRes, error) {
	return client.Upload(ctx, "video", path)
}

func (client *Client) UploadThumb(ctx context.Context, path string) (*response.UploadMediaRes, error) {
	return client.Upload(ctx, "thumb", path)
}

// 上传临时素材
// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
func (client *Client) Upload(ctx context.Context, mediaType string, path string) (*response.UploadMediaRes, error) {

	_, err := os.Stat(path)
	if (err != nil && os.IsExist(err)) && (err != nil && os.IsPermission(err)) {
		return nil, errors.New(fmt.Sprintf("File does not exist, or the file is unreadable: \"%s\"", path))
	}

	if !utils.Contains[string](client.AllowTypes, mediaType) {
		return nil, errors.New(fmt.Sprintf("Unsupported media type: '%s'", mediaType))
	}

	outResponse := &response.UploadMediaRes{}
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
	}, nil, outResponse)

	return outResponse, err
}

// 获取临时素材
// https://work.weixin.qq.com/api/doc/90000/90135/90254
func (client *Client) Get(ctx context.Context, mediaID string) (*http.Response, error) {

	header := &response2.HeaderMediaRes{}
	res, err := client.RequestRaw(ctx, "cgi-bin/media/get", http.MethodPost, &object.HashMap{
		"query": &object.StringMap{
			"media_id": mediaID,
		},
	}, header, nil)

	return res, err

}

// 获取高清语音素材
// https://work.weixin.qq.com/api/doc/90000/90135/90255
func (client *Client) GetJSSDK(ctx context.Context, mediaID string) (*http.Response, error) {

	header := &response2.HeaderMediaRes{}
	res, err := client.RequestRaw(ctx, "cgi-bin/media/get/jssdk", http.MethodPost, &object.HashMap{
		"query": &object.StringMap{
			"media_id": mediaID,
		},
	}, header, nil)

	return res, err

}
