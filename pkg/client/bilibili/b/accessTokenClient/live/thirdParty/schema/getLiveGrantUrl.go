package schema

// BiliBiliLiveThirdPartyGetLiveGrantUrlReq 表示获取第三方直播授权URL请求参数
// 接口文档参考：https://member.bilibili.com/liveopen/fn/live/thirdPartyLive/grantUrl
// 请求示例：
// GET https://member.bilibili.com/liveopen/fn/live/thirdPartyLive/grantUrl?biz_code=xxx&open_id=xxx&live_area_id=1&third_live_uuid=xxx
type BiliBiliLiveThirdPartyGetLiveGrantUrlReq struct {
	BizCode       string `json:"biz_code"`        // 申请接入时获得的biz_code
	OpenID        string `json:"open_id"`         // 开播用户的open_id
	LiveAreaID    int    `json:"live_area_id"`    // 目标开播分区
	ThirdLiveUUID string `json:"third_live_uuid"` // 第三方提供本次开播的唯一ID
}

type LiveThirdPartyGetLiveGrantUrlData struct {
	GrantURL string `json:"grant_url"` // 授权链接地址
}

// BiliBiliLiveThirdPartyGetLiveGrantUrlRes 表示获取第三方直播授权URL响应结构
// 响应示例：
//
//	{
//	  "code": 0,
//	  "message": "0",
//	  "data": {
//	    "grant_url": ""
//	  }
//	}
type BiliBiliLiveThirdPartyGetLiveGrantUrlRes struct {
	Code    int                               `json:"code"`
	Message string                            `json:"message"`
	Data    LiveThirdPartyGetLiveGrantUrlData `json:"data"`
}
