package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangNoteGetDirectLinkListReq 表示 GET /api/open/jg/direct-link/list 的请求参数
type JuGuangNoteGetDirectLinkListReq struct {
	AdvertiserID *int64  `json:"advertiser_id,omitempty"` // 广告主id
	ID           *int64  `json:"id,omitempty"`            // 直达链接id
	PageNum      *int    `json:"page_num,omitempty"`      // 页码，从1开始
	PageSize     *int    `json:"page_size,omitempty"`     // 页大小，最大100
	RemarkName   *string `json:"remark_name,omitempty"`   // 备注名称，模糊匹配
	TypeStr      *string `json:"type_str,omitempty"`      // 类型，1-deeplink，2-ulk
	StatusStr    *string `json:"status_str,omitempty"`    // 状态，1-审核中，2-审核通过，3-审核拒绝
}

// JuGuangNoteGetDirectLinkListRes 表示 GET /api/open/jg/direct-link/list 的响应
type JuGuangNoteGetDirectLinkListRes struct {
	response.JuGuangRes
	Data DirectLinkListData `json:"data"` // 数据
}

// DirectLinkListData 表示直达链接列表数据
type DirectLinkListData struct {
	Total          int64        `json:"total"`            // 直达链接总数
	DirectLinkList []DirectLink `json:"direct_link_list"` // 直达链接列表
}

// DirectLink 表示单个直达链接
type DirectLink struct {
	ID         int64  `json:"id"`          // 直达链接id
	URL        string `json:"url"`         // url内容
	Type       int    `json:"type"`        // 类型，1-deeplink，2-ulk
	Status     int    `json:"status"`      // 状态，1-审核中，2-审核通过，3-审核拒绝
	RemarkName string `json:"remark_name"` // 备注名称
}
