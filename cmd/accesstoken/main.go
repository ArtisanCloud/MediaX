package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	ytclient "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient"
	playlistsSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/playlists/schema"
	searchSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/search/schema"
	videoSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video/schema"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	loggerconfig "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

const (
	defaultConfigPath = "config.yaml"
	providerGoogle    = "google"
	maxResultsMin     = 1
	maxResultsMax     = 50
)

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
	tokenSource    string
}

func parseFlags() *options {
	configPath := flag.String("config", envString("ACCESSTOKEN_CONFIG", ""), "Path to config.yaml (defaults to config.yaml or MEDIA_X_CONFIG env)")
	action := flag.String("action", envString("ACCESSTOKEN_ACTION", ""), "Action to execute: videos.list | search.list | playlists.list")
	part := flag.String("part", envString("ACCESSTOKEN_PART", ""), "Value for the YouTube part parameter")
	ids := flag.String("ids", envString("ACCESSTOKEN_IDS", ""), "Comma separated resource IDs (videos/playlists)")
	chart := flag.String("chart", envString("ACCESSTOKEN_CHART", ""), "Chart parameter for videos.list (e.g. mostPopular)")
	category := flag.String("category", envString("ACCESSTOKEN_CATEGORY", ""), "videoCategoryId for videos.list")
	region := flag.String("region", envString("ACCESSTOKEN_REGION", ""), "Region code (ISO 3166-1 alpha-2)")
	pageToken := flag.String("page-token", envString("ACCESSTOKEN_PAGE_TOKEN", ""), "Page token for pagination")
	maxResults := flag.Int("max-results", envInt("ACCESSTOKEN_MAX_RESULTS", 5), "Max results (1-50)")
	query := flag.String("query", envString("ACCESSTOKEN_QUERY", ""), "Keyword for search.list")
	channelID := flag.String("channel-id", envString("ACCESSTOKEN_CHANNEL_ID", ""), "Channel ID for search.list/playlists.list")
	mine := flag.Bool("mine", envBool("ACCESSTOKEN_MINE", false), "Set mine=true for playlists.list")
	searchMine := flag.Bool("search-mine", envBool("ACCESSTOKEN_SEARCH_MINE", false), "Set forMine=true for search.list")
	searchType := flag.String("search-type", envString("ACCESSTOKEN_SEARCH_TYPE", ""), "Search type for search.list (video,channel,playlist)")
	accessToken := flag.String("access-token", envString("ACCESSTOKEN_ACCESS_TOKEN", ""), "YouTube OAuth access token (overrides config/env)")
	accessTokenTTL := flag.Int("access-token-ttl", envInt("ACCESSTOKEN_ACCESS_TOKEN_TTL", 3600), "Access token TTL seconds for cache metadata")

	flag.Parse()

	return &options{
		configPath:     strings.TrimSpace(*configPath),
		action:         strings.ToLower(strings.TrimSpace(*action)),
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

func validateOptions(opts *options) error {
	if opts == nil {
		return errors.New("options are required")
	}
	if opts.action == "" {
		return errors.New("missing required flag: -action")
	}
	switch opts.action {
	case "videos.list", "search.list", "playlists.list":
	default:
		return fmt.Errorf("unsupported action %q", opts.action)
	}
	if strings.TrimSpace(opts.part) == "" {
		return errors.New("missing required flag: -part")
	}
	if opts.maxResults < maxResultsMin || opts.maxResults > maxResultsMax {
		return fmt.Errorf("invalid -max-results %d: must be between %d and %d", opts.maxResults, maxResultsMin, maxResultsMax)
	}

	switch opts.action {
	case "videos.list":
		if opts.ids == "" && opts.chart == "" {
			return errors.New("videos.list: provide at least -ids or -chart")
		}
	case "search.list":
		if opts.query == "" && opts.channelID == "" && !opts.searchMine {
			return errors.New("search.list: provide -query, -channel-id, or -search-mine")
		}
	case "playlists.list":
		if opts.ids == "" && opts.channelID == "" && !opts.mine {
			return errors.New("playlists.list: provide at least -ids, -channel-id, or -mine")
		}
		if opts.mine && opts.ids != "" {
			return errors.New("playlists.list: -mine and -ids are mutually exclusive")
		}
	}

	return nil
}

func run(opts *options) error {
	if err := validateOptions(opts); err != nil {
		return err
	}

	configPath := resolveConfigPath(opts.configPath)
	localConfig := &config.LocalConfig{}
	if err := utils.LoadYAML(configPath, localConfig); err != nil {
		return fmt.Errorf("load config %s: %w", configPath, err)
	}
	if localConfig.GoogleYouTubeConfig == nil {
		return errors.New("missing google_youtube_config in config file")
	}

	accessToken, source := resolveAccessToken(opts.accessToken, localConfig.GoogleYouTubeConfig)
	if accessToken == "" {
		return errors.New("missing access token: provide --access-token, GOOGLE_YOUTUBE_ACCESS_TOKEN env, or oauth.access_token in config")
	}
	opts.accessToken = accessToken
	opts.tokenSource = source

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
	logInvocation(mediaX.Logger, opts, localConfig.GoogleYouTubeConfig.OauthKey)

	ytClient, err := mediaX.CreateGoogleYouTubeACClient(localConfig.GoogleYouTubeConfig)
	if err != nil {
		return fmt.Errorf("create youtube client: %w", err)
	}

	ctx := context.Background()
	var data any
	switch opts.action {
	case "videos.list":
		data, err = execVideosList(ctx, mediaX.Logger, ytClient, opts)
	case "search.list":
		data, err = execSearchList(ctx, mediaX.Logger, ytClient, opts)
	case "playlists.list":
		data, err = execPlaylistsList(ctx, mediaX.Logger, ytClient, opts)
	default:
		return fmt.Errorf("unsupported action %q", opts.action)
	}
	if err != nil {
		return err
	}
	return printJSON(data)
}

func execVideosList(ctx context.Context, log *logger.Logger, yt *ytclient.GoogleYouTubeACClient, opts *options) (*videoSchema.YouTubeVideoListRes, error) {
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
	req.MaxResults = fmt.Sprintf("%d", opts.maxResults)
	if log != nil {
		log.InfoF(
			"accesstoken: request provider=%s action=%s ids=%s chart=%s region=%s category=%s max_results=%s page_token=%s",
			providerGoogle,
			opts.action,
			valueOrDash(opts.ids),
			valueOrDash(opts.chart),
			valueOrDash(opts.region),
			valueOrDash(opts.videoCategory),
			req.MaxResults,
			valueOrDash(opts.pageToken),
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

func execSearchList(ctx context.Context, log *logger.Logger, yt *ytclient.GoogleYouTubeACClient, opts *options) (*searchSchema.YouTubeSearchListRes, error) {
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
	if log != nil {
		log.InfoF(
			"accesstoken: request provider=%s action=%s query=%s search_type=%s search_mine=%t channel_id=%s region=%s max_results=%d page_token=%s",
			providerGoogle,
			opts.action,
			valueOrDash(opts.query),
			valueOrDash(opts.searchType),
			opts.searchMine,
			valueOrDash(opts.channelID),
			valueOrDash(opts.region),
			opts.maxResults,
			valueOrDash(opts.pageToken),
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

func execPlaylistsList(ctx context.Context, log *logger.Logger, yt *ytclient.GoogleYouTubeACClient, opts *options) (*playlistsSchema.YouTubePlaylistsListRes, error) {
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
	if log != nil {
		log.InfoF(
			"accesstoken: request provider=%s action=%s channel_id=%s ids=%s mine=%t max_results=%d page_token=%s",
			providerGoogle,
			opts.action,
			valueOrDash(opts.channelID),
			valueOrDash(opts.ids),
			opts.mine,
			req.MaxResults,
			valueOrDash(opts.pageToken),
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

func resolveConfigPath(flagPath string) string {
	if strings.TrimSpace(flagPath) != "" {
		return flagPath
	}
	if env := os.Getenv("MEDIA_X_CONFIG"); strings.TrimSpace(env) != "" {
		return env
	}
	return defaultConfigPath
}

func resolveAccessToken(flagToken string, cfg *config.GoogleYouTubeConfig) (string, string) {
	if token := strings.TrimSpace(flagToken); token != "" {
		return token, "flag"
	}
	if env := firstNonEmptyEnv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "YOUTUBE_ACCESS_TOKEN"); env != "" {
		return env, "env"
	}
	if cfg != nil && cfg.ClientConfig != nil && cfg.OAuthConfig != nil {
		if token := strings.TrimSpace(cfg.OAuthConfig.AccessToken); token != "" {
			return token, "config"
		}
	}
	return "", ""
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

func logInvocation(log *logger.Logger, opts *options, oauthKey string) {
	if log == nil || opts == nil {
		return
	}
	log.InfoF(
		"accesstoken: provider=%s action=%s part=%s ids=%s channel_id=%s query=%s mine=%t search_mine=%t search_type=%s max_results=%d region=%s page_token=%s chart=%s category=%s token_source=%s token=%s oauth_key=%s",
		providerGoogle,
		opts.action,
		valueOrDash(opts.part),
		valueOrDash(opts.ids),
		valueOrDash(opts.channelID),
		valueOrDash(opts.query),
		opts.mine,
		opts.searchMine,
		valueOrDash(opts.searchType),
		opts.maxResults,
		valueOrDash(opts.region),
		valueOrDash(opts.pageToken),
		valueOrDash(opts.chart),
		valueOrDash(opts.videoCategory),
		valueOrDash(opts.tokenSource),
		maskToken(opts.accessToken),
		valueOrDash(oauthKey),
	)
}

func logActionSuccess(log *logger.Logger, opts *options, count int, nextPage string) {
	if log == nil || opts == nil {
		return
	}
	log.InfoF(
		"accesstoken: success provider=%s action=%s part=%s results=%d next_page=%s",
		providerGoogle,
		opts.action,
		valueOrDash(opts.part),
		count,
		valueOrDash(nextPage),
	)
}

func logActionError(log *logger.Logger, opts *options, err error) {
	if log == nil || err == nil {
		return
	}
	action := "-"
	if opts != nil && opts.action != "" {
		action = opts.action
	}
	log.ErrorF("accesstoken: provider=%s action=%s error=%v", providerGoogle, action, err)
}

func valueOrDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}

func maskToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return "-"
	}
	if len(token) <= 4 {
		return "***"
	}
	return token[:2] + "***" + token[len(token)-2:]
}

func envString(key, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return fallback
}

func envInt(key string, fallback int) int {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	if parsed, err := strconv.Atoi(val); err == nil {
		return parsed
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	switch strings.ToLower(val) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}
