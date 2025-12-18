package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type clientTokenServerOptions struct {
	localConfig *config.LocalConfig
	configPath  string
	listenAddr  string
	redis       redis.Cmdable
	mediaX      *client.MediaX
	cache       cache.ICache
}

type clientTokenServer struct {
	logger      *logger.Logger
	mediaX      *client.MediaX
	cache       cache.ICache
	redis       redis.Cmdable
	localConfig *config.LocalConfig
	configPath  string
	listenAddr  string
	apiToken    string

	tokenStore    *tokenCacheStore
	callbackStore *callbackLogStore
}

func newClientTokenServer(opts *clientTokenServerOptions) (*clientTokenServer, error) {
	if opts == nil || opts.localConfig == nil {
		return nil, fmt.Errorf("clienttoken: options/localConfig is nil")
	}
	if opts.mediaX == nil {
		return nil, fmt.Errorf("clienttoken: mediaX is nil")
	}
	apiToken := resolveAPIToken(opts.localConfig)
	server := &clientTokenServer{
		logger:        opts.mediaX.Logger,
		mediaX:        opts.mediaX,
		cache:         opts.cache,
		redis:         opts.redis,
		localConfig:   opts.localConfig,
		configPath:    opts.configPath,
		listenAddr:    opts.listenAddr,
		apiToken:      apiToken,
		tokenStore:    newTokenCacheStore(opts.redis),
		callbackStore: newCallbackLogStore(50),
	}
	return server, nil
}

func resolveAPIToken(localCfg *config.LocalConfig) string {
	if val := strings.TrimSpace(osGetenv("CLIENTTOKEN_API_TOKEN")); val != "" {
		return val
	}
	if localCfg != nil && localCfg.ClientTokenProviders != nil {
		if _, _, mode := localCfg.ClientTokenProviders.FirstSelection(); mode != nil && mode.WechatOfficialAccountConfig != nil {
			if token := strings.TrimSpace(mode.WechatOfficialAccountConfig.APIToken); token != "" {
				return token
			}
		}
	}
	return "dev-clienttoken"
}

// osGetenv extracted for testing.
var osGetenv = func(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func (s *clientTokenServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug", s.handleDebugPage)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/client-token/token", s.requireAPIToken(http.HandlerFunc(s.handleToken)))
	mux.Handle("/client-token/cache", s.requireAPIToken(http.HandlerFunc(s.handleCache)))
	mux.Handle("/client-token/call", s.requireAPIToken(http.HandlerFunc(s.handleCall)))
	mux.Handle("/client-token/message/validate", s.requireAPIToken(http.HandlerFunc(s.handleMessageValidate)))
	mux.Handle("/client-token/message/callback", s.requireAPIToken(http.HandlerFunc(s.handleMessageCallback)))
	return mux
}

func (s *clientTokenServer) requireAPIToken(next http.Handler) http.Handler {
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

type providerContext struct {
	ProviderCode         string
	AppCode              string
	ModeKey              string
	ConfigPath           string
	Config               *config.ClientTokenProviderConfig
	CacheKey             string
	TTLSeconds           int
	RefreshBeforeSeconds int
}

func (ctx *providerContext) AppID() string {
	if ctx == nil || ctx.Config == nil {
		return ""
	}
	return ctx.Config.AppIDValue()
}

func (s *clientTokenServer) resolveProviderConfig(providerCode, appCode, modeKey, override string) (*providerContext, error) {
	localCfg := s.localConfig
	configPath := s.configPath
	if strings.TrimSpace(override) != "" && !strings.EqualFold(strings.TrimSpace(override), strings.TrimSpace(s.configPath)) {
		tmpCfg := &config.LocalConfig{}
		if err := utils.LoadYAML(strings.TrimSpace(override), tmpCfg); err != nil {
			return nil, fmt.Errorf("load config %s: %w", override, err)
		}
		if tmpCfg.ClientTokenProviders == nil {
			return nil, fmt.Errorf("config %s 缺少 client_token_providers", override)
		}
		localCfg = tmpCfg
		configPath = strings.TrimSpace(override)
	}
	if localCfg.ClientTokenProviders == nil {
		return nil, fmt.Errorf("client_token_providers 未配置")
	}
	code := strings.TrimSpace(providerCode)
	app := strings.TrimSpace(appCode)
	mode := strings.TrimSpace(modeKey)
	var provider *config.ClientTokenProvider
	var appCfg *config.ClientTokenProviderApp
	if code != "" {
		provider, appCfg = localCfg.ClientTokenProviders.FindAppByProviderCode(code)
	}
	if appCfg == nil && app != "" {
		provider, appCfg = localCfg.ClientTokenProviders.FindApp("", app)
	}
	if appCfg == nil {
		provider, appCfg = localCfg.ClientTokenProviders.FindApp("", "")
	}
	if provider == nil || appCfg == nil {
		return nil, fmt.Errorf("未找到 ClientToken Provider/App (provider=%s app=%s)", providerCode, appCode)
	}
	modeCfg := appCfg.FindMode(mode)
	if modeCfg == nil || modeCfg.WechatOfficialAccountConfig == nil {
		return nil, fmt.Errorf("Provider %s app %s 缺少 wechat_official_account_config", provider.Code, appCfg.Code)
	}
	cfg := modeCfg.WechatOfficialAccountConfig
	cacheKey := cfg.EffectiveRedisKey(cfg.AppIDValue())
	ctx := &providerContext{
		ProviderCode:         provider.Code,
		AppCode:              appCfg.Code,
		ModeKey:              modeCfg.Key,
		ConfigPath:           configPath,
		Config:               cfg,
		CacheKey:             cacheKey,
		TTLSeconds:           cfg.EffectiveTTLSeconds(),
		RefreshBeforeSeconds: cfg.EffectiveRefreshBefore(),
	}
	return ctx, nil
}

type tokenRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode"`
	ConfigPath       string `json:"config_path"`
	ForceRefresh     bool   `json:"force_refresh"`
}

