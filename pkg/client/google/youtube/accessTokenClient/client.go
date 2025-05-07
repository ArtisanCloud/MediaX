package accessTokenClient

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/activities"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/captions"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/channelBanners"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/channelSections"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/channels"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/commentThreads"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/comments"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/i18nLanguages"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/i18nRegions"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/members"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/membershipsLevels"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/playlistImages"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/playlistItems"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/playlists"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/search"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/subscriptions"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/thumbnails"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/videoAbuseReportReasons"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/videoCategory"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/watermarks"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// GoogleYouTubeACClient YouTube访问Token客户端
// 提供YouTube API的各种功能接口，使用访问Token进行认证
// 文档：https://developers.google.com/youtube/v3/docs?hl=zh-cn
type GoogleYouTubeACClient struct {
	GoogleClient       *core.GoogleClient             // Google基础客户端
	YouTubeConfig      *config.GoogleYouTubeConfig    // YouTube配置
	AccessTokenHandler *core.GoogleAccessTokenHandler // 访问Token处理器

	// clients
	video                   *video.YoutubeVideoClient                                     // 视频管理客户端
	videoAbuseReportReasons *videoAbuseReportReasons.YoutubeVideoAbuseReportReasonsClient // 视频举报原因管理客户端
	videoCategories         *videoCategory.YoutubeVideoCategoryClient                     // 视频分类管理客户端
	activities              *activities.YoutubeActivitiesClient                           // 活动管理客户端
	captions                *captions.YoutubeCaptionsClient                               // 字幕管理客户端
	channelBanners          *channelBanners.YoutubeChannelBannersClient                   // 频道横幅管理客户端
	channels                *channels.YoutubeChannelsClient                               // 频道管理客户端
	comments                *comments.YoutubeCommentsClient                               // 评论管理客户端
	commentThreads          *commentThreads.YoutubeCommentThreadsClient                   // 评论线程管理客户端
	i18nLanguages           *i18nLanguages.YoutubeI18nLanguagesClient                     // 国际化语言管理客户端
	i18nRegions             *i18nRegions.YoutubeI18nRegionsClient                         // 国际化地区管理客户端
	members                 *members.YoutubeMembersClient                                 // 会员管理客户端
	membershipsLevels       *membershipsLevels.YoutubeMembershipsLevelsClient             // 会员等级管理客户端
	playlistImages          *playlistImages.YoutubePlaylistImagesClient                   // 播放列表图像管理客户端
	playlistItems           *playlistItems.YoutubePlaylistItemsClient                     // 播放列表项管理客户端
	playlists               *playlists.YoutubePlaylistsClient                             // 播放列表管理客户端
	search                  *search.YoutubeSearchClient                                   // 搜索管理客户端
	subscriptions           *subscriptions.YoutubeSubscriptionsClient                     // 订阅管理客户端
	thumbnails              *thumbnails.YoutubeThumbnailsClient                           // 缩略图管理客户端
	channelSections         *channelSections.YoutubeChannelSectionsClient                 // 频道分区管理客户端
	watermarks              *watermarks.YoutubeWatermarksClient                           // 水印管理客户端
}

// NewGoogleYouTubeACClient 创建新的YouTube访问Token客户端实例
// cfg: YouTube配置
// logger: 日志记录器
// cache: 缓存接口
func NewGoogleYouTubeACClient(cfg *config.GoogleYouTubeConfig, logger *logger.Logger, cache cache.ICache) (*GoogleYouTubeACClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = config.GoogleAppAPIUrl
	}
	c, err := core.NewGoogleClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewGoogleAccessTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.AccessTokenHandler

	// override get custom token
	c.TokenHandler.GetCustomToken = cfg.GetOAuthToken

	return &GoogleYouTubeACClient{
		GoogleClient:       c,
		YouTubeConfig:      cfg,
		AccessTokenHandler: handler,
	}, nil
}

// GetVideoClient 获取视频管理客户端
func (client *GoogleYouTubeACClient) GetVideoClient() *video.YoutubeVideoClient {
	if client.video == nil {
		client.video = video.NewClient(client.GoogleClient.BaseClient)
	}
	return client.video
}

// GetVideoCategoriesClient 获取视频分类管理客户端
func (client *GoogleYouTubeACClient) GetVideoCategoriesClient() *videoCategory.YoutubeVideoCategoryClient {
	if client.videoCategories == nil {
		client.videoCategories = videoCategory.NewClient(client.GoogleClient.BaseClient)
	}
	return client.videoCategories
}

// GetActivitiesClient 获取活动管理客户端
func (client *GoogleYouTubeACClient) GetActivitiesClient() *activities.YoutubeActivitiesClient {
	if client.activities == nil {
		client.activities = activities.NewClient(client.GoogleClient.BaseClient)
	}
	return client.activities
}

// GetCaptionsClient 获取字幕管理客户端
func (client *GoogleYouTubeACClient) GetCaptionsClient() *captions.YoutubeCaptionsClient {
	if client.captions == nil {
		client.captions = captions.NewClient(client.GoogleClient.BaseClient)
	}
	return client.captions
}

