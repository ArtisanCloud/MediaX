package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	app "github.com/ArtisanCloud/MediaX/cmd/accesstoken/internal/app"
	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type cacheCloser func()

type providerAppMeta struct {
	Code         string                `json:"code"`
	Name         string                `json:"name"`
	ProviderCode string                `json:"provider_code"`
	AppKey       string                `json:"app_key,omitempty"`
	ConfigPath   string                `json:"config_path"`
	APIVersion   string                `json:"api_version,omitempty"`
	Modes        []providerAppModeMeta `json:"modes"`
}

type providerAppModeMeta struct {
	Key          string `json:"key"`
	Label        string `json:"label"`
	ProviderCode string `json:"provider_code"`
	OauthKey     string `json:"oauth_key,omitempty"`
}

const (
	providerGoogleYouTube   = "google_youtube"
	providerGoogleBlogger   = "google_blogger"
	providerByteDanceDouYin = "byte_dance_douyin"
	providerRedBookJuGuang  = "redbook_juguang"
	providerBilibili        = "bilbili"
	flowIndexPrefix         = "accesstoken:oauth:flow:"

	storageBackendRedis  = "redis"
	storageBackendMemory = "memory"

	defaultFlowTTLSeconds = 24 * 60 * 60
	maxCallbackBodyLen    = 4096
)

var (
	errFlowNotFound = errors.New("flow_not_found")
	errFlowExpired  = errors.New("flow_expired")
)

type providerContext struct {
	GroupCode    string
	GroupName    string
	AppCode      string
	AppName      string
	ModeKey      string
	ProviderCode string
	ConfigPath   string
	APIVersion   string
	Mode         *config.AccessTokenAuthMode
	App          *config.AccessTokenProviderApp
	Provider     *config.AccessTokenProvider
	Youtube      *config.GoogleYouTubeConfig
	Blogger      *config.GoogleBloggerConfig
	Douyin       *config.ByteDanceDouYinConfig
	RedBook      *config.RedBookJuGuangConfig
	Bilibili     *config.BiliBiliConfig
}

type providerMeta struct {
	Code string            `json:"code"`
	Name string            `json:"name"`
	Apps []providerAppMeta `json:"apps"`
}

type oauthState struct {
	ProviderCode string
	ProviderApp  string
	ModeKey      string
	ConfigPath   string
	CreatedAt    time.Time
}

type oauthTokenRecord struct {
	ProviderCode   string          `json:"provider_code"`
	ProviderApp    string          `json:"provider_app"`
	AuthMode       string          `json:"provider_auth_mode"`
	ConfigPath     string          `json:"config_path,omitempty"`
	FlowID         string          `json:"flow_id"`
	AccessToken    string          `json:"access_token"`
	RefreshToken   string          `json:"refresh_token,omitempty"`
	TokenType      string          `json:"token_type,omitempty"`
	Scope          string          `json:"scope,omitempty"`
	ExpiresIn      int             `json:"expires_in"`
	StoredAt       time.Time       `json:"stored_at"`
	ExpireAt       time.Time       `json:"expire_at"`
	Source         string          `json:"source"`
	TokenExpireAt  time.Time       `json:"token_expire_at"`
	FlowTTLSeconds int             `json:"flow_ttl_seconds"`
	FlowExpireAt   time.Time       `json:"flow_expire_at"`
	StorageBackend string          `json:"storage_backend"`
	Callback       *callbackRecord `json:"callback,omitempty"`
}

type accessTokenServer struct {
	logger              *logger.Logger
	mediaX              *client.MediaX
	cache               cache.ICache
	redis               redis.Cmdable
	defaultConfigPath   string
	apiToken            string
	callbackStore       *callbackLogStore
	listenAddr          string
	defaultCallbackURL  string
	providers           []providerMeta
	defaultProviderCode string
	defaultProviderUI   string
	defaultApp          string
	defaultMode         string
	oauthStates         map[string]*oauthState
	oauthStateMu        sync.Mutex
	oauthTokens         map[string]*oauthTokenRecord
	oauthTokenMu        sync.RWMutex
	providersJSON       template.JS
	storageBackend      string
	flowTTLSeconds      int
	flowTTL             time.Duration
	listenAddrPublic    bool
}

func newAccessTokenServer(defaultConfigPath, listenAddr string) (*accessTokenServer, cacheCloser, error) {
	logCfg := app.BuildFileLogConfig(resolveLogLevel(), "logs/accesstoken-server-info.log", "logs/accesstoken-server-error.log")
	redisCfg := loadAccessTokenRedisConfig(defaultConfigPath)
	cacheStore, redisClient, closer, backend, err := buildCacheStore(redisCfg)
	if err != nil {
		return nil, nil, err
	}
	mediaX := client.NewMediaX(&config.MediaXConfig{Logger: logCfg}, cacheStore)
	flowTTLSeconds := resolveFlowTTLSeconds()
	if flowTTLSeconds <= 0 {
		flowTTLSeconds = defaultFlowTTLSeconds
	}
	server := &accessTokenServer{
		logger:            mediaX.Logger,
		mediaX:            mediaX,
		cache:             cacheStore,
		redis:             redisClient,
		defaultConfigPath: defaultConfigPath,
		apiToken:          resolveAPIToken(),
		callbackStore:     newCallbackLogStore(50),
		listenAddr:        listenAddr,
		oauthStates:       make(map[string]*oauthState),
		oauthTokens:       make(map[string]*oauthTokenRecord),
		storageBackend:    backend,
		flowTTLSeconds:    flowTTLSeconds,
		flowTTL:           time.Duration(flowTTLSeconds) * time.Second,
	}
	if backend == storageBackendMemory && server.logger != nil {
		server.logger.WarnF("accesstoken-server: redis disabled, Flow 数据将存储在内存中，服务重启后会被清理")
	}
	server.defaultCallbackURL = server.buildDefaultCallbackURL()
	if err := server.initProviders(); err != nil {
		return nil, nil, err
	}
	return server, closer, nil
}

