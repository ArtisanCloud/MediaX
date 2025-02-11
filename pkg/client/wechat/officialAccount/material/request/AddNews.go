package request

import (
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type MaterialAddNewsReq struct {
	Articles []*object.HashMap `json:"articles"`
}
