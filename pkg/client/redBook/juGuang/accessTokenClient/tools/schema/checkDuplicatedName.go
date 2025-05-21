package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"
)

// JuGuangToolCheckDuplicatedNameReq 检查名称重复请求参数
// 字段说明：
//
//	advertiser_id - 广告主ID（必填）
//	type - 查询类型（必填，1：计划，2：单元）
//	name - 名称列表（必填，最多100个）
type JuGuangToolCheckDuplicatedNameReq struct {
	AdvertiserID int64    `json:"advertiser_id"` // 广告主ID
	Type         int      `json:"type"`          // 查询类型
	Name         []string `json:"name"`          // 名称列表
}

// JuGuangToolCheckDuplicatedNameRes 检查名称重复响应结构体
// 字段说明：
//
//	code - 返回码
//	msg - 返回信息
//	success - 接口是否成功
//	data - 业务数据
type JuGuangToolCheckDuplicatedNameRes struct {
	response.RedBookAccessTokenRes
	Data *JuGuangToolCheckDuplicatedNameData `json:"data"` // 业务数据
}

// JuGuangToolCheckDuplicatedNameData 检查名称重复数据结构体
type JuGuangToolCheckDuplicatedNameData struct {
	Type        int             `json:"type"`         // 查询类型
	CheckResult map[string]bool `json:"check_result"` // 检查结果
}
