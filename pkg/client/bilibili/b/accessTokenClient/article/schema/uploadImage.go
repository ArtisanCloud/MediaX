package schema

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel/request"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// BiliBiliArticleUploadImageReq 表示上传图片的请求参数
type BiliBiliArticleUploadImageReq struct {
	Files     *object.HashMap     `json:"file"`      // 图片文件
	Form      *request.UploadForm `json:"form"`      // 表单数据
	Watermark bool                `json:"watermark"` // 是否带水印
}

// BiliBiliArticleUploadImageRes 表示上传图片的响应参数
type BiliBiliArticleUploadImageRes struct {
	response.BiliBiliRes
	Data *UploadImageData `json:"data"` // 图片上传结果数据
}

// UploadImageData 表示图片上传结果的具体数据
type UploadImageData struct {
	URL  string `json:"url"`  // 图片URL
	Size int    `json:"size"` // 图片大小
}
