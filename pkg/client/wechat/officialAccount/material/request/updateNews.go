package request

import (
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type MaterialUpdateNewsReq struct {
	MediaID  int64             `json:"media_id"`
	Index    int64             `json:"index"`
	Articles []*object.HashMap `json:"articles"`
}
