package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMToolUploadImageReq 上传图片请求参数
type DouYinIMToolUploadImageRes struct {
	ImageId string `json:"image_id"` // 图片ID
	Width   int    `json:"width"`    // 图片宽度
	Height  int    `json:"height"`   // 图片高度
	Size    int    `json:"size"`     // 图片大小
	// Md5 图片的MD5值，用于校验图片的完整性。
	// 注意：这个字段可能为空，具体取决于图片是否上传成功。
	// 如果图片上传失败，这个字段可能为空。
	// 如果图片上传成功，这个字段可能包含图片的MD5值。
	// 你可以根据这个字段来判断图片是否上传成功。
	Md5 string `json:"md5"` // 图片的MD5值，用于校验图片的完整性。

	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
