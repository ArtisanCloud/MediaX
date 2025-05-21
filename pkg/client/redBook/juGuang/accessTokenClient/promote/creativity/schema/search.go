package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

type JuGuangPromoteCreativitySearchReq struct {
	AdvertiserID  int64   `json:"advertiser_id"`
	CampaignID    *int64  `json:"campaign_id,omitempty"`
	UnitID        *int64  `json:"unit_id,omitempty"`
	CreativityIDs []int64 `json:"creativity_ids,omitempty"`
	Status        *int    `json:"status,omitempty"`
	StartTime     *string `json:"start_time,omitempty"`
	EndTime       *string `json:"end_time,omitempty"`
	Page          *Page   `json:"page,omitempty"`
	NoteID        *string `json:"note_id,omitempty"`
}

type Page struct {
	PageIndex *int `json:"page_index,omitempty"`
	PageSize  *int `json:"page_size,omitempty"`
}

type JuGuangPromoteCreativitySearchRes struct {
	response.RedBookAccessTokenRes
}

type JuGuangPromoteCreativitySearchData struct {
	Page           *Page            `json:"page"`
	CreativityDTOs []*CreativityDTO `json:"creativity_dtos"`
}

type CreativityDTO struct {
	AdvertiserID             int64    `json:"advertiser_id"`
	CampaignID               *int64   `json:"campaign_id,omitempty"`
	UnitID                   *int64   `json:"unit_id,omitempty"`
	CreativityID             int64    `json:"creativity_id"`
	CreativityName           string   `json:"creativity_name"`
	CreativityEnable         *int     `json:"creativity_enable,omitempty"`
	CreativityFilterState    *int     `json:"creativity_filter_state,omitempty"`
	CreativityCreateTime     string   `json:"creativity_create_time"`
	MaterialType             *int     `json:"material_type,omitempty"`
	ConversionType           *int     `json:"conversion_type,omitempty"`
	NoteID                   *string  `json:"note_id,omitempty"`
	NoteType                 *int     `json:"note_type,omitempty"`
	CustomMask               *int     `json:"custom_mask,omitempty"`
	CustomTitle              *int     `json:"custom_title,omitempty"`
	TitleFills               []string `json:"title_fills,omitempty"`
	MaskGen                  *int     `json:"mask_gen,omitempty"`
	TitleGen                 *int     `json:"title_gen,omitempty"`
	MaskPrefer               *bool    `json:"mask_prefer,omitempty"`
	TitleMaskPrefer          *bool    `json:"title_mask_prefer,omitempty"`
	AuditStatus              *int     `json:"audit_status,omitempty"`
	AuditComment             *string  `json:"audit_comment,omitempty"`
	PageID                   *string  `json:"page_id,omitempty"`
	ClickURLs                []string `json:"click_urls,omitempty"`
	ExpoURLs                 []string `json:"expo_urls,omitempty"`
	JumpURL                  *string  `json:"jump_url,omitempty"`
	BarContent               *string  `json:"bar_content,omitempty"`
	Image                    *string  `json:"image,omitempty"`
	ItemInvalidReason        *int     `json:"item_invalid_reason,omitempty"`
	ConversionComponentTypes []int64  `json:"conversion_component_types,omitempty"`
	Comment                  *string  `json:"comment,omitempty"`
	Programmatic             *int     `json:"programmatic,omitempty"`
	CreativityExtraInfo      *string  `json:"creativity_extra_info,omitempty"`
	IntoShopParam            *string  `json:"into_shop_param,omitempty"`
	BootScreenInfo           *string  `json:"boot_screen_info,omitempty"`
	PoiID                    *string  `json:"poi_id,omitempty"`
	PoiJumpType              *string  `json:"poi_jump_type,omitempty"`
	MonitorCompany           *string  `json:"monitor_company,omitempty"`
	MonitorParams            *string  `json:"monitor_params,omitempty"`
	ItemID                   *string  `json:"item_id,omitempty"`
	Title                    *string  `json:"title,omitempty"`
	GoodsSellingPoint        *string  `json:"goods_selling_point,omitempty"`
	DataPostURL              *string  `json:"data_post_url,omitempty"`
	KosMsgType               *int     `json:"kos_msg_type,omitempty"`
	QualInfo                 *string  `json:"qual_info,omitempty"`
	MiniProgramPath          *string  `json:"mini_program_path,omitempty"`
	PrimaryTitle             *string  `json:"primary_title,omitempty"`
	ActionButtonContent      *string  `json:"action_button_content,omitempty"`
	HorseRacingResult        *string  `json:"horse_racing_result,omitempty"`
}
