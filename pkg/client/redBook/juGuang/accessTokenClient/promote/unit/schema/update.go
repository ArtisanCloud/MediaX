package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// SpuNoteConfig 商品笔记配置
// SpuId 商品ID
// NoteIds 笔记ID列表
type SpuNoteConfig struct {
	SpuId   string   `json:"spu_id"`
	NoteIds []string `json:"note_ids"`
}

// KeywordWithBid 关键词及出价配置
// Keyword 关键词
// Bid 出价
// KeywordSource 关键词来源
// PhraseMatchType 短语匹配类型
// FeedBid Feed出价
type KeywordWithBid struct {
	Keyword         string `json:"keyword"`
	Bid             int    `json:"bid"`
	KeywordSource   int    `json:"keyword_source"`
	PhraseMatchType int    `json:"phrase_match_type"`
	FeedBid         int    `json:"feed_bid"`
}

type JuGuangPromoteUnitUpdateReq struct {
	AdvertiserId        int64            `json:"advertiser_id"`
	UnitId              int64            `json:"unit_id"`
	UnitName            *string          `json:"unit_name,omitempty"`
	EventBid            *int             `json:"event_bid,omitempty"`
	NoteIds             []string         `json:"note_ids,omitempty"`
	TargetType          *int             `json:"target_type,omitempty"`
	TargetConfig        *TargetConfig    `json:"target_config,omitempty"`
	KeywordTargetPeriod *int             `json:"keyword_target_period,omitempty"`
	KeywordTargetAction []int            `json:"keyword_target_action,omitempty"`
	BusinessTreeName    *string          `json:"business_tree_name,omitempty"`
	SpuNoteInfo         []SpuNoteConfig  `json:"spu_note_info,omitempty"`
	KeywordWithBid      []KeywordWithBid `json:"keyword_with_bid,omitempty"`
	SubstitutedUserId   *string          `json:"substituted_user_id,omitempty"`
	PageId              *string          `json:"page_id,omitempty"`
	LandingPageUrl      *string          `json:"landing_page_url,omitempty"`
	UnitExternalPageUrl *string          `json:"unit_external_page_url,omitempty"`
	UnitLandingPageDesc []string         `json:"unit_landing_page_desc,omitempty"`
	KeywordGenType      *int             `json:"keyword_gen_type,omitempty"`
	TargetTemplateId    *int64           `json:"target_template_id,omitempty"`
}

type UnitUpdateData struct {
	UnitId int64 `json:"unit_id"`
}

type JuGuangPromoteUnitUpdateRes struct {
	response.RedBookAccessTokenRes
	// UnitUpdateData 单元更新数据
	Data UnitUpdateData `json:"data"`
}
