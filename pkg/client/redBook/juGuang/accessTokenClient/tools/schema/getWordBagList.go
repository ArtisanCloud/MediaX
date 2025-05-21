package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangToolGetWordBagListReq 获取词包列表请求参数
// 字段说明：
//
//	advertiser_id - 广告主ID（必填）
//	name - 搜索名称（可选）
//	category - 分类（可选）
//	page_num - 页码（可选，默认1）
//	page_size - 页大小（可选，默认10）
//	start_time - 开始时间（可选）
//	end_time - 结束时间（可选）
type JuGuangToolGetWordBagListReq struct {
	AdvertiserID int64   `json:"advertiser_id"`        // 广告主ID
	Name         *string `json:"name,omitempty"`       // 搜索名称
	Category     *string `json:"category,omitempty"`   // 分类
	PageNum      *int    `json:"page_num,omitempty"`   // 页码
	PageSize     *int    `json:"page_size,omitempty"`  // 页大小
	StartTime    *string `json:"start_time,omitempty"` // 开始时间
	EndTime      *string `json:"end_time,omitempty"`   // 结束时间
}

// JuGuangToolGetWordBagListRes 获取词包列表响应结构体
// 字段说明：
//
//	code - 返回码
//	msg - 返回信息
//	success - 接口是否成功
//	request_id - 请求ID
//	data - 业务数据
type JuGuangToolGetWordBagListRes struct {
	response.RedBookAccessTokenRes
	Data *JuGuangToolGetWordBagListData `json:"data"` // 业务数据
}

// JuGuangToolGetWordBagListData 词包列表数据结构体
type JuGuangToolGetWordBagListData struct {
	Page           *JuGuangToolGetWordBagListPage `json:"page"`              // 分页信息
	WordTagDtoList []JuGuangToolWordTagDto        `json:"word_tag_dto_list"` // 词包列表
}

// JuGuangToolGetWordBagListPage 分页信息结构体
type JuGuangToolGetWordBagListPage struct {
	PageNum    int `json:"page_num"`    // 当前页码
	TotalCount int `json:"total_count"` // 总词包数
}

// JuGuangToolWordTagDto 词包信息结构体
type JuGuangToolWordTagDto struct {
	Name          string                `json:"name"`           // 词包名称
	Source        int                   `json:"source"`         // 词包来源
	WordList      []JuGuangToolWordItem `json:"word_list"`      // 词列表
	CreateAudit   string                `json:"create_audit"`   // 创建账号
	CreateTime    string                `json:"create_time"`    // 创建时间
	KeywordSource int                   `json:"keyword_source"` // 词来源
}

// JuGuangToolWordItem 词项结构体
type JuGuangToolWordItem struct {
	Keyword          string   `json:"keyword"`           // 关键词
	Source           int      `json:"source"`            // 词来源
	Bid              int      `json:"bid"`               // 市场出价
	CompetitionLevel string   `json:"competition_level"` // 竞争指数
	RecommendReason  []string `json:"recommend_reason"`  // 推荐理由
	MonthPv          int      `json:"monthpv"`           // 月均搜索指数
}