func (s *accessTokenServer) routes() http.Handler {
	mux := http.NewServeMux()
	pageHandler := http.HandlerFunc(s.handleDebugPage)
	mux.Handle("/debug/accesstoken", pageHandler)
	mux.Handle("/debug", pageHandler)
	mux.Handle("/debug/callback", http.HandlerFunc(s.handleDebugCallback))
	mux.Handle("/api/callbacks", s.requireAPIToken(http.HandlerFunc(s.handleListCallbacks)))
	mux.Handle("/api/callbacks/clear", s.requireAPIToken(http.HandlerFunc(s.handleClearCallbacks)))
	mux.Handle("/debug/static/", http.StripPrefix("/debug/static/", s.staticFileHandler()))
	mux.Handle("/api/oauth/flow-indexes", s.requireAPIToken(http.HandlerFunc(s.handleListFlowIndexes)))
	oauthStart := s.requireAPIToken(http.HandlerFunc(s.handleOAuthStart))
	mux.Handle("/api/oauth/start", oauthStart)
	mux.Handle("/accesstoken/oauth/start", oauthStart)
	flowList := s.requireAPIToken(http.HandlerFunc(s.handleListFlows))
	mux.Handle("/api/oauth/tokens", flowList)
	mux.Handle("/accesstoken/flows", flowList)
	mux.Handle("/accesstoken/flow/replay", s.requireAPIToken(http.HandlerFunc(s.handleFlowReplay)))
	mux.Handle("/accesstoken/token", s.requireAPIToken(http.HandlerFunc(s.handleToken)))
	mux.Handle("/accesstoken/call", s.requireAPIToken(http.HandlerFunc(s.handleCall)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func buildCacheStore(redisCfg *config.AccessTokenRedisConfig) (cache.ICache, redis.Cmdable, cacheCloser, string, error) {
	addr := strings.TrimSpace(os.Getenv("ACCESSTOKEN_REDIS_ADDR"))
	if addr == "" && redisCfg != nil {
		addr = strings.TrimSpace(redisCfg.Addr)
	}
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	if strings.EqualFold(addr, "memory") {
		return cache.NewMemoryCache(), nil, func() {}, storageBackendMemory, nil
	}

	db := 0
	if v := strings.TrimSpace(os.Getenv("ACCESSTOKEN_REDIS_DB")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			db = n
		}
	} else if redisCfg != nil && redisCfg.DB >= 0 {
		db = redisCfg.DB
	}
	username := strings.TrimSpace(os.Getenv("ACCESSTOKEN_REDIS_USERNAME"))
	if username == "" && redisCfg != nil {
		username = strings.TrimSpace(redisCfg.Username)
	}
	password := strings.TrimSpace(os.Getenv("ACCESSTOKEN_REDIS_PASSWORD"))
	if password == "" && redisCfg != nil {
		password = strings.TrimSpace(redisCfg.Password)
	}
	opts := &redis.Options{
		Addr:     addr,
		DB:       db,
		Username: username,
		Password: password,
	}
	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		log.Printf("accesstoken-server: ping redis %s failed: %v, fallback to in-memory cache", addr, err)
		return cache.NewMemoryCache(), nil, func() {}, storageBackendMemory, nil
	}
	return cache.NewRedisCache(client), client, func() { _ = client.Close() }, storageBackendRedis, nil
}

func resolveListenAddr(flagListen, flagPort string) string {
	normalize := func(val string) string {
		val = strings.TrimSpace(val)
		if val == "" {
			return val
		}
		if strings.HasPrefix(val, ":") || strings.Contains(val, ":") {
			return val
		}
		if _, err := strconv.Atoi(val); err == nil {
			return defaultListenHost + ":" + val
		}
		return val
	}
	if trimmed := normalize(flagListen); trimmed != "" {
		return trimmed
	}
	if trimmed := normalize(flagPort); trimmed != "" {
		return trimmed
	}
	if env := strings.TrimSpace(os.Getenv("ACCESSTOKEN_LISTEN_ADDR")); env != "" {
		return normalize(env)
	}
	return defaultListenAddr
}

func resolveLogLevel() string {
	if env := strings.TrimSpace(os.Getenv("ACCESSTOKEN_LOG_LEVEL")); env != "" {
		return env
	}
	return "info"
}

func resolveAPIToken() string {
	if env := strings.TrimSpace(os.Getenv("ACCESSTOKEN_API_TOKEN")); env != "" {
		return env
	}
	return "dev-accesstoken"
}

func resolveFlowTTLSeconds() int {
	if env := strings.TrimSpace(os.Getenv("ACCESSTOKEN_FLOW_TTL_SECONDS")); env != "" {
		if v, err := strconv.Atoi(env); err == nil && v > 0 {
			return v
		}
	}
	return defaultFlowTTLSeconds
}

func (s *accessTokenServer) buildDefaultCallbackURL() string {
	addr := strings.TrimSpace(s.listenAddr)
	if addr == "" {
		addr = defaultListenAddr
	}
	host := "127.0.0.1"
	port := ":7071"
	switch {
	case strings.HasPrefix(addr, ":"):
		port = addr
	case strings.Contains(addr, ":"):
		parts := strings.Split(addr, ":")
		if len(parts) >= 2 {
			if candidate := strings.TrimSpace(parts[0]); candidate != "" {
				host = candidate
			}
			port = ":" + strings.TrimSpace(parts[len(parts)-1])
		}
	default:
		point := strings.TrimSpace(addr)
		if strings.Contains(point, ".") {
			host = point
		} else if point != "" {
			port = ":" + point
		}
	}
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = ":7071"
	}
	return fmt.Sprintf("http://%s%s/debug/callback", host, port)
}

