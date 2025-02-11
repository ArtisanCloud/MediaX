package request

import (
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type MaterialAddMaterialReq struct {
	Type         string          `json:"type"`
	Media        *object.HashMap `json:"media"`
	Title        string          `json:"title"`
	Introduction string          `json:"introduction"`
}

// ---------------------------------------------------

type Article struct {
	Title              string `json:"title"`
	ThumbMediaID       string `json:"thumb_media_id"`
	Author             string `json:"author"`
	Digest             string `json:"digest"`
	ShowCoverPic       string `json:"show_cover_pic"`
	Content            string `json:"content"`
	ContentSourceUrl   string `json:"content_source_url"`
	NeedOpenComment    string `json:"need_open_comment"`
	OnlyFansCanComment string `json:"only_fans_can_comment"`
}

type AddArticlesReq struct {
	Articles []*Article `json:"articles"`
}
