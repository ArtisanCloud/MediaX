package schema

// YoutubeChannelSectionsInsertReq 表示 POST /youtube/v3/channelSections API 的请求参数
type YoutubeChannelSectionsInsertReq struct {
	// Required parameters
	Part string `json:"part" binding:"required"` // 必填，指定API响应包含的channelSection资源属性（如contentDetails,id,snippet）

	// Optional parameters
	OnBehalfOfContentOwner        string `json:"onBehalfOfContentOwner,omitempty"`        // 可选，代表内容所有者执行操作
	OnBehalfOfContentOwnerChannel string `json:"onBehalfOfContentOwnerChannel,omitempty"` // 可选，代表内容所有者频道执行操作

	// Request body
	Snippet ChannelSectionsSnippet `json:"snippet"` // 频道版块的元数据信息

	ContentDetails ContentDetails `json:"contentDetails,omitempty"` // 可选，频道版块的内容详细信息
}

// YoutubeChannelSectionsInsertRes 表示 POST /youtube/v3/channelSections API 的响应
type YoutubeChannelSectionsInsertRes struct {
	// If successful, this method returns a channelSection resource in the response body
	Kind           string                 `json:"kind"`                     // 资源类型，固定为"youtube#channelSection"
	Etag           string                 `json:"etag"`                     // 资源的ETag，用于缓存控制
	Id             string                 `json:"id"`                       // 频道版块的唯一标识符
	Snippet        ChannelSectionsSnippet `json:"snippet"`                  // 频道版块的元数据信息
	ContentDetails ContentDetails         `json:"contentDetails,omitempty"` // 可选，频道版块的内容详细信息
}