func (s *clientTokenServer) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()
	var req tokenRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid json: %v", err)
		s.logMetric("token", "", "", "error", "", start, err)
		return
	}
	ctx, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		s.logMetric("token", req.ProviderCode, ctxAppID(ctx), "error", "", start, err)
		return
	}
	record, err := s.refreshToken(r.Context(), ctx)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "refresh token failed: %v", err)
		s.logMetric("token", ctx.ProviderCode, ctx.AppID(), "error", "", start, err)
		return
	}
	resp := map[string]any{
		"access_token": record.AccessToken,
		"masked_token": maskToken(record.AccessToken),
		"appid":        ctx.AppID(),
		"cache_key":    ctx.CacheKey,
		"stored_at":    record.StoredAt,
		"expire_at":    record.ExpireAt,
		"ttl_seconds":  maxInt(int(record.TTLSeconds(time.Now())), 0),
		"token_source": record.Source,
		"config_path":  ctx.ConfigPath,
	}
	s.writeJSON(w, http.StatusOK, resp)
	s.logMetric("token", ctx.ProviderCode, ctx.AppID(), "success", record.Source, start, nil)
}

type cacheRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode"`
	ConfigPath       string `json:"config_path"`
}

func (s *clientTokenServer) handleCache(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	req := cacheRequest{
		ProviderCode:     r.URL.Query().Get("provider_code"),
		ProviderApp:      r.URL.Query().Get("provider_app"),
		ProviderAuthMode: r.URL.Query().Get("provider_auth_mode"),
		ConfigPath:       r.URL.Query().Get("config_path"),
	}
	ctx, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		s.logMetric("cache", req.ProviderCode, ctxAppID(ctx), "error", "", start, err)
		return
	}
	record, source, err := s.tokenStore.fetch(r.Context(), ctx.CacheKey)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "load cache failed: %v", err)
		s.logMetric("cache", ctx.ProviderCode, ctx.AppID(), "error", "", start, err)
		return
	}
	if record == nil {
		s.writeJSON(w, http.StatusOK, map[string]any{
			"cached":       false,
			"token_source": source,
			"cache_key":    ctx.CacheKey,
		})
		s.logMetric("cache", ctx.ProviderCode, ctx.AppID(), "miss", "", start, nil)
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]any{
		"cached":       true,
		"token_source": source,
		"cache_key":    ctx.CacheKey,
		"appid":        record.AppID,
		"stored_at":    record.StoredAt,
		"expire_at":    record.ExpireAt,
		"ttl_seconds":  maxInt(int(record.TTLSeconds(time.Now())), 0),
		"masked_token": maskToken(record.AccessToken),
	})
	s.logMetric("cache", ctx.ProviderCode, ctx.AppID(), "hit", source, start, nil)
}

type callRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode"`
	ConfigPath       string `json:"config_path"`
	Action           string `json:"action"`
	Method           string `json:"method"`
	Query            string `json:"query"`
	Body             string `json:"body"`
}

