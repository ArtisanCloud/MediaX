package schema

import (
	"time"
)

// ## GroupReprot 获取分组报告
//
// 接口文档参考：
// https://ad.xiaohongshu.com/docs-center?bizType=943&articleId=4065
//
// 请求参数
// 字段	类型	是否必填	说明	备注
// advertiser_id	String	是	广告主ID
// startDate	String	是	开始时间 YYYY-MM-DD
// endDate	String	是	结束时间 YYYY-MM-DD
// timeUnit	String	是	时间维度： "DAY"：分天 "SUMMARY"：汇总,默认分天
// columns	List<String>	是	请求列,   必须包含 "groupName", "groupId"
// pageSize	int	否	页面大小	默认20
// pageNum	int	否	页面index	默认1
// filters	List<Struct>	否	过滤选项
// sorts	List<Struct>	否	排序选项
// needTotal	boolean	否	是否需要综合数据
// reportType	String	否	查询维度	默认为 USER_GROUP 人群包维度
// 其他维度：
// CAMPAIGN 计划
// UNIT 单元
// CREATIVITY 创意
// NOTE 笔记
// SPU SPU

// JuGuangDataReportOfflineGroupReprotReq 表示 POST /api/idea/group_report API 的请求参数
type JuGuangDataReportOfflineGroupReprotReq struct {
	AdvertiserID string   `json:"advertiser_id"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	TimeUnit     string   `json:"timeUnit"`
	Columns      []string `json:"columns"`
	PageSize     int      `json:"pageSize,omitempty"`
	PageNum      int      `json:"pageNum,omitempty"`
	Filters      []Filter `json:"filters,omitempty"`
	Sorts        []Sort   `json:"sorts,omitempty"`
	NeedTotal    bool     `json:"needTotal,omitempty"`
	ReportType   string   `json:"reportType,omitempty"`
}

// Filter 表示过滤选项
type Filter struct {
	Column   string   `json:"column"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// Sort 表示排序选项
type Sort struct {
	Column string `json:"column"`
	Sort   string `json:"sort"`
}

// JuGuangDataReportOfflineGroupReprotRes 表示 POST /api/idea/group_report API 的响应
type JuGuangDataReportOfflineGroupReprotRes struct {
	Page            PageInfo        `json:"page"`
	AggregationData AggregationData `json:"aggregationData"`
	DataList        []DataItem      `json:"dataList"`
}

// PageInfo 表示分页信息
type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

// AggregationData 表示综合数据
type AggregationData struct {
	Time           time.Time `json:"time"`
	GroupName      string    `json:"groupName"`
	GroupID        string    `json:"groupId"`
	CampaignName   string    `json:"campaignName"`
	CampaignID     string    `json:"campaignId"`
	UnitName       string    `json:"unitName"`
	UnitID         string    `json:"unitId"`
	CreativityName string    `json:"creativityName"`
	CreativityID   string    `json:"creativityId"`
	NoteID         string    `json:"noteId"`
	SPUName        string    `json:"spuName"`
	SPUID          string    `json:"spuId"`
	Fee            float64   `json:"fee"`
	Impression     int64     `json:"impression"`
	Click          int64     `json:"click"`
	CTR            float64   `json:"ctr"`
	ACP            float64   `json:"acp"`
	CPM            float64   `json:"cpm"`
	Interaction    int64     `json:"interaction"`
}

// DataItem 表示单个数据项
type DataItem struct {
	AggregationData
}
