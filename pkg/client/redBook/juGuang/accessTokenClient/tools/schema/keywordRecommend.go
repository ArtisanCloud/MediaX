package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangToolKeywordRecommendReq 表示关键词推荐请求参数
type JuGuangToolKeywordRecommendReq struct {
	AdvertiserID          int64    `json:"advertiser_id"`                     // 广告主ID
	RequestType           string   `json:"request_type"`                      // 推词类型
	PromotionTarget       *int     `json:"promotion_target,omitempty"`        // 推广目标
	RecommendReasonFilter []string `json:"recommend_reason_filter,omitempty"` // 推荐原因过滤
	Keyword               *string  `json:"keyword,omitempty"`                 // 关键词
	ItemIDs               []string `json:"item_ids,omitempty"`                // 笔记ID列表
	TaxonomyID            *string  `json:"taxonomy_id,omitempty"`             // 行业ID
	AttributeList         *string  `json:"attribute_list,omitempty"`          // 行业属性列表
	AttributeNameList     *string  `json:"attribute_name_list,omitempty"`     // 行业属性名称列表
	Rank                  *int     `json:"rank,omitempty"`                    // 排序
}

// WordInfo 表示推荐词信息
type WordInfo struct {
	Keyword          string   `json:"keyword"`           // 词名词
	Source           int      `json:"source"`            // 词来源
	Bid              int      `json:"bid"`               // 市场出价
	CompetitionLevel string   `json:"competition_level"` // 竞争指数
	RecommendReason  []string `json:"recommend_reason"`  // 推荐理由
	Monthpv          int      `json:"monthpv"`           // 月均搜索指数
}

type KeywordRecommendData struct {
	BagMonthPv int        `json:"bag_month_pv"` // 月pv
	WordNum    int        `json:"word_num"`     // 推荐词数量
	WordList   []WordInfo `json:"word_list"`    // 推荐词信息
}

// JuGuangToolKeywordRecommendRes 表示关键词推荐响应参数
type JuGuangToolKeywordRecommendRes struct {
	response.RedBookAccessTokenRes
	Data KeywordRecommendData `json:"data"`
}