func (s *clientTokenServer) handleCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()
	var req callRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid json: %v", err)
		s.logMetric("call", "", "", "error", "", start, err)
		return
	}
	req.Action = strings.TrimSpace(req.Action)
	if req.Action == "" {
		s.writeError(w, http.StatusBadRequest, "action 不能为空")
		s.logMetric("call", req.ProviderCode, "", "error", "", start, errors.New("missing action"))
		return
	}
	ctxMeta, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		s.logMetric("call", req.ProviderCode, ctxAppID(ctxMeta), "error", "", start, err)
		return
	}
	tokenRec, err := s.ensureToken(r.Context(), ctxMeta)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "ensure token failed: %v", err)
		s.logMetric("call", ctxMeta.ProviderCode, ctxMeta.AppID(), "error", "", start, err)
		return
	}
	result, err := s.executeAPICall(r.Context(), ctxMeta, &req)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, err.Error())
		s.logMetric("call", ctxMeta.ProviderCode, ctxMeta.AppID(), "error", tokenRec.Source, start, err)
		return
	}
	resp := map[string]any{
		"result":        result,
		"token_source":  tokenRec.Source,
		"masked_token":  maskToken(tokenRec.AccessToken),
		"appid":         ctxMeta.AppID(),
		"action":        req.Action,
		"method":        strings.ToUpper(strings.TrimSpace(req.Method)),
		"cache_key":     ctxMeta.CacheKey,
		"request_time":  start.UTC(),
		"config_path":   ctxMeta.ConfigPath,
		"token_expires": tokenRec.ExpireAt,
	}
	s.writeJSON(w, http.StatusOK, resp)
	s.logMetric("call", ctxMeta.ProviderCode, ctxMeta.AppID(), "success", tokenRec.Source, start, nil)
}

func (s *clientTokenServer) executeAPICall(ctx context.Context, meta *providerContext, req *callRequest) (map[string]any, error) {
	wechatClient, err := s.mediaX.CreateWechatClientTokenClient(meta.Config)
	// ensure token handler uses cache
	if err != nil {
		return nil, fmt.Errorf("create wechat client: %w", err)
	}
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = http.MethodGet
	}
	queryMap := parseQueryString(req.Query)
	var bodyData interface{}
	if strings.TrimSpace(req.Body) != "" {
		if err := json.Unmarshal([]byte(req.Body), &bodyData); err != nil {
			return nil, fmt.Errorf("parse body json: %w", err)
		}
	}
	var result map[string]any
	switch method {
	case http.MethodGet:
		result = map[string]any{}
		_, err = wechatClient.GetBaseClient().HttpGet(ctx, req.Action, queryMap, nil, nil, &result)
	case http.MethodPost:
		result = map[string]any{}
		_, err = wechatClient.GetBaseClient().HttpPost(ctx, req.Action, queryMap, bodyData, nil, &result)
	default:
		return nil, fmt.Errorf("暂不支持 method=%s", method)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

type messageValidateRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode"`
	ConfigPath       string `json:"config_path"`
	AppID            string `json:"appid"`
	Signature        string `json:"signature"`
	Timestamp        string `json:"timestamp"`
	Nonce            string `json:"nonce"`
	EchoStr          string `json:"echostr"`
}

func (s *clientTokenServer) handleMessageValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()
	var req messageValidateRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid json: %v", err)
		s.logMetric("message.validate", "", "", "error", "", start, err)
		return
	}
	ctxMeta, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		s.logMetric("message.validate", req.ProviderCode, ctxAppID(ctxMeta), "error", "", start, err)
		return
	}
	cred := ctxMeta.Config.Credentials()
	if cred == nil || cred.MessageTokenValue() == "" {
		s.writeError(w, http.StatusBadRequest, "配置缺少 message_token")
		s.logMetric("message.validate", ctxMeta.ProviderCode, ctxMeta.AppID(), "error", "", start, fmt.Errorf("missing message token"))
		return
	}
	expected := calcSignature(cred.MessageTokenValue(), req.Timestamp, req.Nonce)
	resp := map[string]any{
		"valid":         strings.EqualFold(expected, strings.TrimSpace(req.Signature)),
		"expected":      expected,
		"signature":     strings.TrimSpace(req.Signature),
		"echo":          req.EchoStr,
		"provider_code": ctxMeta.ProviderCode,
		"appid":         ctxMeta.AppID(),
		"config_path":   ctxMeta.ConfigPath,
		"timestamp":     req.Timestamp,
		"nonce":         req.Nonce,
		"message_token": cred.MessageTokenValue(),
	}
	if resp["valid"].(bool) {
		s.logMetric("message.validate", ctxMeta.ProviderCode, ctxMeta.AppID(), "success", "", start, nil)
	} else {
		s.logMetric("message.validate", ctxMeta.ProviderCode, ctxMeta.AppID(), "error", "", start, fmt.Errorf("signature mismatch"))
	}
	s.writeJSON(w, http.StatusOK, resp)
}

