package schema

import "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/core/response"

// BiliBiliLiveThirdPartyEndLiveReq 表示 POST /liveopen/fn/live/thirdPartyLive/endLive API 的请求参数
type BiliBiliLiveThirdPartyEndLiveReq struct {
	BizCode       string `json:"biz_code"`        // 必填，从对接运营处获取的业务唯一固定码
	OpenID        string `json:"open_id"`         // 必填，用户唯一ID
	LiveKey       string `json:"live_key"`        // 必填，开播后返回的直播live_key
	ThirdLiveUUID string `json:"third_live_uuid"` // 必填，本次开播调用grantUrl时给到的唯一id
}

type LiveThirdPartyEndLiveData struct {
	OpenID  string `json:"open_id"`  // 直播用户唯一ID
	EndTime int64  `json:"end_time"` // 停止直播的时间戳(UTC+8)
}

// BiliBiliLiveThirdPartyEndLiveRes 表示 POST /liveopen/fn/live/thirdPartyLive/endLive API 的响应
type BiliBiliLiveThirdPartyEndLiveRes struct {
	response.BiliBiliRes
	Data LiveThirdPartyEndLiveData `json:"data"`
}
