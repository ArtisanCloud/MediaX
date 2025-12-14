package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	ytclient "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient"
	playlistsSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/playlists/schema"
	searchSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/search/schema"
	videoSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video/schema"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	loggerconfig "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

const defaultConfigPath = "config.yaml"

func main() {
	opts := parseFlags()
	if err := run(opts); err != nil {
		log.Fatalf("accesstoken: %v", err)
	}
}

type options struct {
	configPath     string
	action         string
	part           string
	ids            string
	chart          string
	videoCategory  string
	region         string
	pageToken      string
	maxResults     int
	query          string
	channelID      string
	mine           bool
	searchMine     bool
	searchType     string
	accessToken    string
	accessTokenTTL int
}

func parseFlags() *options {
	configPath := flag.String("config", "", "Path to config.yaml (defaults to config.yaml or MEDIA_X_CONFIG env)")
	action := flag.String("action", "videos.list", "Action to execute: videos.list | search.list | playlists.list")
	part := flag.String("part", "snippet", "Value for the YouTube part parameter")
	ids := flag.String("ids", "", "Comma separated resource IDs (videos/playlists)")
	chart := flag.String("chart", "", "Chart parameter for videos.list (e.g. mostPopular)")
	category := flag.String("category", "", "videoCategoryId for videos.list")
	region := flag.String("region", "", "Region code (ISO 3166-1 alpha-2)")
	pageToken := flag.String("page-token", "", "Page token for pagination")
	maxResults := flag.Int("max-results", 5, "Max results (1-50)")
	query := flag.String("query", "", "Keyword for search.list")
	channelID := flag.String("channel-id", "", "Channel ID for search.list/playlists.list")
	mine := flag.Bool("mine", false, "Set mine=true for playlists.list")
	searchMine := flag.Bool("search-mine", false, "Set forMine=true for search.list")
	searchType := flag.String("search-type", "", "Search type for search.list (video,channel,playlist)")
	accessToken := flag.String("access-token", "", "YouTube OAuth access token (overrides config/env)")
	accessTokenTTL := flag.Int("access-token-ttl", 3600, "Access token TTL seconds for cache metadata")

	flag.Parse()

	return &options{
		configPath:     strings.TrimSpace(*configPath),
		action:         strings.TrimSpace(*action),
		part:           strings.TrimSpace(*part),
		ids:            strings.TrimSpace(*ids),
		chart:          strings.TrimSpace(*chart),
		videoCategory:  strings.TrimSpace(*category),
		region:         strings.TrimSpace(*region),
		pageToken:      strings.TrimSpace(*pageToken),
		maxResults:     *maxResults,
		query:          strings.TrimSpace(*query),
		channelID:      strings.TrimSpace(*channelID),
		mine:           *mine,
		searchMine:     *searchMine,
		searchType:     strings.TrimSpace(*searchType),
		accessToken:    strings.TrimSpace(*accessToken),
		accessTokenTTL: *accessTokenTTL,
	}
}

func run(opts *options) error {
	configPath := resolveConfigPath(opts.configPath)
	localConfig := &config.LocalConfig{}
	if err := utils.LoadYAML(configPath, localConfig); err != nil {
		return fmt.Errorf("load config %s: %w", configPath, err)
	}
	if localConfig.GoogleYouTubeConfig == nil {
		return errors.New("missing google_youtube_config in config file")
	}

	accessToken := resolveAccessToken(opts.accessToken, localConfig.GoogleYouTubeConfig)
	if accessToken == "" {
		return errors.New("missing access token: provide --access-token, GOOGLE_YOUTUBE_ACCESS_TOKEN env, or oauth.access_token in config")
	}

	if opts.accessTokenTTL <= 0 {
		opts.accessTokenTTL = 3600
	}

	localConfig.GoogleYouTubeConfig.GetOAuthToken = func(key string, refresh bool) object.HashMap {
		return object.HashMap{
			"access_token": accessToken,
			"expires_in":   float64(opts.accessTokenTTL),
		}
	}

	cacheStore := cache.NewMemoryCache()
	mediaX := client.NewMediaX(&config.MediaXConfig{Logger: buildLogConfig()}, cacheStore)

	ytClient, err := mediaX.CreateGoogleYouTubeACClient(localConfig.GoogleYouTubeConfig)
	if err != nil {
		return fmt.Errorf("create youtube client: %w", err)
	}

	ctx := context.Background()
	var data any
	switch opts.action {
	case "videos.list":
		data, err = execVideosList(ctx, ytClient, opts)
	case "search.list":
		data, err = execSearchList(ctx, ytClient, opts)
	case "playlists.list":
		data, err = execPlaylistsList(ctx, ytClient, opts)
	default:
		return fmt.Errorf("unsupported action %q", opts.action)
	}
	if err != nil {
		return err
	}
	return printJSON(data)
}

