package schema

import (
	"github.com/ArtisanCloud/MediaX/internal/kernel/request"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type YoutubeChannelbannersInsertReq struct {
	VideoId                string              `json:"videoId"`                          // 视频ID（必填）
	OnBehalfOfContentOwner string              `json:"onBehalfOfContentOwner,omitempty"` // 内容所有者（可选，仅供 YouTube 内容合作伙伴使用）
	Files                  *object.HashMap     // 文件路径
	Form                   *request.UploadForm // 表单数据
}

type ChannelBanners struct {
	Kind string `json:"kind"` // 资源类型
	Etag string `json:"etag"` // 资源的ETag
	Url  string `json:"url"`  // 资源的URL
}

// 插入频道横幅
type YoutubeChannelbannersInsertRes struct {
	ChannelBanners
}