func (s *accessTokenServer) initProviders() error {
	localCfg, configPath, err := s.loadLocalConfig("")
	if err != nil {
		return err
	}
	if localCfg.AccessTokenProviders == nil || len(localCfg.AccessTokenProviders.Providers) == 0 {
		return errors.New("accesstoken-server: access_token_providers 未配置")
	}
	var providers []providerMeta
	for _, group := range localCfg.AccessTokenProviders.Providers {
		if group == nil || len(group.Apps) == 0 {
			continue
		}
		meta := providerMeta{
			Code: strings.TrimSpace(group.Code),
			Name: firstNonEmpty(strings.TrimSpace(group.Name), strings.TrimSpace(group.Code)),
		}
		for _, appCfg := range group.Apps {
			if appCfg == nil || len(appCfg.AuthModes) == 0 {
				continue
			}
			appMeta := providerAppMeta{
				Code:         strings.TrimSpace(appCfg.Code),
				Name:         firstNonEmpty(strings.TrimSpace(appCfg.Name), strings.TrimSpace(appCfg.Code)),
				ProviderCode: appCfg.ProviderCodeValue(),
				AppKey:       appCfg.ProviderCodeValue(),
				ConfigPath:   configPath,
				APIVersion:   strings.TrimSpace(appCfg.ApiVersion),
			}
			for _, mode := range appCfg.AuthModes {
				if mode == nil {
					continue
				}
				modeKey := strings.TrimSpace(mode.Key)
				if modeKey == "" {
					modeKey = "default"
				}
				appMeta.Modes = append(appMeta.Modes, providerAppModeMeta{
					Key:          modeKey,
					Label:        firstNonEmpty(strings.TrimSpace(mode.Label), modeKey),
					ProviderCode: mode.EffectiveProviderCode(appCfg),
					OauthKey:     extractOauthKey(mode),
				})
			}
			if len(appMeta.Modes) > 0 {
				meta.Apps = append(meta.Apps, appMeta)
			}
		}
		if len(meta.Apps) > 0 {
			providers = append(providers, meta)
		}
	}
	if len(providers) == 0 {
		return errors.New("accesstoken-server: access_token_providers 中没有可用的 Provider/App")
	}
	s.providers = providers
	if group, app, mode := localCfg.AccessTokenProviders.FirstSelection(); app != nil && mode != nil {
		s.defaultProviderUI = strings.TrimSpace(group.Code)
		s.defaultApp = strings.TrimSpace(app.Code)
		s.defaultMode = strings.TrimSpace(mode.Key)
		s.defaultProviderCode = mode.EffectiveProviderCode(app)
	}
	buf, err := json.Marshal(s.providers)
	if err != nil {
		return err
	}
	s.providersJSON = template.JS(buf)
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func extractOauthKey(mode *config.AccessTokenAuthMode) string {
	if mode == nil {
		return ""
	}
	switch mode.ConfigKind() {
	case providerGoogleYouTube:
		if mode.GoogleYouTubeConfig != nil {
			return strings.TrimSpace(mode.GoogleYouTubeConfig.OauthKey)
		}
	case providerGoogleBlogger:
		if mode.GoogleBloggerConfig != nil {
			return strings.TrimSpace(mode.GoogleBloggerConfig.OauthKey)
		}
	}
	return ""
}

func (s *accessTokenServer) requireAPIToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.apiToken == "" {
			next.ServeHTTP(w, r)
			return
		}
		token := s.extractAPIToken(r)
		if token == "" || token != s.apiToken {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *accessTokenServer) extractAPIToken(r *http.Request) string {
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = strings.TrimSpace(token[7:])
	} else if header := strings.TrimSpace(r.Header.Get("X-API-Token")); header != "" {
		token = header
	} else if query := strings.TrimSpace(r.URL.Query().Get("api_token")); query != "" {
		token = query
	}
	return token
}

func (s *accessTokenServer) loadLocalConfig(override string) (*config.LocalConfig, string, error) {
	path := strings.TrimSpace(override)
	if path == "" {
		path = s.defaultConfigPath
	}
	localCfg := &config.LocalConfig{}
	if err := utils.LoadYAML(path, localCfg); err != nil {
		return nil, path, fmt.Errorf("load config %s: %w", path, err)
	}
	return localCfg, path, nil
}