func execVideosList(ctx context.Context, yt *ytclient.GoogleYouTubeACClient, opts *options) (*videoSchema.YouTubeVideoListRes, error) {
	if opts.part == "" {
		return nil, errors.New("videos.list: part is required")
	}
	req := &videoSchema.YouTubeVideoListReq{
		Part:            opts.part,
		Chart:           opts.chart,
		ID:              opts.ids,
		RegionCode:      opts.region,
		VideoCategoryID: opts.videoCategory,
		PageToken:       opts.pageToken,
	}
	if opts.maxResults > 0 {
		req.MaxResults = fmt.Sprintf("%d", opts.maxResults)
	}
	result, err := yt.GetVideoClient().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("videos.list failed: %w", err)
	}
	return result, nil
}

func execSearchList(ctx context.Context, yt *ytclient.GoogleYouTubeACClient, opts *options) (*searchSchema.YouTubeSearchListRes, error) {
	if opts.part == "" {
		return nil, errors.New("search.list: part is required")
	}
	req := &searchSchema.YouTubeSearchListReq{
		Part:       opts.part,
		Q:          opts.query,
		Type:       opts.searchType,
		ChannelId:  opts.channelID,
		RegionCode: opts.region,
		ForMine:    opts.searchMine,
		MaxResults: opts.maxResults,
		PageToken:  opts.pageToken,
	}
	res, err := yt.GetSearchClient().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search.list failed: %w", err)
	}
	return res, nil
}

func execPlaylistsList(ctx context.Context, yt *ytclient.GoogleYouTubeACClient, opts *options) (*playlistsSchema.YouTubePlaylistsListRes, error) {
	if opts.part == "" {
		return nil, errors.New("playlists.list: part is required")
	}
	req := &playlistsSchema.YouTubePlaylistsListReq{
		Part:      opts.part,
		ChannelId: opts.channelID,
		Id:        opts.ids,
		Mine:      opts.mine,
		PageToken: opts.pageToken,
	}
	if opts.maxResults > 0 {
		req.MaxResults = opts.maxResults
	}
	res, err := yt.GetPlaylistsClient().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("playlists.list failed: %w", err)
	}
	return res, nil
}

func resolveConfigPath(flagPath string) string {
	if strings.TrimSpace(flagPath) != "" {
		return flagPath
	}
	if env := os.Getenv("MEDIA_X_CONFIG"); strings.TrimSpace(env) != "" {
		return env
	}
	return defaultConfigPath
}

func resolveAccessToken(flagToken string, cfg *config.GoogleYouTubeConfig) string {
	if strings.TrimSpace(flagToken) != "" {
		return strings.TrimSpace(flagToken)
	}
	if env := firstNonEmptyEnv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "YOUTUBE_ACCESS_TOKEN"); env != "" {
		return env
	}
	if cfg != nil && cfg.ClientConfig != nil && cfg.OAuthConfig != nil {
		if token := strings.TrimSpace(cfg.OAuthConfig.AccessToken); token != "" {
			return token
		}
	}
	return ""
}

func firstNonEmptyEnv(keys ...string) string {
	for _, k := range keys {
		if val := strings.TrimSpace(os.Getenv(k)); val != "" {
			return val
		}
	}
	return ""
}

func buildLogConfig() *loggerconfig.LogConfig {
	return &loggerconfig.LogConfig{
		Level:   "info",
		Console: true,
	}
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
