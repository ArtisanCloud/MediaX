package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"

// DouYinIMGroupSetUserEnterAuditReq 设置用户入群审核请求参数
type DouYinIMGroupSetUserEnterAuditReq struct {
	ApplyId string `json:"apply_id"` // 入群申请ID，来源于查询群主所在群的用户入群申请状态
	Status  int64  `json:"status"`   // 目标申请状态，2-通过，3-拒绝
}

// DouYinIMGroupSetUserEnterAuditRes 设置用户入群审核响应参数
type DouYinIMGroupSetUserEnterAuditRes struct {
	// Success 操作是否成功
	Success bool `json:"success"`
	// Extra 包含通用的响应扩展字段，如 log_id、now、error_code 等。
	Extra response.DouYinRes `json:"extra,omitempty"`
	// Data 业务数据主体，包含具体的业务响应信息和视频列表。
	Data struct {
		// Data通用返回信息
		response.DouYinRes
	} `json:"data"`
}
