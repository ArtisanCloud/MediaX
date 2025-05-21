package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteGetNegativeWordListReq 表示获取聚光否定词列表的请求参数
// 参考：https://adapi.xiaohongshu.com/api/open/jg/negative/keyword/list
// 字段说明：
//
//	advertiser_id: 广告主ID（必填）
//	unit_id: 单元ID（必填）
//	page_num: 页数，默认1（可选）
//	page_size: 页大小，默认20（可选）
type JuGuangNoteGetNegativeWordListReq struct {
	AdvertiserID int64 `json:"advertiser_id"`       // 广告主ID
	UnitID       int64 `json:"unit_id"`             // 单元ID
	PageNum      int   `json:"page_num,omitempty"`  // 页数，默认1
	PageSize     int   `json:"page_size,omitempty"` // 页大小，默认20
}

type GetNegativeWordListData struct {
	Page             PageRespDTO              `json:"page"`              // 分页信息
	NegativeKeywords []NegativeKeywordItemDTO `json:"negative_keywords"` // 否定词详情列表
}

// JuGuangNoteGetNegativeWordListRes 表示获取聚光否定词列表的响应结构体
type JuGuangNoteGetNegativeWordListRes struct {
	response.RedBookAccessTokenRes
	Data GetNegativeWordListData `json:"data"`
}

// PageRespDTO 分页信息
type PageRespDTO struct {
	PageIndex  int `json:"page_index"`  // 页码
	TotalCount int `json:"total_count"` // 总数量
}

// NegativeKeywordItemDTO 否定词详情
type NegativeKeywordItemDTO struct {
	CampaignID        int64  `json:"campaign_id"`         // 计划id
	UnitID            int64  `json:"unit_id"`             // 单元id
	NegativeKeywordID int64  `json:"negative_keyword_id"` // 否定词id
	Keyword           string `json:"keyword"`             // 否定词
	PhraseMatchType   int    `json:"phrase_match_type"`   // 匹配方式，0-精确匹配，1-短语匹配
}
