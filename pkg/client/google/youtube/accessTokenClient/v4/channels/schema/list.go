package schema

// YoutubeChannelsListReq 表示获取YouTube频道列表的请求参数
type YoutubeChannelsListReq struct {
	//
	Part string `json:"part"`

	CategoryId  string `json:"categoryId,omitempty"`
	ForHandle   string `json:"forHandle,omitempty"`
	ForUsername string `json:"forUsername,omitempty"`
	Id          string `json:"id,omitempty"`
	ManagedByMe bool   `json:"managedByMe,omitempty"`
	Mine        bool   `json:"mine,omitempty"`

	// Optional parameters
	Hl                     string `json:"hl,omitempty"`                     // Language for localized resource metadata
	MaxResults             uint   `json:"maxResults,omitempty"`             // Max items to return (0-50, default 5)
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner,omitempty"` // Content owner making request
	PageToken              string `json:"pageToken,omitempty"`              // Token for specific page of results
}

// PageInfo 包含分页信息
// TotalResults 表示总结果数
// ResultsPerPage 表示每页结果数
type PageInfo struct {
	TotalResults   uint `json:"totalResults"`
	ResultsPerPage uint `json:"resultsPerPage"`
}

// YoutubeChannelsListRes 表示获取YouTube频道列表的响应结果
// Kind 表示资源类型
// Etag 表示资源的ETag
// NextPageToken 表示下一页的分页令牌
// PrevPageToken 表示上一页的分页令牌
// PageInfo 包含分页信息
// Items 包含频道列表数据
type YoutubeChannelsListRes struct {
	Kind          string        `json:"kind"`
	Etag          string        `json:"etag"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
	PrevPageToken string        `json:"prevPageToken,omitempty"`
	PageInfo      PageInfo      `json:"pageInfo"`
	Items         []ChannelItem `json:"items"`
}

// ChannelItem 表示单个YouTube频道的详细信息
// Kind 表示资源类型
// Etag 表示资源的ETag
// Id 表示频道ID
// Snippet 包含频道的基本信息
// ContentDetails 包含频道的内容详情
// Statistics 包含频道的统计信息
// TopicDetails 包含频道的主题信息
// Status 包含频道的状态信息
// BrandingSettings 包含频道的品牌设置
// AuditDetails 包含频道的审核信息
// ContentOwnerDetails 包含内容所有者信息
// Localizations 包含频道的本地化信息
type ChannelItem struct {
	Kind                string                         `json:"kind"`
	Etag                string                         `json:"etag"`
	Id                  string                         `json:"id"`
	Snippet             ChannelSnippet                 `json:"snippet"`
	ContentDetails      ChannelContentDetails          `json:"contentDetails"`
	Statistics          ChannelStatistics              `json:"statistics"`
	TopicDetails        ChannelTopicDetails            `json:"topicDetails,omitempty"`
	Status              ChannelStatus                  `json:"status"`
	BrandingSettings    ChannelBrandingSettings        `json:"brandingSettings"`
	AuditDetails        AuditDetails                   `json:"auditDetails"`
	ContentOwnerDetails ChannelContentOwnerDetails     `json:"contentOwnerDetails,omitempty"`
	Localizations       map[string]ChannelLocalization `json:"localizations,omitempty"`
}

// ChannelSnippet 包含YouTube频道的基本信息
// Title 表示频道标题
// Description 表示频道描述
// CustomUrl 表示频道的自定义URL
// PublishedAt 表示频道的发布时间
// Thumbnails 包含频道的缩略图信息
// DefaultLanguage 表示频道的默认语言
// Localized 包含频道的本地化信息
// Country 表示频道的国家代码
type ChannelSnippet struct {
	Title           string                      `json:"title"`
	Description     string                      `json:"description"`
	CustomUrl       string                      `json:"customUrl"`
	PublishedAt     string                      `json:"publishedAt"`
	Thumbnails      map[string]ChannelThumbnail `json:"thumbnails"`
	DefaultLanguage string                      `json:"defaultLanguage,omitempty"`
	Localized       ChannelLocalization         `json:"localized"`
	Country         string                      `json:"country,omitempty"`
}

// ChannelThumbnail 表示YouTube频道的缩略图信息
// Url 表示缩略图的URL
// Width 表示缩略图的宽度
// Height 表示缩略图的高度
type ChannelThumbnail struct {
	Url    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}

// ChannelContentDetails 包含YouTube频道的内容详情
// RelatedPlaylists 包含相关播放列表信息
type ChannelContentDetails struct {
	RelatedPlaylists ChannelRelatedPlaylists `json:"relatedPlaylists"`
}

// ChannelRelatedPlaylists 包含YouTube频道相关播放列表信息
// Likes 表示点赞视频的播放列表
// Favorites 表示收藏视频的播放列表
// Uploads 表示上传视频的播放列表
type ChannelRelatedPlaylists struct {
	Likes     string `json:"likes"`
	Favorites string `json:"favorites"`
	Uploads   string `json:"uploads"`
}

// ChannelStatistics 包含YouTube频道的统计信息
// ViewCount 表示频道的总观看次数
// SubscriberCount 表示频道的订阅者数量
// HiddenSubscriberCount 表示是否隐藏订阅者数量
// VideoCount 表示频道的视频数量
type ChannelStatistics struct {
	ViewCount             string `json:"viewCount"`
	SubscriberCount       string `json:"subscriberCount"`
	HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
	VideoCount            string `json:"videoCount"`
}

// ChannelTopicDetails 包含YouTube频道的主题信息
// TopicIds 表示频道相关的主题ID列表
// TopicCategories 表示频道相关的主题类别列表
type ChannelTopicDetails struct {
	TopicIds        []string `json:"topicIds,omitempty"`
	TopicCategories []string `json:"topicCategories,omitempty"`
}

// ChannelStatus 包含YouTube频道的状态信息
// PrivacyStatus 表示频道的隐私状态
// IsLinked 表示频道是否已关联
// LongUploadsStatus 表示长视频上传状态
// MadeForKids 表示频道是否为儿童制作
// SelfDeclaredMadeForKids 表示频道是否声明为儿童制作
type ChannelStatus struct {
	PrivacyStatus           string `json:"privacyStatus"`
	IsLinked                bool   `json:"isLinked"`
	LongUploadsStatus       string `json:"longUploadsStatus"`
	MadeForKids             bool   `json:"madeForKids"`
	SelfDeclaredMadeForKids bool   `json:"selfDeclaredMadeForKids"`
}

// ChannelBrandingSettings 包含YouTube频道的品牌设置
// Channel 包含频道的基本品牌设置
// Watch 包含观看页面的品牌设置
type ChannelBrandingSettings struct {
	Channel BrandingChannel `json:"channel"`
	Watch   BrandingWatch   `json:"watch"`
}

// BrandingChannel 包含YouTube频道的基本品牌设置
// Title 表示频道的品牌标题
// Description 表示频道的品牌描述
// Keywords 表示频道的品牌关键词
// TrackingAnalyticsAccountId 表示跟踪分析账户ID
// UnsubscribedTrailer 表示未订阅用户的预告片
// DefaultLanguage 表示频道的默认语言
// Country 表示频道的国家代码
type BrandingChannel struct {
	Title                      string `json:"title"`
	Description                string `json:"description"`
	Keywords                   string `json:"keywords"`
	TrackingAnalyticsAccountId string `json:"trackingAnalyticsAccountId"`
	UnsubscribedTrailer        string `json:"unsubscribedTrailer"`
	DefaultLanguage            string `json:"defaultLanguage"`
	Country                    string `json:"country"`
}

// BrandingWatch 包含YouTube频道观看页面的品牌设置
// TextColor 表示文本颜色
// BackgroundColor 表示背景颜色
// FeaturedPlaylistId 表示特色播放列表ID
type BrandingWatch struct {
	TextColor          string `json:"textColor"`
	BackgroundColor    string `json:"backgroundColor"`
	FeaturedPlaylistId string `json:"featuredPlaylistId"`
}

// AuditDetails 包含YouTube频道的审核信息
// OverallGoodStanding 表示频道整体是否合规
// CommunityGuidelinesGoodStanding 表示频道是否遵守社区准则
// CopyrightStrikesGoodStanding 表示频道是否遵守版权规则
// ContentIdClaimsGoodStanding 表示频道是否遵守内容ID声明规则
type AuditDetails struct {
	OverallGoodStanding             bool `json:"overallGoodStanding"`
	CommunityGuidelinesGoodStanding bool `json:"communityGuidelinesGoodStanding"`
	CopyrightStrikesGoodStanding    bool `json:"copyrightStrikesGoodStanding"`
	ContentIdClaimsGoodStanding     bool `json:"contentIdClaimsGoodStanding"`
}

// ChannelContentOwnerDetails 包含YouTube频道的内容所有者信息
// ContentOwner 表示内容所有者
// TimeLinked 表示关联时间
type ChannelContentOwnerDetails struct {
	ContentOwner string `json:"contentOwner"`
	TimeLinked   string `json:"timeLinked"`
}

// ChannelLocalization 包含YouTube频道的本地化信息
// Title 表示本地化的标题
// Description 表示本地化的描述
type ChannelLocalization struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
