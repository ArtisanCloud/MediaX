package schema

type YouTubeVideoRateReq struct {
	ID     string `json:"id"`
	Rating string `json:"rating"`
}
type YouTubeVideoRateRes struct {
}

type YouTubeVideoGetRatingReq struct {
	ID                     string `json:"id"`
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"`
}

type Item struct {
	VideoId string `json:"videoId"`
	Rating  string `json:"rating"`
}

type YouTubeVideoGetRatingRes struct {
	Kind  string `json:"kind"`
	Etag  string `json:"etag"`
	Items []Item `json:"items"`
}
