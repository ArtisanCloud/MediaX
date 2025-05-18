package schema

// JuGuangNoteListLandingPageReq 表示获取落地页列表的请求参数
type JuGuangNoteListLandingPageReq struct {
	AdvertiserID *int64  `json:"advertiser_id"`       // 广告主ID
	Page         *int    `json:"page,omitempty"`      // 请求的页码，默认为1
	PageSize     *int    `json:"page_size,omitempty"` // 每页行数，默认为10
	Keyword      *string `json:"keyword,omitempty"`   // 搜索词
	Status       *int    `json:"status"`              // 固定传2(审核通过)
}

// LandingPageVO 表示落地页信息
type LandingPageVO struct {
	AuditComment        *string  `json:"audit_comment"`          // 审核结果备注
	Content             *string  `json:"content"`                // 页面内容json str
	DataUpdated         *bool    `json:"data_updated"`           // 数据是否已更新
	DraftPreviewPath    *string  `json:"draft_preview_path"`     // 草稿页面内容json str
	ID                  *int64   `json:"id"`                     // 页面ID
	LeadsAPI            *string  `json:"leads_api"`              // leads_api
	ObjectID            *string  `json:"object_id"`              // object_id
	OnlineDomain        *string  `json:"online_domain"`          // 在线域名
	PageName            *string  `json:"page_name"`              // 页面名称
	Status              *int     `json:"status"`                 // 落地页状态
	UnitLandingPageDesc []string `json:"unit_landing_page_desc"` // 表单落地页描述
	WebDomain           *string  `json:"web_domain"`             // 域名
	WebPath             *string  `json:"web_path"`               // path
}

// JuGuangNoteListLandingPageRes 表示获取落地页列表的响应
type JuGuangNoteListLandingPageRes struct {
	Code    *int    `json:"code"`    // 返回码
	Msg     *string `json:"msg"`     // 返回信息
	Success *bool   `json:"success"` // 接口是否成功
	Data    struct {
		Total *int            `json:"total"` // 总数
		List  []LandingPageVO `json:"list"`  // 落地页信息
	} `json:"data"`
}
