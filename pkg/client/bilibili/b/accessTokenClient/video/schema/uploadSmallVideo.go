package schema

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel/request"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type BiliBiliVideoUploadSmallVideoReq struct {
	UploadToken string              `json:"upload_token"` // 上传令牌
	Files       *object.HashMap     `json:"files"`        // 文件路径
	Form        *request.UploadForm `json:"form"`         // 表单数据
}

type BiliBiliVideoUploadSmallVideoRes struct {
	Code    int    `json:"code"`    // 返回码
	Message string `json:"message"` // 返回消息
}
