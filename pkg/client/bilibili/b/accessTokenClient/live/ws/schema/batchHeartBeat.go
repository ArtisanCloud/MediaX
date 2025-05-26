package schema

// BiliBiliLiveWSBatchHeartBeatReq 表示批量心跳请求参数
// 接口文档参考：https://member.bilibili.com/arcopen/fn/live/room/ws-batch-heartbeat
// 请求示例：
//
//	{
//	  "conn_ids": ["00000000-0000-0000-0000-000000000000"]
//	}
type BiliBiliLiveWSBatchHeartBeatReq struct {
	ConnIDs []string `json:"conn_ids"` // 获取长链地址时得到的心跳ID列表
}

type LiveWSBatchHeartBeat struct {
	FailedConnIDs []string `json:"failed_conn_ids"` // 心跳失败的心跳ID列表
}

// BiliBiliLiveWSBatchHeartBeatRes 表示批量心跳响应结构
// 响应示例：
//
//	{
//	  "code": 0,
//	  "message": "success",
//	  "data": {
//	    "failed_conn_ids": []
//	  }
//	}
type BiliBiliLiveWSBatchHeartBeatRes struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    LiveWSBatchHeartBeat `json:"data"`
}
