package schema

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel/request"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type YouTubeThumbnailsSetReq struct {
	VideoId                string              `json:"videoId"`                          // 视频ID（必填）
	OnBehalfOfContentOwner string              `json:"onBehalfOfContentOwner,omitempty"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	Files                  *object.HashMap     // 文件路径
	Form                   *request.UploadForm // 表单数据
}

type Default struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Medium struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type High struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Standard struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type MaxRes struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type YouTubeThumbnailsSetRes struct {
	Default  Default  `json:"default"`
	Medium   Medium   `json:"medium"`
	High     High     `json:"high"`
	Standard Standard `json:"standard"`
	MaxRes   MaxRes   `json:"maxres"`
}
