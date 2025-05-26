package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliLiveThirdPartyStartLiveReq 表示 POST /liveopen/fn/live/thirdPartyLive/startLive API 的请求参数
type BiliBiliLiveThirdPartyStartLiveReq struct {
	BizCode   string `json:"biz_code"`   // 必填，从对接运营处获取的业务唯一固定码
	OpenID    string `json:"open_id"`    // 必填，开播授权回调给到的用户唯一ID
	LiveToken string `json:"live_token"` // 必填，开播授权回调给到的第三方开播的一次性token
}

type LiveThirdPartyStartLiveData struct {
	PushURL   string `json:"push_url"`   // 开播推流地址
	LiveKey   string `json:"live_key"`   // 本次开播的唯一ID
	StartTime int64  `json:"start_time"` // 开播时间戳(UTC+8)
}

// BiliBiliLiveThirdPartyStartLiveRes 表示 POST /liveopen/fn/live/thirdPartyLive/startLive API 的响应
type BiliBiliLiveThirdPartyStartLiveRes struct {
	response.BiliBiliRes
	Data LiveThirdPartyStartLiveData `json:"data"`
}
