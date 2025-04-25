package schema

type YouTubeVideoReportAbuseReq struct {
	VideoId           string `json:"videoId"`
	ReasonId          string `json:"reasonId"`
	SecondaryReasonId string `json:"secondaryReasonId"`
	Comments          string `json:"comments"`
	Language          string `json:"language"`
}
