package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

type JuGuangPromoteCreativityUpdateReq struct {
	AdvertiserId             int64    `json:"advertiser_id"`
	CreativityId             int64    `json:"creativity_id"`
	CreativityName           string   `json:"creativity_name,omitempty"`
	ClickUrls                []string `json:"click_urls,omitempty"`
	ExpoUrls                 []string `json:"expo_urls,omitempty"`
	MaskPerfer               bool     `json:"mask_perfer,omitempty"`
	TitleMaskPerfer          bool     `json:"title_mask_perfer,omitempty"`
	JumpUrl                  string   `json:"jump_url,omitempty"`
	BarContent               string   `json:"bar_content,omitempty"`
	ItemId                   string   `json:"item_id,omitempty"`
	H5Infos                  string   `json:"h5_infos,omitempty"`
	ConversionComponentTypes []int    `json:"conversion_component_types,omitempty"`
	Comment                  string   `json:"comment,omitempty"`
	H5MaterialInfo           string   `json:"h5_material_info,omitempty"`
	PoiId                    string   `json:"poi_id,omitempty"`
	PoiJumpType              string   `json:"poi_jump_type,omitempty"`
	MonitorCompany           string   `json:"monitor_company,omitempty"`
	MonitorParams            string   `json:"monitor_params,omitempty"`
	AdBizItemId              string   `json:"ad_biz_item_id,omitempty"`
	AppCompIcon              string   `json:"app_comp_icon,omitempty"`
	FallBackJumpUrl          string   `json:"fall_back_jump_url,omitempty"`
	PrimaryTitle             string   `json:"primary_title,omitempty"`
	IosDownloadLink          string   `json:"ios_download_link,omitempty"`
	AndroidDownloadLink      string   `json:"android_download_link,omitempty"`
}

type JuGuangPromoteCreativityUpdateRes struct {
	response.RedBookAccessTokenRes
}
