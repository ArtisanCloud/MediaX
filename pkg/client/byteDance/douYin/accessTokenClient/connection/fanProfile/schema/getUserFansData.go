package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// ActiveDaysDistribution 活跃天数分布
type ActiveDaysDistribution struct {
	Item  string `json:"item"`  // 活跃天数区间，如 "0-3天"
	Value int    `json:"value"` // 该区间的粉丝数量
}

// AgeDistribution 年龄分布
type AgeDistribution struct {
	Item  string `json:"item"`  // 年龄区间，如 "18-23岁"
	Value int    `json:"value"` // 该区间的粉丝数量
}

// DeviceDistribution 设备分布
type DeviceDistribution struct {
	Item  string `json:"item"`  // 设备类型，如 "iPhone"
	Value int    `json:"value"` // 使用该设备的粉丝数量
}

// FlowContribution 流量贡献
type FlowContribution struct {
	AllSum  int    `json:"all_sum"`  // 总流量
	FansSum int    `json:"fans_sum"` // 粉丝贡献的流量
	Flow    string `json:"flow"`     // 流量类型，如 "播放量"
}

// GenderDistribution 性别分布
type GenderDistribution struct {
	Item  string `json:"item"`  // 性别，如 "男"
	Value int    `json:"value"` // 该性别的粉丝数量
}

// GeographicalDistribution 地域分布
type GeographicalDistribution struct {
	Item  string `json:"item"`  // 地域名称，如 "北京"
	Value int    `json:"value"` // 该地域的粉丝数量
}

// InterestDistribution 兴趣分布
type InterestDistribution struct {
	Item  string `json:"item"`  // 兴趣标签，如 "美食"
	Value int    `json:"value"` // 对该兴趣感兴趣的粉丝数量
}

// FansData 粉丝画像数据
type FansData struct {
	ActiveDaysDistributions   []ActiveDaysDistribution   `json:"active_days_distributions"`  // 活跃天数分布
	AgeDistributions          []AgeDistribution          `json:"age_distributions"`          // 年龄分布
	AllFansNum                int                        `json:"all_fans_num"`               // 总粉丝数
	DeviceDistributions       []DeviceDistribution       `json:"device_distributions"`       // 设备分布
	FlowContributions         []FlowContribution         `json:"flow_contributions"`         // 流量贡献
	GenderDistributions       []GenderDistribution       `json:"gender_distributions"`       // 性别分布
	GeographicalDistributions []GeographicalDistribution `json:"geographical_distributions"` // 地域分布
	InterestDistributions     []InterestDistribution     `json:"interest_distributions"`     // 兴趣分布
}

// DouYinConnectionFanProfileGetUserFansDataRes 获取用户粉丝数据接口响应结构
type DouYinConnectionFanProfileGetUserFansDataRes struct {
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
		FansData FansData `json:"fans_data"` // 粉丝画像数据
	} `json:"data"`
}
