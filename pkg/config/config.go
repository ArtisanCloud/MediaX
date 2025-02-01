package config

// PluginMediaXPlatform 枚举定义
type PluginMediaXPlatform string

const (
	PluginMediaX PluginMediaXPlatform = "PluginMediaX"

	// WeChat
	WechatOfficialAccount PluginMediaXPlatform = "PluginMediaXWechatOfficialAccount"
	WechatMiniProgram     PluginMediaXPlatform = "PluginMediaXWechatMiniProgram"
	WechatMoments         PluginMediaXPlatform = "PluginMediaXWechatMoments"
	WechatVideo           PluginMediaXPlatform = "PluginMediaXWechatVideo"
	WechatLive            PluginMediaXPlatform = "PluginMediaXWechatLive"

	// YouTube
	YouTubeChannel  PluginMediaXPlatform = "PluginMediaXYouTubeChannel"
	YouTubeVideo    PluginMediaXPlatform = "PluginMediaXYouTubeVideo"
	YouTubePlaylist PluginMediaXPlatform = "PluginMediaXYouTubePlaylist"
	YouTubeLive     PluginMediaXPlatform = "PluginMediaXYouTubeLive"

	// Instagram
	InstagramFeed  PluginMediaXPlatform = "PluginMediaXInstagramFeed"
	InstagramStory PluginMediaXPlatform = "PluginMediaXInstagramStory"
	InstagramPost  PluginMediaXPlatform = "PluginMediaXInstagramPost"
	InstagramVideo PluginMediaXPlatform = "PluginMediaXInstagramVideo"

	// Facebook
	FacebookPage  PluginMediaXPlatform = "PluginMediaXFacebookPage"
	FacebookPost  PluginMediaXPlatform = "PluginMediaXFacebookPost"
	FacebookLive  PluginMediaXPlatform = "PluginMediaXFacebookLive"
	FacebookGroup PluginMediaXPlatform = "PluginMediaXFacebookGroup"

	// TikTok
	TikTokVideo PluginMediaXPlatform = "PluginMediaXTikTokVideo"
	TikTokLive  PluginMediaXPlatform = "PluginMediaXTikTokLive"
	TikTokDuet  PluginMediaXPlatform = "PluginMediaXTikTokDuet"

	// LinkedIn
	LinkedInPost    PluginMediaXPlatform = "PluginMediaXLinkedInPost"
	LinkedInArticle PluginMediaXPlatform = "PluginMediaXLinkedInArticle"
	LinkedInVideo   PluginMediaXPlatform = "PluginMediaXLinkedInVideo"

	// Snapchat
	SnapchatStory     PluginMediaXPlatform = "PluginMediaXSnapchatStory"
	SnapchatSpotlight PluginMediaXPlatform = "PluginMediaXSnapchatSpotlight"

	// Pinterest
	PinterestPin   PluginMediaXPlatform = "PluginMediaXPinterestPin"
	PinterestBoard PluginMediaXPlatform = "PluginMediaXPinterestBoard"

	// Vimeo
	VimeoVideo PluginMediaXPlatform = "PluginMediaXVimeoVideo"
	VimeoLive  PluginMediaXPlatform = "PluginMediaXVimeoLive"

	// Reddit
	RedditPost    PluginMediaXPlatform = "PluginMediaXRedditPost"
	RedditComment PluginMediaXPlatform = "PluginMediaXRedditComment"
)
