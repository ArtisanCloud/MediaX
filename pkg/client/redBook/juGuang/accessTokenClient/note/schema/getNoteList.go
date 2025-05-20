package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteGetNoteListReq 表示获取笔记列表的请求参数
type JuGuangNoteGetNoteListReq struct {
	AdvertiserID    int64   `json:"advertiser_id"`               // 广告主ID
	NoteType        int     `json:"note_type"`                   // 笔记类型
	Keyword         *string `json:"keyword,omitempty"`           // 搜索关键词
	OrderField      *string `json:"order_field,omitempty"`       // 排序字段
	OrderType       *string `json:"order_type,omitempty"`        // 排序类型
	NoteContentType *int    `json:"note_content_type,omitempty"` // 笔记内容类型
	PlacementType   *int    `json:"placement_type,omitempty"`    // 推广场景
	SpuID           *string `json:"spu_id,omitempty"`            // spu_id
	FilterTaobao    *int    `json:"filter_taobao,omitempty"`     // 是否只展示小红星笔记
	MarketTarget    *int    `json:"market_target,omitempty"`     // 营销诉求
	SpuType         *int    `json:"spu_type,omitempty"`          // spu类型
	Page            *int    `json:"page,omitempty"`              // 页码
	PageSize        *int    `json:"page_size,omitempty"`         // 每页行数
	BaseOnly        bool    `json:"base_only"`                   // 是否只拉取笔记基本信息
}

// JuGuangNoteGetNoteListRes 表示获取笔记列表的响应
type JuGuangNoteGetNoteListRes struct {
	response.RedBookAccessTokenRes
	Data JuGuangNoteGetNoteListResData `json:"data"` // 返回数据
}

// JuGuangNoteGetNoteListResData 表示获取笔记列表的响应数据
type JuGuangNoteGetNoteListResData struct {
	Total int               `json:"total"` // 总数
	Notes []BaseNoteItemDTO `json:"notes"` // 笔记信息
}

// BaseNoteItemDTO 表示笔记基本信息
type BaseNoteItemDTO struct {
	NoteID                 string                `json:"note_id"`                  // 笔记ID
	Image                  string                `json:"image"`                    // 图片
	Desc                   string                `json:"desc"`                     // 笔记内容
	CreateTime             int64                 `json:"create_time"`              // 创建时间
	Author                 string                `json:"author"`                   // 笔记作者
	AuthorImage            string                `json:"author_image"`             // 作者头像
	Status                 int                   `json:"status"`                   // 笔记状态
	NoteContentType        int                   `json:"note_content_type"`        // 笔记类型
	CooperateState         bool                  `json:"cooperate_state"`          // 是否合作笔记
	Title                  string                `json:"title"`                    // 标题
	ImageList              []string              `json:"image_list"`               // 图片列表
	CooperateComponentType int                   `json:"cooperate_component_type"` // 合作组件类型
	CrowdCreationNote      bool                  `json:"crowd_creation_note"`      // 是否为共创笔记
	ItemID                 string                `json:"item_id"`                  // 商品ID
	ReadCount              int                   `json:"read_count"`               // 阅读数
	ReadRate               string                `json:"read_rate"`                // 阅读率
	InteractCount          int                   `json:"interact_count"`           // 互动数
	InteractRate           string                `json:"interact_rate"`            // 互动率
	HighQuality            int                   `json:"high_quality"`             // 优质笔记
	HighPotential          int                   `json:"high_potential"`           // 高潜笔记
	OutsideShopVisit       int                   `json:"outside_shop_visit"`       // 站外进店量
	OutsideShopVisitRate   string                `json:"outside_shop_visitRate"`   // 站外进店率
	ItemIDs                []string              `json:"item_ids"`                 // 笔记挂接的商品
	NoteMultiSpuInfo       []NoteMultiSpuInfoDTO `json:"note_multi_spu_info"`      // 笔记绑定spu信息
	Taxonomy1              string                `json:"taxonomy1"`                // 一级类目信息
	Taxonomy2              string                `json:"taxonomy2"`                // 二级类目信息
	Taxonomy3              string                `json:"taxonomy3"`                // 三级类目信息
	IsHitStrategy          int                   `json:"is_hit_strategy"`          // 是否命中策略
	HitStrategyContent     string                `json:"hit_strategy_content"`     // 命中策略内容
	WinHorseNote           bool                  `json:"win_horse_note"`           // 是否为优胜笔记
	StaffTag               string                `json:"staff_tag"`                // 员工标签
	StaffArea              string                `json:"staff_area"`               // 地域
	NoteURL                string                `json:"note_url"`                 // 笔记链接
	HasShopCard            bool                  `json:"has_shop_card"`            // 是否为商品笔记
}

// NoteMultiSpuInfoDTO 表示笔记绑定spu信息
type NoteMultiSpuInfoDTO struct {
	BindID                   int64  `json:"bind_id"`                       // 绑定ID
	SpuID                    string `json:"spu_id"`                        // spu_id
	SpuName                  string `json:"spu_name"`                      // spu名称
	ExceedModifyLimitIn30Day bool   `json:"exceed_modify_limit_in_30_day"` // 是否超出30内修改绑定三次限制
	ExceedModifyLimitToday   bool   `json:"exceed_modify_limit_today"`     // 是否超出一天内修改绑定一次限制
	BindByCurAccount         bool   `json:"bind_by_cur_account"`           // 是否为当前账户绑定
	BindAuditStatus          int    `json:"bind_audit_status"`             // 绑定关系审核状态
	BindAuditReason          string `json:"bind_audit_reason"`             // 绑定关系审核拒绝原因
	SpuType                  int    `json:"spu_type"`                      // spu类型
	SpuSubName               string `json:"spu_sub_name"`                  // 非标、品牌名称
	SeriesID                 string `json:"series_id"`                     // 系列 ID
}
