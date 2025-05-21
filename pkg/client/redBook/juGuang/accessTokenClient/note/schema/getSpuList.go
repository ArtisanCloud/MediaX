package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteGetSpuListReq 表示获取SPU列表的请求参数
type JuGuangNoteGetSpuListReq struct {
	AdvertiserID int64   `json:"advertiser_id"`       // 广告主ID
	CanBind      *bool   `json:"can_bind,omitempty"`  // 是否可以绑定
	Keyword      *string `json:"keyword,omitempty"`   // 搜索关键词
	Page         *int    `json:"page,omitempty"`      // 页码
	PageSize     *int    `json:"page_size,omitempty"` // 页大小
}

type GetSpuListData struct {
	Total int         `json:"total"` // 总数
	Spu   []SpuDetail `json:"spu"`   // SPU列表
}

// JuGuangNoteGetSpuListRes 表示获取SPU列表的响应
type JuGuangNoteGetSpuListRes struct {
	response.RedBookAccessTokenRes
	Data GetSpuListData `json:"data"`
}

// SpuDetail 表示SPU详细信息
type SpuDetail struct {
	MainSpuID      int64    `json:"main_spu_id"`      // 主SPU ID
	SpuID          int64    `json:"spu_id"`           // SPU ID
	SpuName        string   `json:"spu_name"`         // SPU名称
	SpuStatus      int      `json:"spu_status"`       // SPU状态
	TaxonomyCode   string   `json:"taxonomy_code"`    // SPU类目code
	BrandID        string   `json:"brand_id"`         // 品牌ID
	NickNameList   []string `json:"nick_name_list"`   // SPU昵称列表
	SeriesList     []string `json:"series_list"`      // SPU系列名称列表
	PicUrlList     []string `json:"pic_url_list"`     // SPU图片链接列表
	SpuAuditReason string   `json:"spu_audit_reason"` // 审核拒绝原因
}
