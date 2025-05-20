package schema

// JuGuangNoteBatchAddNegativeWordReq 表示批量添加否定词的请求参数
// 参考：https://adapi.xiaohongshu.com/api/open/jg/negative/keyword/batch/add
// 字段说明：
//
//	advertiser_id: 广告主ID（必填）
//	unit_id: 单元ID（必填）
//	keywords: 否定词列表（必填，元素为 NegativeKeywordAddItemDTO）
type JuGuangNoteBatchAddNegativeWordReq struct {
	AdvertiserID int64                       `json:"advertiser_id"` // 广告主ID
	UnitID       int64                       `json:"unit_id"`       // 单元ID
	Keywords     []NegativeKeywordAddItemDTO `json:"keywords"`      // 否定词列表
}

// NegativeKeywordAddItemDTO 否定词项结构体
//
// keyword: 否定词（必填）
// phrase_match_type: 匹配方式，0-精确匹配，1-短语匹配（必填）
type NegativeKeywordAddItemDTO struct {
	Keyword         string `json:"keyword"`           // 否定词
	PhraseMatchType int    `json:"phrase_match_type"` // 匹配方式
}

// JuGuangNoteBatchAddNegativeWordRes 表示批量添加否定词的响应结构体
// 字段说明：
//
// code: 返回码
// msg: 返回信息
// success: 接口是否成功
// request_id: 请求ID
type JuGuangNoteBatchAddNegativeWordRes struct {
	Code      int    `json:"code"`       // 返回码
	Msg       string `json:"msg"`        // 返回信息
	Success   bool   `json:"success"`    // 接口是否成功
	RequestID string `json:"request_id"` // 请求ID
}
