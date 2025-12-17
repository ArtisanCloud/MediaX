package channelBanners

import (
	"context"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v4/channelBanners/schema"
)

// YoutubeChannelBannersClient 是YouTube频道横幅API的客户端
type YoutubeChannelBannersClient struct {
	*kernel.BaseClient
}

// NewClient 创建新的YoutubeChannelBannersClient实例
func NewClient(c *kernel.BaseClient) *YoutubeChannelBannersClient {
	return &YoutubeChannelBannersClient{
		BaseClient: c,
	}
}

// ChannelBanners: insert - 频道横幅
// https://developers.google.cn/youtube/v3/docs/channelBanners/insert?hl=zh-cn
//
// 功能：上传并设置频道的横幅图片
//
// 参数：
//
//	ctx  - 请求上下文
//	data - 请求参数，包含以下字段：
//	  • Files: 要上传的图片文件
//	  • Form: 其他表单参数
//
// 返回值：
//
//	*schema.YoutubeChannelbannersInsertRes 包含以下字段：
//	  • Kind: 资源类型
//	  • ETag: 资源的ETag
//	  • URL: 上传后的图片URL
//	error 调用过程中遇到的错误（如有）
func (c *YoutubeChannelBannersClient) Insert(ctx context.Context, data *schema.YoutubeChannelbannersInsertReq) (*schema.YoutubeChannelbannersInsertRes, error) {
	result := &schema.YoutubeChannelbannersInsertRes{}
	_, err := c.BaseClient.HttpUpload(ctx, "channelBanners", data.Files, data.Form, nil, nil, nil, result)
	return result, err
}
