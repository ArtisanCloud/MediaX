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

	app "github.com/ArtisanCloud/MediaX/cmd/accesstoken/internal/app"
	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

func main() {
	opts := parseFlags()
	if err := run(opts); err != nil {
		log.Fatalf("accesstoken: %v", err)
	}
}

func parseFlags() *app.Options {
	configPath := flag.String("config", envString("ACCESSTOKEN_CONFIG", ""), "Path to config.yaml (defaults to config.yaml or MEDIA_X_CONFIG env)")
	provider := flag.String("provider", envString("ACCESSTOKEN_PROVIDER", ""), "Provider group code (e.g. google)")
	providerApp := flag.String("provider-app", envString("ACCESSTOKEN_PROVIDER_APP", ""), "Provider app code (e.g. youtube)")
	authMode := flag.String("auth-mode", envString("ACCESSTOKEN_AUTH_MODE", ""), "Auth mode key (default: default)")
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

	return &app.Options{
		ConfigPath:     strings.TrimSpace(*configPath),
		Provider:       strings.TrimSpace(*provider),
		ProviderApp:    strings.TrimSpace(*providerApp),
		AuthMode:       strings.TrimSpace(*authMode),
		Action:         strings.TrimSpace(*action),
		Part:           strings.TrimSpace(*part),
		IDs:            strings.TrimSpace(*ids),
		Chart:          strings.TrimSpace(*chart),
		VideoCategory:  strings.TrimSpace(*category),
		Region:         strings.TrimSpace(*region),
		PageToken:      strings.TrimSpace(*pageToken),
		MaxResults:     *maxResults,
		Query:          strings.TrimSpace(*query),
		ChannelID:      strings.TrimSpace(*channelID),
		Mine:           *mine,
		SearchMine:     *searchMine,
		SearchType:     strings.TrimSpace(*searchType),
		AccessToken:    strings.TrimSpace(*accessToken),
		AccessTokenTTL: *accessTokenTTL,
	}
}

func run(opts *app.Options) error {
	if opts == nil {
		return errors.New("options are required")
	}
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		return err
	}

	configPath := app.ResolveConfigPath(opts.ConfigPath)
	localConfig := &config.LocalConfig{}
	if err := utils.LoadYAML(configPath, localConfig); err != nil {
		return fmt.Errorf("load config %s: %w", configPath, err)
	}
	if localConfig.AccessTokenProviders == nil {
		return errors.New("missing access_token_providers in config file")
	}
	provider, providerApp := localConfig.AccessTokenProviders.FindApp(opts.Provider, opts.ProviderApp)
	if providerApp == nil {
		provider, providerApp = localConfig.AccessTokenProviders.FindAppByProviderCode(opts.Provider)
	}
	var mode *config.AccessTokenAuthMode
	if providerApp != nil {
		mode = providerApp.FindMode(opts.AuthMode)
	}
	if providerApp == nil || mode == nil {
		var defaultProvider *config.AccessTokenProvider
		defaultProvider, providerApp, mode = localConfig.AccessTokenProviders.FirstSelection()
		if providerApp == nil || mode == nil {
			return errors.New("no access token provider/app configured in config file")
		}
		if provider == nil {
			provider = defaultProvider
		}
	}
	opts.Provider = provider.Code
	opts.ProviderApp = providerApp.Code
	if opts.AuthMode == "" && mode != nil {
		opts.AuthMode = mode.Key
	}
	opts.ProviderCode = mode.EffectiveProviderCode(providerApp)
	if opts.ProviderCode == "" {
		return fmt.Errorf("provider %s/%s missing provider_code", opts.Provider, opts.ProviderApp)
	}
	if mode.ConfigKind() != "google_youtube" || mode.GoogleYouTubeConfig == nil {
		return fmt.Errorf("CLI currently只支持 Google YouTube app，选择了 %s", mode.ConfigKind())
	}
	ytCfg := mode.GoogleYouTubeConfig

	accessToken, source := app.ResolveAccessToken(opts.AccessToken, opts.ProviderCode, ytCfg.ClientConfig)
	if accessToken == "" {
		return fmt.Errorf("missing access token: provide --access-token, 环境变量 %v 或配置文件 oauth.access_token", app.EnvKeysForProvider(opts.ProviderCode))
	}
	opts.AccessToken = accessToken
	opts.TokenSource = source

	ytCfg.GetOAuthToken = func(key string, refresh bool) object.HashMap {
		return object.HashMap{
			"access_token": accessToken,
			"expires_in":   float64(opts.AccessTokenTTL),
		}
	}

	cacheStore := cache.NewMemoryCache()
	mediaX := client.NewMediaX(&config.MediaXConfig{Logger: app.BuildBasicLogConfig("")}, cacheStore)
	app.LogInvocation(mediaX.Logger, opts, ytCfg.OauthKey)

	ytClient, err := mediaX.CreateGoogleYouTubeACClient(ytCfg)
	if err != nil {
		return fmt.Errorf("create youtube client: %w", err)
	}

	ctx := context.Background()
	data, err := app.ExecuteAction(ctx, mediaX.Logger, ytClient, opts)
	if err != nil {
		return err
	}
	return printJSON(data)
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}
	fmt.Println(string(data))
	return nil
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
