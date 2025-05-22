package schema

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel/request"
	"github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// BiliBiliVideoCoverUploadReq 视频封面上传请求
type BiliBiliVideoCoverUploadReq struct {
	Files *object.HashMap     `json:"file"` // 图片文件
	Form  *request.UploadForm `json:"form"` // 表单数据
}

// BiliBiliVideoCoverUploadRes 视频封面上传响应
type BiliBiliVideoCoverUploadRes struct {
	response.BiliBiliRes
	Url string `json:"url"` // 封面地址
}
