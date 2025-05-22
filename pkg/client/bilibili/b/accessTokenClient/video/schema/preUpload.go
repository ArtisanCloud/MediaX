package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type BiliBiliVideoPreUploadReq struct {
	Name  string `json:"name"`            // 文件名，需携带正确的扩展名
	Utype string `json:"utype,omitempty"` // 上传类型：0-多分片，1-单个小文件
}

type PreUploadData struct {
	UploadToken string `json:"upload_token"` // 上传凭证
}

type BiliBiliVideoPreUploadRes struct {
	response.BiliBiliRes
	Data PreUploadData `json:"data"`
}
