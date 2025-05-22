package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

type BiliBiliUserGetUserUnionIdRes struct {
	response.BiliBiliRes
	Data *BiliBiliUserUnionIdData `json:"data"`
}

type BiliBiliUserUnionIdData struct {
	UnionId string `json:"union_id"`
}
