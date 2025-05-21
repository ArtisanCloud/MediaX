package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/redBook/core/response"

// JuGuangToolGetKeywordMatchReq 关键词匹配请求参数
// 字段说明：
//
//	advertiser_id - 广告主ID（必填）
//	keywords - 关键词列表（必填，最多150个）
type JuGuangToolGetKeywordMatchReq struct {
	AdvertiserID int64    `json:"advertiser_id"` // 广告主ID
	Keywords     []string `json:"keywords"`      // 关键词列表
}

// JuGuangToolGetKeywordMatchRes 关键词匹配响应结构体
// 字段说明：
//
//	code - 返回码
//	msg - 返回信息
//	success - 接口是否成功
//	data - 业务数据
type JuGuangToolGetKeywordMatchRes struct {
	response.RedBookAccessTokenRes
	Data *JuGuangToolGetKeywordMatchData `json:"data"` // 业务数据
}

// JuGuangToolGetKeywordMatchData 关键词匹配数据结构体
type JuGuangToolGetKeywordMatchData struct {
	MatchDistinctCount int64                      `json:"match_distinct_count"` // 匹配个数
	MatchInfos         []JuGuangToolMatchInfoItem `json:"match_infos"`          // 匹配信息列表
}

// JuGuangToolMatchInfoItem 匹配信息项结构体
type JuGuangToolMatchInfoItem struct {
	Keyword     string `json:"keyword"`      // 关键词
	InThesaurus bool   `json:"in_thesaurus"` // 是否匹配词库
}