func (s *accessTokenServer) resolveProviderConfig(providerCode, providerApp, modeKey, overridePath string) (*providerContext, error) {
	localCfg, path, err := s.loadLocalConfig(overridePath)
	if err != nil {
		return nil, err
	}
	if localCfg.AccessTokenProviders == nil {
		return nil, fmt.Errorf("missing access_token_providers in %s", path)
	}
	code := strings.TrimSpace(providerCode)
	appCode := strings.TrimSpace(providerApp)
	modeKey = strings.TrimSpace(modeKey)
	var provider *config.AccessTokenProvider
	var app *config.AccessTokenProviderApp
	if code != "" {
		provider, app = localCfg.AccessTokenProviders.FindAppByProviderCode(code)
	}
	if app == nil && appCode != "" {
		provider, app = localCfg.AccessTokenProviders.FindApp("", appCode)
	}
	if app == nil {
		defaultProvider, defaultApp, defaultMode := localCfg.AccessTokenProviders.FirstSelection()
		provider = defaultProvider
		app = defaultApp
		if modeKey == "" && defaultMode != nil {
			modeKey = defaultMode.Key
		}
	}
	if provider == nil || app == nil {
		return nil, fmt.Errorf("未找到 Provider/App (provider=%s app=%s) in %s", providerCode, providerApp, path)
	}
	mode := app.FindMode(modeKey)
	if mode == nil {
		return nil, fmt.Errorf("provider %s app %s 缺少授权模式 %s", provider.Code, app.Code, modeKey)
	}
	ctx := &providerContext{
		GroupCode:    strings.TrimSpace(provider.Code),
		GroupName:    firstNonEmpty(strings.TrimSpace(provider.Name), strings.TrimSpace(provider.Code)),
		AppCode:      strings.TrimSpace(app.Code),
		AppName:      firstNonEmpty(strings.TrimSpace(app.Name), strings.TrimSpace(app.Code)),
		ModeKey:      strings.TrimSpace(mode.Key),
		ProviderCode: mode.EffectiveProviderCode(app),
		ConfigPath:   path,
		APIVersion:   strings.TrimSpace(app.ApiVersion),
		Mode:         mode,
		App:          app,
		Provider:     provider,
	}
	switch mode.ConfigKind() {
	case providerGoogleYouTube:
		if mode.GoogleYouTubeConfig == nil {
			return nil, fmt.Errorf("provider %s app %s 缺少 google_youtube_config", provider.Code, app.Code)
		}
		ctx.Youtube = mode.GoogleYouTubeConfig
	case providerGoogleBlogger:
		if mode.GoogleBloggerConfig == nil {
			return nil, fmt.Errorf("provider %s app %s 缺少 google_blogger_config", provider.Code, app.Code)
		}
		ctx.Blogger = mode.GoogleBloggerConfig
	case providerByteDanceDouYin:
		if mode.ByteDanceDouYinConfig == nil {
			return nil, fmt.Errorf("provider %s app %s 缺少 byte_dance_douyin_config", provider.Code, app.Code)
		}
		ctx.Douyin = mode.ByteDanceDouYinConfig
	case providerRedBookJuGuang:
		if mode.RedBookJuGuangConfig == nil {
			return nil, fmt.Errorf("provider %s app %s 缺少 redbook_juguang_config", provider.Code, app.Code)
		}
		ctx.RedBook = mode.RedBookJuGuangConfig
	case providerBilibili:
		if mode.BiliConfig() == nil {
			return nil, fmt.Errorf("provider %s app %s 缺少 bilibili_config", provider.Code, app.Code)
		}
		ctx.Bilibili = mode.BiliConfig()
	default:
		return nil, fmt.Errorf("unsupported provider %s", mode.ConfigKind())
	}
	return ctx, nil
}

func (ctx *providerContext) clientConfig() *config.ClientConfig {
	if ctx == nil || ctx.Mode == nil {
		return nil
	}
	return ctx.Mode.ClientConfig()
}

func (ctx *providerContext) displayName() string {
	if ctx == nil {
		return "-"
	}
	name := strings.TrimSpace(ctx.GroupName)
	app := strings.TrimSpace(ctx.AppName)
	if name != "" && app != "" {
		return name + " · " + app
	}
	if app != "" {
		return app
	}
	if name != "" {
		return name
	}
	return "-"
}

func (s *accessTokenServer) resolveAccessTokenValue(ctx *providerContext, explicit string) (string, string, string, *oauthTokenRecord) {
	if ctx == nil {
		return "", "", "", nil
	}
	if token := strings.TrimSpace(explicit); token != "" {
		return token, "request", "payload", nil
	}
	if keys := app.EnvKeysForProvider(ctx.ProviderCode); len(keys) > 0 {
		for _, key := range keys {
			if val := strings.TrimSpace(os.Getenv(key)); val != "" {
				return val, "env", key, nil
			}
		}
	}
	if cfg := ctx.clientConfig(); cfg != nil && cfg.OAuthConfig != nil {
		if token := strings.TrimSpace(cfg.OAuthConfig.AccessToken); token != "" {
			return token, "config", "config", nil
		}
	}
	if rec := s.latestOAuthToken(ctx.ProviderCode, ctx.AppCode, ctx.ModeKey); rec != nil {
		if token := strings.TrimSpace(rec.AccessToken); token != "" {
			detail := rec.FlowID
			if detail == "" {
				detail = "oauth_cache"
			}
			return token, "flow", detail, rec
		}
	}
	return "", "", "", nil
}

