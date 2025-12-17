package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient"
	playlistsSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v4/playlists/schema"
	searchSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v4/search/schema"
	videoSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v4/video/schema"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// ExecuteAction routes a request to the matching YouTube API client and returns raw data.
func ExecuteAction(ctx context.Context, log *logger.Logger, yt *accessTokenClient.GoogleYouTubeACClient, opts *Options) (any, error) {
	if opts == nil {
		return nil, errors.New("options are required")
	}
	switch opts.Action {
	case "videos.list":
		return execVideosList(ctx, log, yt, opts)
	case "search.list":
		return execSearchList(ctx, log, yt, opts)
	case "playlists.list":
		return execPlaylistsList(ctx, log, yt, opts)
	default:
		return nil, fmt.Errorf("unsupported action %q", opts.Action)
	}
}

func execVideosList(ctx context.Context, log *logger.Logger, yt *accessTokenClient.GoogleYouTubeACClient, opts *Options) (*videoSchema.YouTubeVideoListRes, error) {
	if opts.Part == "" {
		return nil, errors.New("videos.list: part is required")
	}
	req := &videoSchema.YouTubeVideoListReq{
		Part:            opts.Part,
		Chart:           opts.Chart,
		ID:              opts.IDs,
		RegionCode:      opts.Region,
		VideoCategoryID: opts.VideoCategory,
		PageToken:       opts.PageToken,
	}
	req.MaxResults = fmt.Sprintf("%d", opts.MaxResults)
	if log != nil {
		log.InfoF(
			"accesstoken: request provider=%s action=%s ids=%s chart=%s region=%s category=%s max_results=%s page_token=%s",
			ValueOrDash(opts.ProviderCode),
			opts.Action,
			ValueOrDash(opts.IDs),
			ValueOrDash(opts.Chart),
			ValueOrDash(opts.Region),
			ValueOrDash(opts.VideoCategory),
			req.MaxResults,
			ValueOrDash(opts.PageToken),
		)
	}
	result, err := yt.GetVideoClient().List(ctx, req)
	if err != nil {
		logActionError(log, opts, err)
		return nil, fmt.Errorf("videos.list failed: %w", err)
	}
	logActionSuccess(log, opts, len(result.Items), result.NextPageToken)
	return result, nil
}

func execSearchList(ctx context.Context, log *logger.Logger, yt *accessTokenClient.GoogleYouTubeACClient, opts *Options) (*searchSchema.YouTubeSearchListRes, error) {
	if opts.Part == "" {
		return nil, errors.New("search.list: part is required")
	}
	req := &searchSchema.YouTubeSearchListReq{
		Part:       opts.Part,
		Q:          opts.Query,
		Type:       opts.SearchType,
		ChannelId:  opts.ChannelID,
		RegionCode: opts.Region,
		ForMine:    opts.SearchMine,
		MaxResults: opts.MaxResults,
		PageToken:  opts.PageToken,
	}
	if log != nil {
		log.InfoF(
			"accesstoken: request provider=%s action=%s query=%s search_type=%s search_mine=%t channel_id=%s region=%s max_results=%d page_token=%s",
			ValueOrDash(opts.ProviderCode),
			opts.Action,
			ValueOrDash(opts.Query),
			ValueOrDash(opts.SearchType),
			opts.SearchMine,
			ValueOrDash(opts.ChannelID),
			ValueOrDash(opts.Region),
			opts.MaxResults,
			ValueOrDash(opts.PageToken),
		)
	}
	res, err := yt.GetSearchClient().List(ctx, req)
	if err != nil {
		logActionError(log, opts, err)
		return nil, fmt.Errorf("search.list failed: %w", err)
	}
	logActionSuccess(log, opts, len(res.Items), res.NextPageToken)
	return res, nil
}

func execPlaylistsList(ctx context.Context, log *logger.Logger, yt *accessTokenClient.GoogleYouTubeACClient, opts *Options) (*playlistsSchema.YouTubePlaylistsListRes, error) {
	if opts.Part == "" {
		return nil, errors.New("playlists.list: part is required")
	}
	req := &playlistsSchema.YouTubePlaylistsListReq{
		Part:      opts.Part,
		ChannelId: opts.ChannelID,
		Id:        opts.IDs,
		Mine:      opts.Mine,
		PageToken: opts.PageToken,
	}
	if opts.MaxResults > 0 {
		req.MaxResults = opts.MaxResults
	}
	if log != nil {
		log.InfoF(
			"accesstoken: request provider=%s action=%s channel_id=%s ids=%s mine=%t max_results=%d page_token=%s",
			ValueOrDash(opts.ProviderCode),
			opts.Action,
			ValueOrDash(opts.ChannelID),
			ValueOrDash(opts.IDs),
			opts.Mine,
			req.MaxResults,
			ValueOrDash(opts.PageToken),
		)
	}
	res, err := yt.GetPlaylistsClient().List(ctx, req)
	if err != nil {
		logActionError(log, opts, err)
		return nil, fmt.Errorf("playlists.list failed: %w", err)
	}
	logActionSuccess(log, opts, len(res.Items), res.NextPageToken)
	return res, nil
}

func logActionSuccess(log *logger.Logger, opts *Options, count int, nextPage string) {
	if log == nil || opts == nil {
		return
	}
	log.InfoF(
		"accesstoken: success provider=%s action=%s part=%s results=%d next_page=%s",
		ValueOrDash(opts.ProviderCode),
		opts.Action,
		ValueOrDash(opts.Part),
		count,
		ValueOrDash(nextPage),
	)
}

func logActionError(log *logger.Logger, opts *Options, err error) {
	if log == nil || err == nil {
		return
	}
	action := "-"
	if opts != nil && opts.Action != "" {
		action = opts.Action
	}
	log.ErrorF("accesstoken: provider=%s action=%s error=%v", ValueOrDash(opts.ProviderCode), action, err)
}
