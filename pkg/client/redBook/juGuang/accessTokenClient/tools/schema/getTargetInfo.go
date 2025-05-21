package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// ## JuGuangToolGetTargetInfoReq 获取定向信息请求参数
type JuGuangToolGetTargetInfoReq struct {
	AdvertiserID    int64 `json:"advertiser_id"`    // 广告主ID
	MarketingTarget int   `json:"marketing_target"` // 营销目标
}

// ## JuGuangToolGetTargetInfoRes 获取定向信息响应结构体
type JuGuangToolGetTargetInfoRes struct {
	response.RedBookAccessTokenRes
}

// RtbTargetInfo 定向信息结构体
type RtbTargetInfo struct {
	IndustryInterestTarget IndustryInterestTarget `json:"industry_interest_target"` // 行业兴趣
	CrowdTarget            CrowdTarget            `json:"crowd_target"`             // 人群包
	GenderTargets          []CodeNamePair         `json:"gender_targets"`           // 性别
	AgeTargets             []CodeNamePair         `json:"age_targets"`              // 年龄
	AreaTargets            []CodeNamePair         `json:"area_targets"`             // 地域
	DeviceTargets          []CodeNamePair         `json:"device_targets"`           // 设备
}

// CrowdPackageVO 人群包VO结构体
type CrowdPackageVO struct {
	Value      string `json:"value"`       // 人群包ID
	Name       string `json:"name"`        // 人群包名称
	GroupID    string `json:"group_id"`    // 人群包真实ID
	SyncStatus int    `json:"sync_status"` // 同步状态
	Status     int    `json:"status"`      // 删除可用状态
	Type       string `json:"type"`        // 人群包类型
	Tag        string `json:"tag"`         // 行业人群一级类目名称
	Desc       string `json:"desc"`        // 人群包描述
}