func (s *accessTokenServer) buildOAuthAuthorizeURL(ctx *providerContext) (string, string, error) {
	if ctx == nil || ctx.Mode == nil {
		return "", "", errors.New("provider 配置缺失")
	}
	clientCfg := ctx.clientConfig()
	if clientCfg == nil || clientCfg.OAuthConfig == nil {
		return "", "", errors.New("provider 缺少 OAuth 配置")
	}
	oauthCfg := clientCfg.OAuthConfig
	redirect := strings.TrimSpace(oauthCfg.RedirectUrl)
	if redirect == "" {
		redirect = s.defaultCallbackURL
	}
	if strings.TrimSpace(oauthCfg.ClientID) == "" {
		return "", "", errors.New("OAuth client_id 未配置")
	}
	scope := strings.TrimSpace(oauthCfg.Scope)
	if scope == "" {
		return "", "", errors.New("OAuth scope 未配置")
	}
	authEndpoint := strings.TrimSpace(oauthCfg.OAuthUrl)
	if authEndpoint == "" && ctx.Mode.ConfigKind() == providerGoogleYouTube {
		authEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	}
	if authEndpoint == "" {
		return "", "", errors.New("OAuth oauth_url 未配置")
	}
	state := s.registerOAuthState(ctx, ctx.ConfigPath)
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", oauthCfg.ClientID)
	query.Set("redirect_uri", redirect)
	query.Set("scope", scope)
	query.Set("state", state)
	switch ctx.Mode.ConfigKind() {
	case providerGoogleYouTube, providerGoogleBlogger:
		query.Set("access_type", "offline")
		query.Set("include_granted_scopes", "true")
		if strings.Contains(scope, "youtube") {
			query.Set("prompt", "consent")
		}
	}
	return fmt.Sprintf("%s?%s", authEndpoint, query.Encode()), state, nil
}

func (s *accessTokenServer) registerOAuthState(ctx *providerContext, configPath string) string {
	buf := make([]byte, 18)
	if _, err := rand.Read(buf); err != nil {
		timestamp := time.Now().UnixNano()
		state := fmt.Sprintf("accesstoken-%d", timestamp)
		s.oauthStateMu.Lock()
		s.oauthStates[state] = &oauthState{
			ProviderCode: ctx.ProviderCode,
			ProviderApp:  ctx.AppCode,
			ModeKey:      ctx.ModeKey,
			ConfigPath:   configPath,
			CreatedAt:    time.Now(),
		}
		s.oauthStateMu.Unlock()
		return state
	}
	state := base64.RawURLEncoding.EncodeToString(buf)
	s.oauthStateMu.Lock()
	s.oauthStates[state] = &oauthState{
		ProviderCode: ctx.ProviderCode,
		ProviderApp:  ctx.AppCode,
		ModeKey:      ctx.ModeKey,
		ConfigPath:   configPath,
		CreatedAt:    time.Now(),
	}
	s.oauthStateMu.Unlock()
	return state
}

func (s *accessTokenServer) popOAuthState(state string) *oauthState {
	state = strings.TrimSpace(state)
	if state == "" {
		return nil
	}
	s.oauthStateMu.Lock()
	defer s.oauthStateMu.Unlock()
	payload, ok := s.oauthStates[state]
	if ok {
		delete(s.oauthStates, state)
	}
	return payload
}

func (s *accessTokenServer) completeOAuthFlow(ctx context.Context, state, code string, callback *callbackRecord) (*oauthTokenRecord, error) {
	payload := s.popOAuthState(state)
	if payload == nil {
		return nil, errors.New("state 无效或已过期，请重新发起授权")
	}
	providerCtx, err := s.resolveProviderConfig(payload.ProviderCode, payload.ProviderApp, payload.ModeKey, payload.ConfigPath)
	if err != nil {
		return nil, err
	}
	exchange, err := s.exchangeAuthorizationCode(ctx, providerCtx, code)
	if err != nil {
		return nil, err
	}
	record := &oauthTokenRecord{
		ProviderCode: providerCtx.ProviderCode,
		ProviderApp:  providerCtx.AppCode,
		AuthMode:     providerCtx.ModeKey,
		ConfigPath:   providerCtx.ConfigPath,
		FlowID:       fmt.Sprintf("oauth-%s", state),
		AccessToken:  exchange.AccessToken,
		RefreshToken: exchange.RefreshToken,
		TokenType:    exchange.TokenType,
		Scope:        exchange.Scope,
		ExpiresIn:    exchange.ExpiresIn,
		StoredAt:     time.Now().UTC(),
		Source:       "authorization_code",
		Callback:     cloneCallbackRecord(callback),
	}
	if record.ExpiresIn <= 0 {
		record.ExpiresIn = app.DefaultAccessTokenTTLSeconds
	}
	record.TokenExpireAt = record.StoredAt.Add(time.Duration(record.ExpiresIn) * time.Second)
	record.FlowTTLSeconds = s.flowTTLSeconds
	if record.FlowTTLSeconds <= 0 {
		record.FlowTTLSeconds = defaultFlowTTLSeconds
	}
	record.FlowExpireAt = record.StoredAt.Add(time.Duration(record.FlowTTLSeconds) * time.Second)
	if err := s.saveOAuthTokenRecord(record); err != nil {
		s.logger.ErrorF("accesstoken-server: 保存授权 token 失败: %v", err)
	}
	s.logFlowAction("oauth.complete", providerCtx, record.FlowID, "authorization_code", "")
	return record, nil
}

type oauthExchangePayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
}

func (s *accessTokenServer) exchangeAuthorizationCode(ctx context.Context, pctx *providerContext, code string) (*oauthExchangePayload, error) {
	clientCfg := pctx.clientConfig()
	if clientCfg == nil || clientCfg.OAuthConfig == nil {
		return nil, errors.New("provider 缺少 OAuth 配置")
	}
	oauthCfg := clientCfg.OAuthConfig
	tokenURL := strings.TrimSpace(oauthCfg.AccessTokenUrl)
	if tokenURL == "" {
		return nil, errors.New("未配置 access_token_url")
	}
	redirectURI := strings.TrimSpace(oauthCfg.RedirectUrl)
	if redirectURI == "" {
		redirectURI = s.defaultCallbackURL
	}
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", redirectURI)
	values.Set("client_id", oauthCfg.ClientID)
	values.Set("client_secret", oauthCfg.ClientSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token endpoint 返回 %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	payload := &oauthExchangePayload{}
	if err := json.Unmarshal(body, payload); err != nil {
		return nil, fmt.Errorf("解析 token 响应失败: %w", err)
	}
	if strings.TrimSpace(payload.AccessToken) == "" {
		return nil, errors.New("token 响应缺少 access_token")
	}
	return payload, nil
}

func (s *accessTokenServer) cacheKey(providerCode, appCode, mode string) string {
	providerCode = strings.TrimSpace(providerCode)
	appCode = strings.TrimSpace(appCode)
	mode = strings.TrimSpace(mode)
	if providerCode == "" && appCode == "" {
		return ""
	}
	return fmt.Sprintf("accesstoken:oauth:%s:%s:%s", providerCode, appCode, mode)
}

func (s *accessTokenServer) cacheKeyForToken(rec *oauthTokenRecord) string {
	if rec == nil {
		return ""
	}
	return s.cacheKey(rec.ProviderCode, rec.ProviderApp, rec.AuthMode)
}

func (s *accessTokenServer) flowIndexKey(flowID string) string {
	flowID = strings.TrimSpace(flowID)
	if flowID == "" {
		return ""
	}
	return fmt.Sprintf("accesstoken:oauth:flow:%s", flowID)
}

func (s *accessTokenServer) saveOAuthTokenRecord(rec *oauthTokenRecord) error {
	if rec == nil {
		return nil
	}
	if rec.StorageBackend == "" {
		rec.StorageBackend = s.storageBackend
	}
	s.enrichFlowMetadata(rec)
	key := s.cacheKeyForToken(rec)
	s.oauthTokenMu.Lock()
	if key != "" {
		s.oauthTokens[key] = rec
	}
	s.oauthTokenMu.Unlock()
	if s.cache != nil {
		data, err := json.Marshal(rec)
		if err == nil {
			expiration := s.flowTTL
			if expiration <= 0 {
				expiration = time.Duration(rec.FlowTTLSeconds) * time.Second
			}
			if expiration <= 0 {
				expiration = time.Hour * 24
			}
			ctx := context.Background()
			if key != "" {
				_ = s.cache.Set(ctx, key, data, expiration)
			}
			if flowKey := s.flowIndexKey(rec.FlowID); flowKey != "" && key != "" {
				_ = s.cache.Set(ctx, flowKey, []byte(key), expiration)
			}
		}
	}
	return nil
}

func (s *accessTokenServer) listOAuthTokenRecords(providerCode, appCode, mode string) []*oauthTokenRecord {
	now := time.Now().UTC()
	mode = strings.TrimSpace(mode)
	appCode = strings.TrimSpace(appCode)
	providerCode = strings.TrimSpace(providerCode)
	s.oauthTokenMu.RLock()
	defer s.oauthTokenMu.RUnlock()
	var items []*oauthTokenRecord
	for _, rec := range s.oauthTokens {
		if rec == nil {
			continue
		}
		s.enrichFlowMetadata(rec)
		if rec.FlowExpireAt.Before(now) {
			continue
		}
		if providerCode != "" && !strings.EqualFold(rec.ProviderCode, providerCode) {
			continue
		}
		if appCode != "" && !strings.EqualFold(rec.ProviderApp, appCode) {
			continue
		}
		if mode != "" && !strings.EqualFold(rec.AuthMode, mode) {
			continue
		}
		items = append(items, rec)
	}
	if len(items) == 0 {
		if cached := s.fetchTokenFromCache(providerCode, appCode, mode); cached != nil {
			items = append(items, cached)
		}
	}
	return items
}

func (s *accessTokenServer) latestOAuthToken(providerCode, appCode, mode string) *oauthTokenRecord {
	records := s.listOAuthTokenRecords(providerCode, appCode, mode)
	if len(records) == 0 {
		return nil
	}
	var latest *oauthTokenRecord
	for _, rec := range records {
		if rec == nil {
			continue
		}
		if latest == nil || rec.StoredAt.After(latest.StoredAt) {
			latest = rec
		}
	}
	return latest
}

func (s *accessTokenServer) fetchTokenFromCache(providerCode, appCode, mode string) *oauthTokenRecord {
	if s.cache == nil {
		return nil
	}
	key := s.cacheKey(providerCode, appCode, mode)
	if key == "" {
		return nil
	}
	raw, err := s.cache.Get(context.Background(), key)
	if err != nil || len(raw) == 0 {
		return nil
	}
	rec := &oauthTokenRecord{}
	if err := json.Unmarshal(raw, rec); err != nil {
		return nil
	}
	s.enrichFlowMetadata(rec)
	if rec.FlowExpireAt.Before(time.Now().UTC()) {
		return nil
	}
	s.oauthTokenMu.Lock()
	s.oauthTokens[key] = rec
	s.oauthTokenMu.Unlock()
	return rec
}

func (s *accessTokenServer) fetchTokenByFlowID(flowID string) *oauthTokenRecord {
	rec, err := s.lookupFlowByID(flowID)
	if err != nil {
		return nil
	}
	return rec
}

func (s *accessTokenServer) lookupFlowByID(flowID string) (*oauthTokenRecord, error) {
	flowID = strings.TrimSpace(flowID)
	if flowID == "" {
		return nil, errFlowNotFound
	}
	now := time.Now().UTC()
	s.oauthTokenMu.RLock()
	for _, rec := range s.oauthTokens {
		if rec == nil {
			continue
		}
		s.enrichFlowMetadata(rec)
		if strings.EqualFold(rec.FlowID, flowID) {
			if rec.FlowExpireAt.Before(now) {
				s.oauthTokenMu.RUnlock()
				return nil, errFlowExpired
			}
			s.oauthTokenMu.RUnlock()
			return rec, nil
		}
	}
	s.oauthTokenMu.RUnlock()
	indexKey := s.flowIndexKey(flowID)
	if indexKey == "" {
		return nil, errFlowNotFound
	}
	keyBytes, err := s.cachedGet(context.Background(), indexKey)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) == 0 {
		return nil, errFlowNotFound
	}
	raw, err := s.cachedGet(context.Background(), string(keyBytes))
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, errFlowNotFound
	}
	rec := &oauthTokenRecord{}
	if err := json.Unmarshal(raw, rec); err != nil {
		return nil, err
	}
	s.enrichFlowMetadata(rec)
	if rec.FlowExpireAt.Before(now) {
		return nil, errFlowExpired
	}
	s.oauthTokenMu.Lock()
	s.oauthTokens[string(keyBytes)] = rec
	s.oauthTokenMu.Unlock()
	return rec, nil
}

