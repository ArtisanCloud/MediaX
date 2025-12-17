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
	ProviderCode string    `json:"provider_code"`
	ProviderApp  string    `json:"provider_app"`
	AuthMode     string    `json:"provider_auth_mode"`
	FlowID       string    `json:"flow_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	ExpiresIn    int       `json:"expires_in"`
	StoredAt     time.Time `json:"stored_at"`
	ExpireAt     time.Time `json:"expire_at"`
	Source       string    `json:"source"`
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
}

func newAccessTokenServer(defaultConfigPath, listenAddr string) (*accessTokenServer, cacheCloser, error) {
	logCfg := app.BuildFileLogConfig(resolveLogLevel(), "logs/accesstoken-server-info.log", "logs/accesstoken-server-error.log")
	cacheStore, redisClient, closer, err := buildCacheStore()
	if err != nil {
		return nil, nil, err
	}
	mediaX := client.NewMediaX(&config.MediaXConfig{Logger: logCfg}, cacheStore)
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
	}
	server.defaultCallbackURL = server.buildDefaultCallbackURL()
	if err := server.initProviders(); err != nil {
		return nil, nil, err
	}
	return server, closer, nil
}

func (s *accessTokenServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/accesstoken", s.handleDebugPage)
	mux.HandleFunc("/debug", s.handleDebugPage)
	mux.Handle("/debug/callback", http.HandlerFunc(s.handleDebugCallback))
	mux.Handle("/api/callbacks", s.requireAPIToken(http.HandlerFunc(s.handleListCallbacks)))
	mux.Handle("/api/callbacks/clear", s.requireAPIToken(http.HandlerFunc(s.handleClearCallbacks)))
	mux.Handle("/api/oauth/flow-indexes", s.requireAPIToken(http.HandlerFunc(s.handleListFlowIndexes)))
	mux.Handle("/api/oauth/start", s.requireAPIToken(http.HandlerFunc(s.handleOAuthStart)))
	mux.Handle("/api/oauth/tokens", s.requireAPIToken(http.HandlerFunc(s.handleListOAuthTokens)))
	mux.Handle("/accesstoken/token", s.requireAPIToken(http.HandlerFunc(s.handleToken)))
	mux.Handle("/accesstoken/call", s.requireAPIToken(http.HandlerFunc(s.handleCall)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func buildCacheStore() (cache.ICache, redis.Cmdable, cacheCloser, error) {
	addr := strings.TrimSpace(os.Getenv("ACCESSTOKEN_REDIS_ADDR"))
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	db := 0
	if v := strings.TrimSpace(os.Getenv("ACCESSTOKEN_REDIS_DB")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			db = n
		}
	}
	opts := &redis.Options{
		Addr:     addr,
		DB:       db,
		Username: strings.TrimSpace(os.Getenv("ACCESSTOKEN_REDIS_USERNAME")),
		Password: strings.TrimSpace(os.Getenv("ACCESSTOKEN_REDIS_PASSWORD")),
	}
	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, nil, func() {}, fmt.Errorf("ping redis %s: %w", addr, err)
	}
	return cache.NewRedisCache(client), client, func() { _ = client.Close() }, nil
}

func resolveListenAddr(flagPort string) string {
	if trimmed := strings.TrimSpace(flagPort); trimmed != "" {
		if strings.HasPrefix(trimmed, ":") {
			return trimmed
		}
		return ":" + trimmed
	}
	if env := strings.TrimSpace(os.Getenv("ACCESSTOKEN_LISTEN_ADDR")); env != "" {
		return env
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
		token := strings.TrimSpace(r.Header.Get("Authorization"))
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = strings.TrimSpace(token[7:])
		} else if header := strings.TrimSpace(r.Header.Get("X-API-Token")); header != "" {
			token = header
		} else if query := strings.TrimSpace(r.URL.Query().Get("api_token")); query != "" {
			token = query
		}
		if token == "" || token != s.apiToken {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
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
		if mode.BiliBiliConfig == nil {
			return nil, fmt.Errorf("provider %s app %s 缺少 bilibili_config", provider.Code, app.Code)
		}
		ctx.Bilibili = mode.BiliBiliConfig
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

func (s *accessTokenServer) resolveAccessTokenValue(ctx *providerContext, explicit string) (string, string, *oauthTokenRecord) {
	if ctx == nil {
		return "", "", nil
	}
	if token := strings.TrimSpace(explicit); token != "" {
		return token, "payload", nil
	}
	if keys := app.EnvKeysForProvider(ctx.ProviderCode); len(keys) > 0 {
		for _, key := range keys {
			if val := strings.TrimSpace(os.Getenv(key)); val != "" {
				return val, "env:" + key, nil
			}
		}
	}
	if cfg := ctx.clientConfig(); cfg != nil && cfg.OAuthConfig != nil {
		if token := strings.TrimSpace(cfg.OAuthConfig.AccessToken); token != "" {
			return token, "config", nil
		}
	}
	if rec := s.latestOAuthToken(ctx.ProviderCode, ctx.AppCode, ctx.ModeKey); rec != nil {
		if token := strings.TrimSpace(rec.AccessToken); token != "" {
			source := strings.TrimSpace(rec.Source)
			if rec.FlowID != "" {
				source = fmt.Sprintf("flow:%s", rec.FlowID)
			}
			if source == "" {
				source = "oauth_cache"
			}
			return token, source, rec
		}
	}
	return "", "", nil
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
	if authEndpoint == "" {
		authEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	}
	state := s.registerOAuthState(ctx, ctx.ConfigPath)
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", oauthCfg.ClientID)
	query.Set("redirect_uri", redirect)
	query.Set("scope", scope)
	query.Set("state", state)
	query.Set("access_type", "offline")
	query.Set("include_granted_scopes", "true")
	if strings.Contains(scope, "youtube") {
		query.Set("prompt", "consent")
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

func (s *accessTokenServer) completeOAuthFlow(ctx context.Context, state, code string) (*oauthTokenRecord, error) {
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
		FlowID:       fmt.Sprintf("oauth-%s", state),
		AccessToken:  exchange.AccessToken,
		RefreshToken: exchange.RefreshToken,
		TokenType:    exchange.TokenType,
		Scope:        exchange.Scope,
		ExpiresIn:    exchange.ExpiresIn,
		StoredAt:     time.Now().UTC(),
		Source:       "authorization_code",
	}
	if record.ExpiresIn <= 0 {
		record.ExpiresIn = app.DefaultAccessTokenTTLSeconds
	}
	record.ExpireAt = record.StoredAt.Add(time.Duration(record.ExpiresIn) * time.Second)
	if err := s.saveOAuthTokenRecord(record); err != nil {
		s.logger.ErrorF("accesstoken-server: 保存授权 token 失败: %v", err)
	}
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
	key := s.cacheKeyForToken(rec)
	s.oauthTokenMu.Lock()
	if key != "" {
		s.oauthTokens[key] = rec
	}
	s.oauthTokenMu.Unlock()
	if s.cache != nil {
		data, err := json.Marshal(rec)
		if err == nil {
			expiration := time.Duration(rec.ExpiresIn) * time.Second
			if expiration <= 0 {
				expiration = time.Hour
			}
			ctx := context.Background()
			if key != "" {
				_ = s.cache.Set(ctx, key, data, expiration)
			}
			if flowKey := s.flowIndexKey(rec.FlowID); flowKey != "" && key != "" {
				_ = s.cache.Set(ctx, flowKey, key, expiration)
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
		if rec.ExpireAt.Before(now) {
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
	if rec.ExpireAt.Before(time.Now().UTC()) {
		return nil
	}
	s.oauthTokenMu.Lock()
	s.oauthTokens[key] = rec
	s.oauthTokenMu.Unlock()
	return rec
}

func (s *accessTokenServer) fetchTokenByFlowID(flowID string) *oauthTokenRecord {
	flowID = strings.TrimSpace(flowID)
	if flowID == "" {
		return nil
	}
	now := time.Now().UTC()
	s.oauthTokenMu.RLock()
	for _, rec := range s.oauthTokens {
		if rec == nil {
			continue
		}
		if strings.EqualFold(rec.FlowID, flowID) && rec.ExpireAt.After(now) {
			s.oauthTokenMu.RUnlock()
			return rec
		}
	}
	s.oauthTokenMu.RUnlock()
	indexKey := s.flowIndexKey(flowID)
	if indexKey == "" {
		return nil
	}
	keyBytes, err := s.cachedGet(context.Background(), indexKey)
	if err != nil || len(keyBytes) == 0 {
		return nil
	}
	raw, err := s.cachedGet(context.Background(), string(keyBytes))
	if err != nil || len(raw) == 0 {
		return nil
	}
	rec := &oauthTokenRecord{}
	if err := json.Unmarshal(raw, rec); err != nil {
		return nil
	}
	if rec.ExpireAt.Before(now) {
		return nil
	}
	s.oauthTokenMu.Lock()
	s.oauthTokens[string(keyBytes)] = rec
	s.oauthTokenMu.Unlock()
	return rec
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
