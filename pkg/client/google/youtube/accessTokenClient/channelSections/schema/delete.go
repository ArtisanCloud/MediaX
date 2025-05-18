package schema

type YoutubeChannelSectionsDeleteReq struct {
	// 必填参数：用于唯一标识要删除的频道分区的ID
	Id string `json:"id"`

	// 可选参数：注意：此参数专门用于YouTube内容合作伙伴。
	// onBehalfOfContentOwner参数表示请求的授权凭据标识了代表参数值中指定的内容所有者行事的YouTube CMS用户
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"`
}
type YoutubeChannelSectionsDeleteRes struct {
	ChannelSection
}
