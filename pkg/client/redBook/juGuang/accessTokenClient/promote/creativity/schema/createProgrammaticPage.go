package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

type H5MaterialInfoPhoto struct {
	PhotoURL string `json:"photo_url"`
}

type H5MaterialInfoTitle struct {
	Title string `json:"title"`
}

type H5MaterialInfo struct {
	Photos         []H5MaterialInfoPhoto `json:"photos"`
	Titles         []H5MaterialInfoTitle `json:"titles"`
	ClickURLs      []string              `json:"click_urls,omitempty"`
	ExpoURLs       []string              `json:"expo_urls,omitempty"`
	MonitorCompany string                `json:"monitor_company,omitempty"`
	MonitorParams  string                `json:"monitor_params,omitempty"`
}

type JuGuangPromoteCreativityCreateProgrammaticPageReq struct {
	AdvertiserID   int64          `json:"advertiser_id"`
	UnitID         int64          `json:"unit_id"`
	CreativityName string         `json:"creativity_name,omitempty"`
	H5MaterialInfo H5MaterialInfo `json:"h5_material_info"`
	QualInfo       QualInfo       `json:"qual_info"`
}

type CreateProgrammaticPageData struct {
	CreativityID int64 `json:"creativity_id"`
}

type JuGuangPromoteCreativityCreateProgrammaticPageRes struct {
	response.RedBookAccessTokenRes
	// CreateProgrammaticPageData 创建程序化落地页创意返回数据
	Data CreateProgrammaticPageData `json:"data"`
}
