package schema

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel/request"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type BiliBiliVideoUploadCompleteReq struct {
	Files *object.HashMap     // 文件路径
	Form  *request.UploadForm // 表单数据
}

type BiliBiliVideoUploadCompleteRes struct {
	Code    int    `json:"code"`    // 返回码
	Message string `json:"message"` // 返回信息
}
