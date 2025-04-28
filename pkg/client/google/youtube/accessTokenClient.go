package youtube

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/activities"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/captions"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/channelBanners"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/channelSections"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/channels"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/commentThreads"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/comments"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/i18nLanguages"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/i18nRegions"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/members"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/membershipsLevels"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/playlistItems"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/playlists"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/search"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/subscriptions"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/thumbnails"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/video"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/videoAbuseReportReasons"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/videoCategory"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/watermarks"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// https://developers.google.com/youtube/v3/docs?hl=zh-cn
type GoogleYouTubeACClient struct {
	GoogleClient       *core.GoogleClient
	YouTubeConfig      *config.GoogleYouTubeConfig
	AccessTokenHandler *core.GoogleAccessTokenHandler

	// clients
	video                   *video.YoutubeVideoClient
	videoAbuseReportReasons *videoAbuseReportReasons.YoutubeVideoAbuseReportReasonsClient
	videoCategories         *videoCategory.YoutubeVideoCategoryClient
	activities              *activities.YoutubeActivitiesClient
	captions                *captions.YoutubeCaptionsClient
	channelBanners          *channelBanners.YoutubeChannelBannersClient
	channels                *channels.YoutubeChannelsClient
	comments                *comments.YoutubeCommentsClient
	commentThreads          *commentThreads.YoutubeCommentThreadsClient
	i18nLanguages           *i18nLanguages.YoutubeI18nLanguagesClient
	i18nRegions             *i18nRegions.YoutubeI18nRegionsClient
	members                 *members.YoutubeMembersClient
	membershipsLevels       *membershipsLevels.YoutubeMembershipsLevelsClient
	playlistItems           *playlistItems.YoutubePlaylistItemsClient
	playlists               *playlists.YoutubePlaylistsClient
	search                  *search.YoutubeSearchClient
	subscriptions           *subscriptions.YoutubeSubscriptionsClient
	thumbnails              *thumbnails.YoutubeThumbnailsClient
	channelSections         *channelSections.YoutubeChannelSectionsClient
	watermarks              *watermarks.YoutubeWatermarksClient
}

func NewGoogleYouTubeACClient(cfg *config.GoogleYouTubeConfig, logger *logger.Logger, cache cache.ICache) (*GoogleYouTubeACClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = config.GoogleYoutubeAPIUrl
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

func (client *GoogleYouTubeACClient) GetVideoClient() *video.YoutubeVideoClient {
	if client.video == nil {
		client.video = video.NewClient(client.GoogleClient.BaseClient)
	}
	return client.video
}

func (client *GoogleYouTubeACClient) GetVideoCategoriesClient() *videoCategory.YoutubeVideoCategoryClient {
	if client.videoCategories == nil {
		client.videoCategories = videoCategory.NewClient(client.GoogleClient.BaseClient)
	}
	return client.videoCategories
}

func (client *GoogleYouTubeACClient) GetActivitiesClient() *activities.YoutubeActivitiesClient {
	if client.activities == nil {
		client.activities = activities.NewClient(client.GoogleClient.BaseClient)
	}
	return client.activities
}

func (client *GoogleYouTubeACClient) GetCaptionsClient() *captions.YoutubeCaptionsClient {
	if client.captions == nil {
		client.captions = captions.NewClient(client.GoogleClient.BaseClient)
	}
	return client.captions
}

func (client *GoogleYouTubeACClient) GetChannelBannersClient() *channelBanners.YoutubeChannelBannersClient {
	if client.channelBanners == nil {
		client.channelBanners = channelBanners.NewClient(client.GoogleClient.BaseClient)
	}
	return client.channelBanners
}

func (client *GoogleYouTubeACClient) GetChannelsClient() *channels.YoutubeChannelsClient {
	if client.channels == nil {
		client.channels = channels.NewClient(client.GoogleClient.BaseClient)
	}
	return client.channels
}

func (client *GoogleYouTubeACClient) GetCommentsClient() *comments.YoutubeCommentsClient {
	if client.comments == nil {
		client.comments = comments.NewClient(client.GoogleClient.BaseClient)
	}
	return client.comments
}

func (client *GoogleYouTubeACClient) GetCommentThreadsClient() *commentThreads.YoutubeCommentThreadsClient {
	if client.commentThreads == nil {
		client.commentThreads = commentThreads.NewClient(client.GoogleClient.BaseClient)
	}
	return client.commentThreads
}

func (client *GoogleYouTubeACClient) GetI18nLanguagesClient() *i18nLanguages.YoutubeI18nLanguagesClient {
	if client.i18nLanguages == nil {
		client.i18nLanguages = i18nLanguages.NewClient(client.GoogleClient.BaseClient)
	}
	return client.i18nLanguages
}

func (client *GoogleYouTubeACClient) GetI18nRegionsClient() *i18nRegions.YoutubeI18nRegionsClient {
	if client.i18nRegions == nil {
		client.i18nRegions = i18nRegions.NewClient(client.GoogleClient.BaseClient)
	}
	return client.i18nRegions
}

func (client *GoogleYouTubeACClient) GetMembersClient() *members.YoutubeMembersClient {
	if client.members == nil {
		client.members = members.NewClient(client.GoogleClient.BaseClient)
	}
	return client.members
}

func (client *GoogleYouTubeACClient) GetMembershipsLevelsClient() *membershipsLevels.YoutubeMembershipsLevelsClient {
	if client.membershipsLevels == nil {
		client.membershipsLevels = membershipsLevels.NewClient(client.GoogleClient.BaseClient)
	}
	return client.membershipsLevels
}

func (client *GoogleYouTubeACClient) GetPlaylistItemsClient() *playlistItems.YoutubePlaylistItemsClient {
	if client.playlistItems == nil {
		client.playlistItems = playlistItems.NewClient(client.GoogleClient.BaseClient)
	}
	return client.playlistItems
}

func (client *GoogleYouTubeACClient) GetPlaylistsClient() *playlists.YoutubePlaylistsClient {
	if client.playlists == nil {
		client.playlists = playlists.NewClient(client.GoogleClient.BaseClient)
	}
	return client.playlists
}

func (client *GoogleYouTubeACClient) GetSearchClient() *search.YoutubeSearchClient {
	if client.search == nil {
		client.search = search.NewClient(client.GoogleClient.BaseClient)
	}
	return client.search
}

func (client *GoogleYouTubeACClient) GetSubscriptionsClient() *subscriptions.YoutubeSubscriptionsClient {
	if client.subscriptions == nil {
		client.subscriptions = subscriptions.NewClient(client.GoogleClient.BaseClient)
	}
	return client.subscriptions
}

func (client *GoogleYouTubeACClient) GetThumbnailsClient() *thumbnails.YoutubeThumbnailsClient {
	if client.thumbnails == nil {
		client.thumbnails = thumbnails.NewClient(client.GoogleClient.BaseClient)
	}
	return client.thumbnails
}

func (client *GoogleYouTubeACClient) GetChannelSectionsClient() *channelSections.YoutubeChannelSectionsClient {
	if client.channelSections == nil {
		client.channelSections = channelSections.NewClient(client.GoogleClient.BaseClient)
	}
	return client.channelSections
}

func (client *GoogleYouTubeACClient) GetWatermarksClient() *watermarks.YoutubeWatermarksClient {
	if client.watermarks == nil {
		client.watermarks = watermarks.NewClient(client.GoogleClient.BaseClient)
	}
	return client.watermarks
}
