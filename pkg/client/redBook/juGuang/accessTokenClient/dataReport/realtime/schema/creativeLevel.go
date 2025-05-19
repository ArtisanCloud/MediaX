package schema

// JuGuangDataReportRealtimeCreativeLevelReq 表示获取创意层级实时报表数据的请求结构体
type JuGuangDataReportRealtimeCreativeLevelReq struct {
	AdvertiserID              int64  `json:"advertiser_id"`
	StartDate                 string `json:"start_date"`
	EndDate                   string `json:"end_date"`
	PageNum                   int    `json:"page_num,omitempty"`
	PageSize                  int    `json:"page_size,omitempty"`
	SortColumn                string `json:"sort_column,omitempty"`
	Sort                      string `json:"sort,omitempty"`
	PlacementList             []int  `json:"placement_list,omitempty"`
	CreativityFilterState     int    `json:"creativity_filter_state,omitempty"`
	CreativityCreateBeginTime string `json:"creativity_create_begin_time,omitempty"`
	CreativityCreateEndTime   string `json:"creativity_create_end_time,omitempty"`
	ConversionType            int    `json:"conversion_type,omitempty"`
	ProgrammaticList          []int  `json:"programmatic_list,omitempty"`
	CreativityAuditState      int    `json:"creativity_audit_state,omitempty"`
	Name                      string `json:"name,omitempty"`
	ID                        int    `json:"id,omitempty"`
	DataCaliber               int    `json:"data_caliber,omitempty"`
	NeedHourlyData            bool   `json:"need_hourly_data,omitempty"`
}

// JuGuangDataReportRealtimeCreativeLevelRes 表示获取创意层级实时报表数据的响应结构体
type BaseCreativityDTO struct {
	CreativityID          int64  `json:"creativity_id"`
	CreativityName        string `json:"creativity_name"`
	CreativityFilterState int    `json:"creativity_filter_state"`
	CreativityCreateTime  string `json:"creativity_create_time"`
	CreativityEnable      int    `json:"creativity_enable"`
	AuditStatus           int    `json:"audit_status"`
	UnitID                int64  `json:"unit_id"`
	Programmatic          int    `json:"programmatic"`
	NoteID                string `json:"note_id"`
	CreativityType        int    `json:"creativity_type"`
}

type CreativityDTO struct {
	Data              DataReportDTO     `json:"data"`
	BaseCampaignDTO   BaseCampaignDTO   `json:"base_campaign_dto"`
	BaseUnitDTO       BaseUnitDTO       `json:"base_unit_dto"`
	BaseCreativityDTO BaseCreativityDTO `json:"base_creativity_dto"`
	HourlyData        []DataReportDTO   `json:"hourly_data,omitempty"`
}

type JuGuangDataReportRealtimeCreativeLevelRes struct {
	Code           int             `json:"code"`
	Msg            string          `json:"msg"`
	Success        bool            `json:"success"`
	Page           PageRespDTO     `json:"page"`
	CreativityDTOs []CreativityDTO `json:"creativity_dtos"`
	TotalData      DataReportDTO   `json:"total_data"`
}
