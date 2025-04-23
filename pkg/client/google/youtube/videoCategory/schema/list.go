package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/video/schema"

type Resource struct {
	Kind    string         `json:"kind"`
	Etag    interface{}    `json:"etag"`
	Id      string         `json:"id"`
	Snippet schema.Snippet `json:"snippet"`
}
type YouTubeVideoCategoriesReq struct {
	Part       string `json:"part"`
	Id         string `json:"id,omitempty"`
	RegionCode string `json:"regionCode,omitempty"`
	HL         string `json:"hl,omitempty"`
}

type YouTubeVideoCategoriesRes struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	PrevPageToken string `json:"prevPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []interface{} `json:"items"`
}
