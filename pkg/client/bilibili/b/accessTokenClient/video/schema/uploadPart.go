package schema

type BiliBiliVideoUploadPartReq struct {
	UploadToken string `json:"upload_token"`
	PartNumber  int    `json:"part_number"`
}

type BiliBiliVideoUploadPartRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
