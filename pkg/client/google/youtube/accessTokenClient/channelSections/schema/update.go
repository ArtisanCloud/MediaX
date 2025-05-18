package schema

// YoutubeChannelSectionsUpdateReq 定义更新YouTube频道版块的请求参数
// 官方文档参考：https://developers.google.com/youtube/v3/docs/channelSections/update
type YoutubeChannelSectionsUpdateReq struct {
	// 必需参数
	Part string `json:"part" binding:"required"` // 以逗号分隔的频道分区资源属性列表(contentDetails,id,snippet)

	// 可选参数
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"` // 代表内容所有者执行操作的CMS用户

	// 请求主体
	Snippet ChannelSectionsSnippet `json:"snippet"` // 频道版块的元数据信息

	ContentDetails ContentDetails `json:"contentDetails,omitempty"` // 频道版块的内容详细信息
}

// YoutubeChannelSectionsUpdateRes 定义更新YouTube频道版块的响应结果
type YoutubeChannelSectionsUpdateRes struct {
	// 如果成功，此方法将在响应正文中返回一个channelSection资源
	Kind           string                 `json:"kind"`                     // 资源类型，固定为 "youtube#channelSection"
	Etag           string                 `json:"etag"`                     // 资源的ETag，用于检查资源是否已更改
	Id             string                 `json:"id"`                       // 频道版块的唯一标识符
	Snippet        ChannelSectionsSnippet `json:"snippet"`                  // 频道版块的元数据信息
	ContentDetails ContentDetails         `json:"contentDetails,omitempty"` // 频道版块的内容详细信息
}
