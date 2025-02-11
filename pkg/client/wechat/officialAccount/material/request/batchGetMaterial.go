package request

type MaterialBatchGetMaterialReq struct {
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
	Count  int64  `json:"count"`
}
