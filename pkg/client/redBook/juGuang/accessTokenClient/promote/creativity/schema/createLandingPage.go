package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// H5InfoDto 表示创意图片信息
type H5InfoDto struct {
	PhotoURL       string   `json:"photo_url"`       // 图片链接url
	Content        string   `json:"content"`         // 文案
	ClickURLs      []string `json:"click_urls"`      // 点击链接
	ExpoURLs       []string `json:"expo_urls"`       // 曝光链接
	MonitorCompany string   `json:"monitor_company"` // 监测公司
	MonitorParams  string   `json:"monitor_params"`  // 监测参数配置
}

// PageCreativityInfo 表示前链h5信息
type PageCreativityInfo struct {
	PageID     string      `json:"page_id"`      // 落地页id
	H5InfoDtos []H5InfoDto `json:"h5_info_dtos"` // 创意图片信息
	QualInfo   QualInfo    `json:"qual_info"`    // 资质信息
}

// JuGuangPromoteCreativityCreateLandingPageReq 表示创建落地页创意的请求
type JuGuangPromoteCreativityCreateLandingPageReq struct {
	AdvertiserID        int64                `json:"advertiser_id"`         // 广告主ID
	UnitID              int64                `json:"unit_id"`               // 单元ID
	CreativityName      string               `json:"creativity_name"`       // 创意名称
	PageCreativityInfos []PageCreativityInfo `json:"page_creativity_infos"` // 前链h5信息
}

// JuGuangPromoteCreativityCreateLandingPageRes 表示创建落地页创意的响应
type JuGuangPromoteCreativityCreateLandingPageRes struct {
	response.RedBookAccessTokenRes
	Data struct {
		CreativityID []int64 `json:"creativity_id"` // 创意id集合
	} `json:"data"`
}
