package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

type JuGuangToolGetOperationRecordReq struct {
	AdvertiserID int64   `json:"advertiser_id"`       // 品牌方ID
	StartTime    string  `json:"start_time"`          // 开始时间
	EndTime      string  `json:"end_time"`            // 结束时间
	Page         *int    `json:"page,omitempty"`      // 第几页
	PageSize     *int    `json:"page_size,omitempty"` // 页大小
	IDs          []int64 `json:"ids,omitempty"`       // 操作对象ID列表
	Name         *string `json:"name,omitempty"`      // 操作对象名称
	OptType      *int    `json:"opt_type,omitempty"`  // 操作内容
	Level        *int    `json:"level,omitempty"`     // 操作层级
}

type HistoryOperateRecordDTO struct {
	ID             int64  `json:"id"`               // 操作记录ID
	OptTime        string `json:"opt_time"`         // 操作时间
	OptAccountName string `json:"opt_account_name"` // 操作人
	OptIP          string `json:"opt_ip"`           // ip地址
	OptLevelName   string `json:"opt_level_name"`   // 操作层级
	OptObject      string `json:"opt_object"`       // 操作对象
	OptObjectID    int64  `json:"opt_object_id"`    // 操作对象ID
	OptTypeName    string `json:"opt_type_name"`    // 操作内容
	OldValue       string `json:"old_value"`        // 操作前
	NewValue       string `json:"new_value"`        // 操作后
}

type JuGuangToolGetOperationRecordRes struct {
	response.RedBookAccessTokenRes
	// OperationRecordData 操作记录数据
	Data *OperationRecordData `json:"data"` // 返回数据
}

type OperationRecordData struct {
	Total       int64                     `json:"total"`        // 总记录数
	TotalPage   int64                     `json:"total_page"`   // 总页数
	HistoryList []HistoryOperateRecordDTO `json:"history_list"` // 操作记录列表
}