func calcSignature(token, timestamp, nonce string) string {
	values := []string{strings.TrimSpace(token), strings.TrimSpace(timestamp), strings.TrimSpace(nonce)}
	sort.Strings(values)
	h := sha1.New()
	_, _ = io.WriteString(h, strings.Join(values, ""))
	return hex.EncodeToString(h.Sum(nil))
}

type messageCallbackRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode"`
	ConfigPath       string `json:"config_path"`
}

func (s *clientTokenServer) handleMessageCallback(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	bodyBytes, _ := io.ReadAll(r.Body)
	record := callbackRecord{
		Timestamp: time.Now(),
		Method:    r.Method,
		Query:     r.URL.RawQuery,
		Headers:   sanitizeHeaders(r.Header),
		Body:      string(bodyBytes),
	}
	s.callbackStore.append(record)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
	s.logMetric("message.callback", "", "", "success", "", start, nil)
}

func sanitizeHeaders(headers http.Header) map[string]string {
	if headers == nil {
		return nil
	}
	copied := make(map[string]string, len(headers))
	for k, v := range headers {
		if len(v) > 0 {
			copied[k] = v[0]
		}
	}
	return copied
}

func (s *clientTokenServer) handleDebugPage(w http.ResponseWriter, _ *http.Request) {
	page := `
<!DOCTYPE html>
<html lang="zh-CN">
<head><meta charset="utf-8"><title>ClientToken Debug</title></head>
<body>
<h1>ClientToken Server</h1>
<p>服务已启动，访问 API 请使用开发工具或浏览器脚本。</p>
<p>待 Phase3 完成后将提供图形界面。</p>
</body>
</html>`
	tpl := template.Must(template.New("debug").Parse(page))
	_ = tpl.Execute(w, nil)
}

func (s *clientTokenServer) refreshToken(ctx context.Context, meta *providerContext) (*tokenCacheRecord, error) {
	wechatClient, err := s.mediaX.CreateWechatClientTokenClient(meta.Config)
	if err != nil {
		return nil, err
	}
	tokenRes, err := wechatClient.AccessTokenHandler.ClientTokenHandler.GetRefreshedToken()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	ttl := meta.TTLSeconds
	if tokenRes != nil && tokenRes.ExpiresIn > 0 {
		ttl = int(tokenRes.ExpiresIn)
	}
	if ttl <= 0 {
		ttl = 7000
	}
	record := &tokenCacheRecord{
		AppID:       meta.AppID(),
		AccessToken: tokenRes.AccessToken,
		StoredAt:    now,
		ExpireAt:    now.Add(time.Duration(ttl) * time.Second),
		Source:      "refresh",
	}
	if err := s.tokenStore.save(ctx, meta.CacheKey, record, time.Duration(ttl)*time.Second); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *clientTokenServer) ensureToken(ctx context.Context, meta *providerContext) (*tokenCacheRecord, error) {
	record, _, err := s.tokenStore.fetch(ctx, meta.CacheKey)
	if err != nil {
		s.logger.WarnF("clienttoken: fetch cache failed: %v", err)
	}
	now := time.Now()
	needRefresh := false
	if record == nil {
		needRefresh = true
	} else if ttl := record.ExpireAt.Sub(now); ttl <= 0 {
		needRefresh = true
	} else if ttl <= time.Duration(meta.RefreshBeforeSeconds)*time.Second {
		needRefresh = true
	}
	if needRefresh {
		return s.refreshToken(ctx, meta)
	}
	return record, nil
}

func (s *clientTokenServer) writeError(w http.ResponseWriter, status int, format string, args ...interface{}) {
	w.WriteHeader(status)
	resp := map[string]any{
		"error": fmt.Sprintf(format, args...),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *clientTokenServer) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func decodeJSON(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func parseQueryString(raw string) *object.StringMap {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	values, err := url.ParseQuery(raw)
	if err != nil {
		return nil
	}
	m := object.StringMap{}
	for k, v := range values {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return &m
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ctxAppID(ctx *providerContext) string {
	if ctx == nil {
		return ""
	}
	return ctx.AppID()
}

func (record *tokenCacheRecord) TTLSeconds(now time.Time) float64 {
	if record == nil {
		return 0
	}
	return record.ExpireAt.Sub(now).Seconds()
}

func (s *clientTokenServer) logMetric(action, provider, appid, status, tokenSource string, start time.Time, err error) {
	latency := time.Since(start)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	s.logger.InfoF("clienttoken_metric: action=%s provider=%s appid=%s status=%s latency_ms=%d token_source=%s error=%s",
		action,
		firstOrDash(provider),
		firstOrDash(appid),
		status,
		latency.Milliseconds(),
		firstOrDash(tokenSource),
		firstOrDash(errMsg),
	)
}

func firstOrDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}
