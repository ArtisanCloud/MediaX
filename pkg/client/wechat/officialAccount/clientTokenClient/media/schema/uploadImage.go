package schema

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// HeaderMediaRes 媒体文件头信息响应结构体
type HeaderMediaRes struct {
	response.OfficialAccountRes

	Connection         []string `header:"Connection" json:"Connection"`                   // 连接状态
	ContentType        []string `header:"Content-Type" json:"Content-Type"`               // 内容类型
	ContentDisposition []string `header:"Content-disposition" json:"Content-disposition"` // 内容描述
	Date               []string `header:"Date" json:"Date"`                               // 日期
	CacheControl       []string `header:"Cache-Control" json:"Cache-Control"`             // 缓存控制
	ContentLength      []string `header:"Content-Length" json:"Content-Length"`           // 内容长度
	Content            []string `header:"Content" json:"content"`                         // 内容
}

// UploadImageRes 上传图片响应结构体
type UploadImageRes struct {
	response.OfficialAccountRes

	URL string `json:"url"` // 图片URL
}

// UploadMediaRes 上传媒体文件响应结构体
type UploadMediaRes struct {
	response.OfficialAccountRes

	Item         []*object.HashMap `json:"item"`                     // 媒体文件列表
	Type         string            `json:"type"`                     // 媒体类型
	MediaID      string            `json:"media_id"`                 // 媒体ID
	ThumbMediaID string            `json:"thumb_media_id,omitempty"` // 缩略图媒体ID
	CreatedAt    int               `json:"created_at"`               // 创建时间
}