func (s *accessTokenServer) cachedGet(ctx context.Context, key string) ([]byte, error) {
	if key == "" {
		return nil, nil
	}
	if s.cache != nil {
		data, err := s.cache.Get(ctx, key)
		if err == nil && len(data) > 0 {
			return data, nil
		}
		if err != nil && err != redis.Nil {
			return nil, err
		}
	}
	if s.redis != nil {
		data, err := s.redis.Get(ctx, key).Bytes()
		if err == redis.Nil {
			return nil, nil
		}
		return data, err
	}
	return nil, nil
}

func (s *accessTokenServer) scanFlowRecords(ctx context.Context, provider, appCode, mode string, limit int, cursor string) ([]*oauthTokenRecord, string, error) {
	if s.redis == nil {
		return nil, "", nil
	}
	var start uint64
	if cursor != "" {
		if parsed, err := strconv.ParseUint(cursor, 10, 64); err == nil {
			start = parsed
		}
	}
	records := make([]*oauthTokenRecord, 0, limit)
	seen := make(map[string]struct{}, limit*2)
	nextCursor := start
	for len(records) < limit {
		keys, newCursor, err := s.redis.Scan(ctx, nextCursor, flowIndexPrefix+"*", int64(limit*4)).Result()
		if err != nil {
			return nil, "", err
		}
		nextCursor = newCursor
		for _, key := range keys {
			if !strings.HasPrefix(key, flowIndexPrefix) {
				continue
			}
			flowID := strings.TrimPrefix(key, flowIndexPrefix)
			if flowID == "" {
				continue
			}
			if _, ok := seen[flowID]; ok {
				continue
			}
			seen[flowID] = struct{}{}
			rec, err := s.lookupFlowByID(flowID)
			if err != nil {
				if errors.Is(err, errFlowExpired) || errors.Is(err, errFlowNotFound) {
					continue
				}
				continue
			}
			if !s.matchesProvider(rec, provider, appCode, mode) {
				continue
			}
			records = append(records, rec)
			if len(records) >= limit {
				break
			}
		}
		if nextCursor == 0 {
			break
		}
	}
	var next string
	if nextCursor != 0 {
		next = strconv.FormatUint(nextCursor, 10)
	}
	return records, next, nil
}

func (s *accessTokenServer) matchesProvider(rec *oauthTokenRecord, provider, appCode, mode string) bool {
	if rec == nil {
		return false
	}
	if provider != "" && !strings.EqualFold(rec.ProviderCode, provider) {
		return false
	}
	if appCode != "" && !strings.EqualFold(rec.ProviderApp, appCode) {
		return false
	}
	if mode != "" && !strings.EqualFold(rec.AuthMode, mode) {
		return false
	}
	return true
}

