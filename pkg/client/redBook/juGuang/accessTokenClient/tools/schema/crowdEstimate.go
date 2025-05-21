package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

type CodeNamePair struct {
	Code     string         `json:"code"`
	Name     string         `json:"name"`
	Children []CodeNamePair `json:"children,omitempty"`
}

type CrowdTarget struct {
	CrowdPkg []struct {
		Value string `json:"value"`
		Name  string `json:"name"`
	} `json:"crowd_pkg"`
}

type IndustryInterestTarget struct {
	ContentInterests  []CodeNamePair `json:"content_interests,omitempty"`
	ShoppingInterests []CodeNamePair `json:"shopping_interests,omitempty"`
}

type JuGuangToolCrowdEstimateReq struct {
	AdvertiserID    int64 `json:"advertiser_id"`
	MarketingTarget int   `json:"marketing_target"`
	Placement       int   `json:"placement"`
	OptimizeTarget  int   `json:"optimize_target"`
	TargetType      int   `json:"target_type"`
	TargetConfig    struct {
		TargetGender           string                  `json:"target_gender"`
		TargetAge              string                  `json:"target_age"`
		TargetCity             string                  `json:"target_city"`
		TargetAreaCode         string                  `json:"target_area_code"`
		TargetDevice           string                  `json:"target_device"`
		IndustryInterestTarget *IndustryInterestTarget `json:"industry_interest_target,omitempty"`
		CrowdTarget            *CrowdTarget            `json:"crowd_target,omitempty"`
		InterestKeywords       []string                `json:"interest_keywords,omitempty"`
		Keywords               []string                `json:"keywords,omitempty"`
		KeywordTargetPeriod    *int                    `json:"keyword_target_period,omitempty"`
		KeywordTargetAction    []int                   `json:"keyword_target_action,omitempty"`
	} `json:"target_config"`
}

type CrowdEstimateData struct {
	CrowdScope  int    `json:"crowd_scope"`
	CrowdNum    string `json:"crowd_num"`
	RawCrowdNum int64  `json:"raw_crowd_num"`
}

type JuGuangToolCrowdEstimateRes struct {
	response.RedBookAccessTokenRes
	Data CrowdEstimateData `json:"data"`
}
