package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteBatchDeleteNegativeWordReq 表示批量删除否定词的请求参数
// 参考：https://adapi.xiaohongshu.com/api/open/jg/negative/keyword/batch/delete
// 字段说明：
//
//	advertiser_id: 广告主ID（必填）
//	unit_id: 单元ID（必填）
//	negative_keyword_ids: 否定词id列表（必填）
type JuGuangNoteBatchDeleteNegativeWordReq struct {
	AdvertiserID       int64   `json:"advertiser_id"`        // 广告主ID
	UnitID             int64   `json:"unit_id"`              // 单元ID
	NegativeKeywordIDs []int64 `json:"negative_keyword_ids"` // 否定词id列表
}

// JuGuangNoteBatchDeleteNegativeWordRes 表示批量删除否定词的响应结构体
// 字段说明：
//
//	code: 返回码
//	msg: 返回信息
//	success: 接口是否成功
//	request_id: 请求ID
type JuGuangNoteBatchDeleteNegativeWordRes struct {
	response.RedBookAccessTokenRes
}