func (s *accessTokenServer) buildFlowSummary(rec *oauthTokenRecord) map[string]any {
	if rec == nil {
		return map[string]any{}
	}
	s.enrichFlowMetadata(rec)
	summary := map[string]any{
		"flow_id":            rec.FlowID,
		"provider_code":      rec.ProviderCode,
		"provider_app":       rec.ProviderApp,
		"provider_auth_mode": rec.AuthMode,
		"storage_backend":    rec.StorageBackend,
		"expire_at":          rec.FlowExpireAt,
		"flow_ttl_seconds":   rec.FlowTTLSeconds,
		"status":             s.flowStatus(rec),
		"masked_token":       app.MaskToken(rec.AccessToken),
		"masked_account":     s.maskAccount(rec),
		"source":             rec.Source,
		"config_path":        rec.ConfigPath,
		"token_type":         rec.TokenType,
		"expires_in":         rec.ExpiresIn,
	}
	if rec.Callback != nil && rec.Callback.State != "" {
		summary["flow_state"] = rec.Callback.State
	}
	return summary
}

func (s *accessTokenServer) flowStatus(rec *oauthTokenRecord) string {
	if rec == nil {
		return "invalid"
	}
	now := time.Now().UTC()
	if rec.FlowExpireAt.Before(now) {
		return "expired"
	}
	if strings.TrimSpace(rec.AccessToken) == "" {
		return "pending"
	}
	return "authorized"
}

func (s *accessTokenServer) buildFlowPayload(rec *oauthTokenRecord) map[string]any {
	if rec == nil {
		return nil
	}
	s.enrichFlowMetadata(rec)
	payload := map[string]any{
		"flow_id":            rec.FlowID,
		"provider_code":      rec.ProviderCode,
		"provider_app":       rec.ProviderApp,
		"provider_auth_mode": rec.AuthMode,
		"config_path":        rec.ConfigPath,
		"access_token":       rec.AccessToken,
		"masked_token":       app.MaskToken(rec.AccessToken),
		"refresh_token":      rec.RefreshToken,
		"token_type":         rec.TokenType,
		"scope":              rec.Scope,
		"expires_in":         rec.ExpiresIn,
		"expire_at":          rec.ExpireAt,
		"flow_expire_at":     rec.FlowExpireAt,
		"flow_ttl_seconds":   rec.FlowTTLSeconds,
		"storage_backend":    rec.StorageBackend,
		"token_source":       rec.Source,
	}
	if rec.Callback != nil {
		payload["callback"] = rec.Callback
	}
	return payload
}

func (s *accessTokenServer) maskAccount(rec *oauthTokenRecord) string {
	if rec == nil {
		return "-"
	}
	parts := []string{safeValue(rec.ProviderCode), safeValue(rec.ProviderApp), safeValue(rec.AuthMode)}
	return strings.Join(parts, "/")
}

func (s *accessTokenServer) enrichFlowMetadata(rec *oauthTokenRecord) {
	if rec == nil {
		return
	}
	if rec.FlowTTLSeconds <= 0 {
		rec.FlowTTLSeconds = s.flowTTLSeconds
		if rec.FlowTTLSeconds <= 0 {
			rec.FlowTTLSeconds = defaultFlowTTLSeconds
		}
	}
	if rec.FlowExpireAt.IsZero() {
		if !rec.StoredAt.IsZero() {
			rec.FlowExpireAt = rec.StoredAt.Add(time.Duration(rec.FlowTTLSeconds) * time.Second)
		} else if !rec.ExpireAt.IsZero() {
			rec.FlowExpireAt = rec.ExpireAt
		}
	}
	if rec.TokenExpireAt.IsZero() && rec.ExpiresIn > 0 {
		rec.TokenExpireAt = rec.StoredAt.Add(time.Duration(rec.ExpiresIn) * time.Second)
	}
	if rec.ExpireAt.IsZero() && !rec.FlowExpireAt.IsZero() {
		rec.ExpireAt = rec.FlowExpireAt
	}
	if rec.StorageBackend == "" {
		rec.StorageBackend = s.storageBackend
	}
}

func (s *accessTokenServer) validateListenAddr() error {
	host := listenHostFromAddr(s.listenAddr)
	lower := strings.ToLower(host)
	s.listenAddrPublic = host == "" || lower == "0.0.0.0" || lower == "::"
	if s.listenAddrPublic && s.logger != nil {
		s.logger.WarnF("accesstoken-server: listen_addr=%s is publicly reachable; override ACCESSTOKEN_LISTEN_ADDR to restrict to 127.0.0.1 if needed", s.listenAddr)
	}
	return nil
}

func listenHostFromAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return ""
	}
	if strings.HasPrefix(addr, "[") {
		if idx := strings.Index(addr, "]"); idx != -1 {
			return strings.TrimSpace(addr[1:idx])
		}
	}
	if strings.HasPrefix(addr, ":") {
		return ""
	}
	if strings.Contains(addr, ":") {
		host, _, _ := strings.Cut(addr, ":")
		return strings.TrimSpace(host)
	}
	return addr
}

func loadAccessTokenRedisConfig(configPath string) *config.AccessTokenRedisConfig {
	path := strings.TrimSpace(configPath)
	if path == "" {
		return nil
	}
	localCfg := &config.LocalConfig{}
	if err := utils.LoadYAML(path, localCfg); err != nil {
		return nil
	}
	if localCfg.AccessTokenProviders == nil {
		return nil
	}
	return localCfg.AccessTokenProviders.Redis
}