// GetChannelBannersClient 获取频道横幅管理客户端
func (client *GoogleYouTubeACClient) GetChannelBannersClient() *channelBanners.YoutubeChannelBannersClient {
	if client.channelBanners == nil {
		client.channelBanners = channelBanners.NewClient(client.GoogleClient.BaseClient)
	}
	return client.channelBanners
}

// GetChannelsClient 获取频道管理客户端
func (client *GoogleYouTubeACClient) GetChannelsClient() *channels.YoutubeChannelsClient {
	if client.channels == nil {
		client.channels = channels.NewClient(client.GoogleClient.BaseClient)
	}
	return client.channels
}

// GetCommentsClient 获取评论管理客户端
func (client *GoogleYouTubeACClient) GetCommentsClient() *comments.YoutubeCommentsClient {
	if client.comments == nil {
		client.comments = comments.NewClient(client.GoogleClient.BaseClient)
	}
	return client.comments
}

// GetCommentThreadsClient 获取评论线程管理客户端
func (client *GoogleYouTubeACClient) GetCommentThreadsClient() *commentThreads.YoutubeCommentThreadsClient {
	if client.commentThreads == nil {
		client.commentThreads = commentThreads.NewClient(client.GoogleClient.BaseClient)
	}
	return client.commentThreads
}

// GetI18nLanguagesClient 获取国际化语言管理客户端
func (client *GoogleYouTubeACClient) GetI18nLanguagesClient() *i18nLanguages.YoutubeI18nLanguagesClient {
	if client.i18nLanguages == nil {
		client.i18nLanguages = i18nLanguages.NewClient(client.GoogleClient.BaseClient)
	}
	return client.i18nLanguages
}

// GetI18nRegionsClient 获取国际化地区管理客户端
func (client *GoogleYouTubeACClient) GetI18nRegionsClient() *i18nRegions.YoutubeI18nRegionsClient {
	if client.i18nRegions == nil {
		client.i18nRegions = i18nRegions.NewClient(client.GoogleClient.BaseClient)
	}
	return client.i18nRegions
}

// GetMembersClient 获取会员管理客户端
func (client *GoogleYouTubeACClient) GetMembersClient() *members.YoutubeMembersClient {
	if client.members == nil {
		client.members = members.NewClient(client.GoogleClient.BaseClient)
	}
	return client.members
}

// GetMembershipsLevelsClient 获取会员等级管理客户端
func (client *GoogleYouTubeACClient) GetMembershipsLevelsClient() *membershipsLevels.YoutubeMembershipsLevelsClient {
	if client.membershipsLevels == nil {
		client.membershipsLevels = membershipsLevels.NewClient(client.GoogleClient.BaseClient)
	}
	return client.membershipsLevels
}

func (client *GoogleYouTubeACClient) GetPlaylistImagesClient() *playlistImages.YoutubePlaylistImagesClient {
	if client.playlistImages == nil {
		client.playlistImages = playlistImages.NewClient(client.GoogleClient.BaseClient)
	}
	return client.playlistImages
}

// GetPlaylistItemsClient 获取播放列表项管理客户端
func (client *GoogleYouTubeACClient) GetPlaylistItemsClient() *playlistItems.YoutubePlaylistItemsClient {
	if client.playlistItems == nil {
		client.playlistItems = playlistItems.NewClient(client.GoogleClient.BaseClient)
	}
	return client.playlistItems
}

// GetPlaylistsClient 获取播放列表管理客户端
func (client *GoogleYouTubeACClient) GetPlaylistsClient() *playlists.YoutubePlaylistsClient {
	if client.playlists == nil {
		client.playlists = playlists.NewClient(client.GoogleClient.BaseClient)
	}
	return client.playlists
}

// GetSearchClient 获取搜索管理客户端
func (client *GoogleYouTubeACClient) GetSearchClient() *search.YoutubeSearchClient {
	if client.search == nil {
		client.search = search.NewClient(client.GoogleClient.BaseClient)
	}
	return client.search
}

// GetSubscriptionsClient 获取订阅管理客户端
func (client *GoogleYouTubeACClient) GetSubscriptionsClient() *subscriptions.YoutubeSubscriptionsClient {
	if client.subscriptions == nil {
		client.subscriptions = subscriptions.NewClient(client.GoogleClient.BaseClient)
	}
	return client.subscriptions
}

// GetThumbnailsClient 获取缩略图管理客户端
func (client *GoogleYouTubeACClient) GetThumbnailsClient() *thumbnails.YoutubeThumbnailsClient {
	if client.thumbnails == nil {
		client.thumbnails = thumbnails.NewClient(client.GoogleClient.BaseClient)
	}
	return client.thumbnails
}

// GetChannelSectionsClient 获取频道分区管理客户端
func (client *GoogleYouTubeACClient) GetChannelSectionsClient() *channelSections.YoutubeChannelSectionsClient {
	if client.channelSections == nil {
		client.channelSections = channelSections.NewClient(client.GoogleClient.BaseClient)
	}
	return client.channelSections
}

// GetWatermarksClient 获取水印管理客户端
func (client *GoogleYouTubeACClient) GetWatermarksClient() *watermarks.YoutubeWatermarksClient {
	if client.watermarks == nil {
		client.watermarks = watermarks.NewClient(client.GoogleClient.BaseClient)
	}
	return client.watermarks
}
