package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	handler "github.com/ArtisanCloud/MediaX/internal/accesstoken/handler"
	"github.com/ArtisanCloud/MediaX/pkg/client"
	douyinresponse "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	wechatresponse "github.com/ArtisanCloud/MediaX/pkg/client/wechat/core/response"
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
	mux.Handle("/client-token/message/callbacks", s.requireAPIToken(http.HandlerFunc(s.handleCallbacks)))
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

type providerType string

const (
	providerTypeWechat    providerType = "wechat"
	providerTypeByteDance providerType = "byte_dance_douyin"
)

type providerContext struct {
	ProviderType         providerType
	ProviderCode         string
	AppCode              string
	ModeKey              string
	ConfigPath           string
	WechatConfig         *config.ClientTokenProviderConfig
	ByteDanceConfig      *config.ByteDanceDouYinConfig
	CacheKey             string
	TTLSeconds           int
	RefreshBeforeSeconds int
	RiskControlReady     bool
	RiskControlActive    bool
}

func (ctx *providerContext) AppID() string {
	if ctx == nil {
		return ""
	}
	switch ctx.ProviderType {
	case providerTypeWechat:
		if ctx.WechatConfig == nil {
			return ""
		}
		return ctx.WechatConfig.AppIDValue()
	case providerTypeByteDance:
		if ctx.ByteDanceConfig == nil {
			return ""
		}
		cred := ctx.ByteDanceConfig.ClientTokenCredential()
		if cred == nil {
			return ""
		}
		return cred.ClientKeyValue()
	default:
		return ""
	}
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
	if modeCfg == nil {
		return nil, fmt.Errorf("Provider %s app %s 缺少授权模式", provider.Code, appCfg.Code)
	}
	ctx := &providerContext{ProviderCode: provider.Code, AppCode: appCfg.Code, ModeKey: modeCfg.Key, ConfigPath: configPath}
	switch {
	case modeCfg.WechatOfficialAccountConfig != nil:
		cfg := modeCfg.WechatOfficialAccountConfig
		cacheKey := cfg.EffectiveRedisKey(cfg.AppIDValue())
		ctx.ProviderType = providerTypeWechat
		ctx.WechatConfig = cfg
		ctx.CacheKey = cacheKey
		ctx.TTLSeconds = cfg.EffectiveTTLSeconds()
		ctx.RefreshBeforeSeconds = cfg.EffectiveRefreshBefore()
	case modeCfg.ByteDanceDouYinConfig != nil:
		douyinCfg := modeCfg.ByteDanceDouYinConfig
		douyinCfg.NormalizeClientTokenCredentials()
		if douyinCfg.RiskControlActivated() && !douyinCfg.RiskControlReady() {
			return nil, fmt.Errorf("Provider %s app %s 已设置 device_id/risk_info 其中一项，需同时配置 device_id 与 risk_info", provider.Code, appCfg.Code)
		}
		cacheKey := douyinCfg.EffectiveRedisKey("")
		ctx.ProviderType = providerTypeByteDance
		ctx.ByteDanceConfig = douyinCfg
		ctx.CacheKey = cacheKey
		ctx.TTLSeconds = douyinCfg.EffectiveTTLSeconds()
		ctx.RefreshBeforeSeconds = douyinCfg.EffectiveRefreshBefore()
		ctx.RiskControlActive = douyinCfg.RiskControlActivated()
		ctx.RiskControlReady = douyinCfg.RiskControlReady()
	default:
		return nil, fmt.Errorf("Provider %s app %s 缺少支持的配置类型", provider.Code, appCfg.Code)
	}
	if strings.TrimSpace(ctx.CacheKey) == "" {
		return nil, fmt.Errorf("Provider %s app %s 缺少 redis_key/client_key", provider.Code, appCfg.Code)
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
	subject := apiTokenSubjectFromRequest(r)
	var req tokenRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid json: %v", err)
		s.logMetric("token", "", "", "error", "", start, err)
		return
	}
	ctx, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "%s", err.Error())
		s.logMetric("token", req.ProviderCode, ctxAppID(ctx), "error", "", start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.token.refresh",
			Action:          "refresh",
			ProviderCode:    req.ProviderCode,
			AppCode:         req.ProviderApp,
			Mode:            req.ProviderAuthMode,
			RedisKey:        "",
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
			Provider:        "",
			StorageBackend:  s.storageBackendLabel(),
			TokenSource:     "",
			TTL:             0,
		})
		return
	}
	record, err := s.refreshToken(r.Context(), ctx)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "refresh token failed: %v", err)
		s.logMetric("token", ctx.ProviderCode, ctx.AppID(), "error", "", start, err)
		tokenSrc := ""
		if record != nil {
			tokenSrc = record.Source
		}
		s.logAudit(auditLogFields{
			Event:           "clienttoken.token.refresh",
			Action:          "refresh",
			Provider:        string(ctx.ProviderType),
			ProviderCode:    ctx.ProviderCode,
			AppCode:         ctx.AppCode,
			Mode:            ctx.ModeKey,
			RedisKey:        ctx.CacheKey,
			StorageBackend:  s.storageBackendLabel(),
			TokenSource:     firstOrDash(tokenSrc),
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
		})
		return
	}
	ttlSeconds := maxInt(int(record.TTLSeconds(time.Now())), 0)
	storageBackend := s.storageBackendLabel()
	resp := map[string]any{
		"access_token": record.AccessToken,
		"masked_token": maskToken(record.AccessToken),
		"appid":        ctx.AppID(),
		"cache_key":    ctx.CacheKey,
		"stored_at":    record.StoredAt,
		"expire_at":    record.ExpireAt,
		"ttl_seconds":  ttlSeconds,
		"token_source": record.Source,
		"config_path":  ctx.ConfigPath,
	}
	s.writeJSON(w, http.StatusOK, resp)
	s.logMetric("token", ctx.ProviderCode, ctx.AppID(), "success", record.Source, start, nil)
	s.logAudit(auditLogFields{
		Event:           "clienttoken.token.refresh",
		Action:          "refresh",
		Provider:        string(ctx.ProviderType),
		ProviderCode:    ctx.ProviderCode,
		AppCode:         ctx.AppCode,
		Mode:            ctx.ModeKey,
		RedisKey:        ctx.CacheKey,
		TTL:             ttlSeconds,
		StorageBackend:  storageBackend,
		TokenSource:     firstOrDash(record.Source),
		APITokenSubject: subject,
		Status:          "success",
	})
}

type cacheRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode"`
	ConfigPath       string `json:"config_path"`
}

func (s *clientTokenServer) handleCache(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleCacheGet(w, r)
	case http.MethodDelete:
		s.handleCacheDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *clientTokenServer) handleCacheGet(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	subject := apiTokenSubjectFromRequest(r)
	req := s.parseCacheRequest(r)
	ctx, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "%s", err.Error())
		s.logMetric("cache", req.ProviderCode, ctxAppID(ctx), "error", "", start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.cache.get",
			Action:          "cache.get",
			ProviderCode:    req.ProviderCode,
			AppCode:         req.ProviderApp,
			Mode:            req.ProviderAuthMode,
			RedisKey:        "",
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
		})
		return
	}
	record, storage, ttl, err := s.tokenStore.fetch(r.Context(), ctx.CacheKey)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "load cache failed: %v", err)
		s.logMetric("cache", ctx.ProviderCode, ctx.AppID(), "error", "", start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.cache.get",
			Action:          "cache.get",
			Provider:        string(ctx.ProviderType),
			ProviderCode:    ctx.ProviderCode,
			AppCode:         ctx.AppCode,
			Mode:            ctx.ModeKey,
			RedisKey:        ctx.CacheKey,
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
		})
		return
	}
	storageBackend := s.cacheStorageBackend(storage)
	ttlSeconds := maxInt(int(ttl.Seconds()), 0)
	if record == nil {
		s.writeJSON(w, http.StatusOK, map[string]any{
			"cached":          false,
			"token_source":    firstOrDash(""),
			"cache_key":       ctx.CacheKey,
			"storage_backend": storageBackend,
			"ttl_seconds":     0,
			"refreshed_at":    nil,
		})
		s.logMetric("cache", ctx.ProviderCode, ctx.AppID(), "miss", "", start, nil)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.cache.get",
			Action:          "cache.get",
			Provider:        string(ctx.ProviderType),
			ProviderCode:    ctx.ProviderCode,
			AppCode:         ctx.AppCode,
			Mode:            ctx.ModeKey,
			RedisKey:        ctx.CacheKey,
			TTL:             0,
			StorageBackend:  storageBackend,
			TokenSource:     "-",
			APITokenSubject: subject,
			Status:          "miss",
		})
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]any{
		"cached":          true,
		"token_source":    firstOrDash(record.Source),
		"cache_key":       ctx.CacheKey,
		"appid":           record.AppID,
		"stored_at":       record.StoredAt,
		"refreshed_at":    record.StoredAt,
		"expire_at":       record.ExpireAt,
		"ttl_seconds":     ttlSeconds,
		"masked_token":    maskToken(record.AccessToken),
		"storage_backend": storageBackend,
	})
	s.logMetric("cache", ctx.ProviderCode, ctx.AppID(), "hit", record.Source, start, nil)
	s.logAudit(auditLogFields{
		Event:           "clienttoken.cache.get",
		Action:          "cache.get",
		Provider:        string(ctx.ProviderType),
		ProviderCode:    ctx.ProviderCode,
		AppCode:         ctx.AppCode,
		Mode:            ctx.ModeKey,
		RedisKey:        ctx.CacheKey,
		TTL:             ttlSeconds,
		StorageBackend:  storageBackend,
		TokenSource:     firstOrDash(record.Source),
		APITokenSubject: subject,
		Status:          "hit",
	})
}

func (s *clientTokenServer) handleCacheDelete(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	subject := apiTokenSubjectFromRequest(r)
	req := s.parseCacheRequest(r)
	ctx, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "%s", err.Error())
		s.logMetric("cache.delete", req.ProviderCode, ctxAppID(ctx), "error", "", start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.cache.delete",
			Action:          "cache.delete",
			ProviderCode:    req.ProviderCode,
			AppCode:         req.ProviderApp,
			Mode:            req.ProviderAuthMode,
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
		})
		return
	}
	record, storage, ttl, err := s.tokenStore.fetch(r.Context(), ctx.CacheKey)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "load cache failed: %v", err)
		s.logMetric("cache.delete", ctx.ProviderCode, ctx.AppID(), "error", "", start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.cache.delete",
			Action:          "cache.delete",
			Provider:        string(ctx.ProviderType),
			ProviderCode:    ctx.ProviderCode,
			AppCode:         ctx.AppCode,
			Mode:            ctx.ModeKey,
			RedisKey:        ctx.CacheKey,
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
		})
		return
	}
	if err := s.tokenStore.delete(r.Context(), ctx.CacheKey); err != nil {
		s.writeError(w, http.StatusBadGateway, "delete cache failed: %v", err)
		s.logMetric("cache.delete", ctx.ProviderCode, ctx.AppID(), "error", "", start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.cache.delete",
			Action:          "cache.delete",
			Provider:        string(ctx.ProviderType),
			ProviderCode:    ctx.ProviderCode,
			AppCode:         ctx.AppCode,
			Mode:            ctx.ModeKey,
			RedisKey:        ctx.CacheKey,
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
		})
		return
	}
	storageBackend := s.cacheStorageBackend(storage)
	ttlSeconds := maxInt(int(ttl.Seconds()), 0)
	tokenSource := ""
	if record != nil {
		tokenSource = record.Source
	}
	resp := map[string]any{
		"deleted":         true,
		"cached_before":   record != nil,
		"cache_key":       ctx.CacheKey,
		"provider_code":   ctx.ProviderCode,
		"provider_app":    ctx.AppCode,
		"provider_mode":   ctx.ModeKey,
		"storage_backend": storageBackend,
		"ttl_seconds":     ttlSeconds,
		"token_source":    firstOrDash(tokenSource),
		"refreshed_at":    nil,
	}
	if record != nil {
		resp["appid"] = record.AppID
		resp["refreshed_at"] = record.StoredAt
		resp["expire_at"] = record.ExpireAt
		resp["masked_token"] = maskToken(record.AccessToken)
	}
	s.writeJSON(w, http.StatusOK, resp)
	s.logMetric("cache.delete", ctx.ProviderCode, ctx.AppID(), "success", tokenSource, start, nil)
	s.logAudit(auditLogFields{
		Event:           "clienttoken.cache.delete",
		Action:          "cache.delete",
		Provider:        string(ctx.ProviderType),
		ProviderCode:    ctx.ProviderCode,
		AppCode:         ctx.AppCode,
		Mode:            ctx.ModeKey,
		RedisKey:        ctx.CacheKey,
		TTL:             ttlSeconds,
		StorageBackend:  storageBackend,
		TokenSource:     firstOrDash(tokenSource),
		APITokenSubject: subject,
		Status:          "success",
	})
}

func (s *clientTokenServer) parseCacheRequest(r *http.Request) cacheRequest {
	values := r.URL.Query()
	req := cacheRequest{
		ProviderCode:     values.Get("provider_code"),
		ProviderApp:      values.Get("provider_app"),
		ProviderAuthMode: values.Get("provider_auth_mode"),
		ConfigPath:       values.Get("config_path"),
	}
	if r.Method == http.MethodPost || r.Method == http.MethodDelete {
		var payload cacheRequest
		if err := decodeJSON(r.Body, &payload); err == nil {
			if payload.ProviderCode != "" {
				req.ProviderCode = payload.ProviderCode
			}
			if payload.ProviderApp != "" {
				req.ProviderApp = payload.ProviderApp
			}
			if payload.ProviderAuthMode != "" {
				req.ProviderAuthMode = payload.ProviderAuthMode
			}
			if payload.ConfigPath != "" {
				req.ConfigPath = payload.ConfigPath
			}
		}
	}
	return req
}

func (s *clientTokenServer) storageBackendLabel() string {
	if s.redis == nil {
		return "memory"
	}
	return "redis"
}

func (s *clientTokenServer) cacheStorageBackend(storage string) string {
	if strings.TrimSpace(storage) != "" {
		return strings.TrimSpace(storage)
	}
	return s.storageBackendLabel()
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
	subject := apiTokenSubjectFromRequest(r)
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
		s.writeError(w, http.StatusBadRequest, "%s", err.Error())
		s.logMetric("call", req.ProviderCode, ctxAppID(ctxMeta), "error", "", start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.call",
			Action:          req.Action,
			ProviderCode:    req.ProviderCode,
			AppCode:         req.ProviderApp,
			Mode:            req.ProviderAuthMode,
			RedisKey:        "",
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
		})
		return
	}
	needsRiskFields := ctxMeta.ProviderType == providerTypeByteDance
	riskReady := !needsRiskFields || ctxMeta.RiskControlReady
	riskExtra := ""
	if needsRiskFields && !riskReady {
		riskExtra = "risk_fields=disabled"
	}
	tokenRec, storageBackend, ttlDur, err := s.ensureToken(r.Context(), ctxMeta)
	if err != nil {
		if errors.Is(err, errTokenRefreshFailed) {
			s.writeError(w, http.StatusBadGateway, "token refresh failed: ttl_remaining=%d seconds", maxInt(int(ttlDur.Seconds()), 0))
		} else {
			s.writeError(w, http.StatusBadGateway, "ensure token failed: %v", err)
		}
		s.logMetric("call", ctxMeta.ProviderCode, ctxMeta.AppID(), "error", "", start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.call",
			Action:          req.Action,
			Provider:        string(ctxMeta.ProviderType),
			ProviderCode:    ctxMeta.ProviderCode,
			AppCode:         ctxMeta.AppCode,
			Mode:            ctxMeta.ModeKey,
			RedisKey:        ctxMeta.CacheKey,
			StorageBackend:  storageBackend,
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
			Extra:           riskExtra,
		})
		return
	}
	result, err := s.executeAPICall(r.Context(), ctxMeta, &req)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "%s", err.Error())
		s.logMetric("call", ctxMeta.ProviderCode, ctxMeta.AppID(), "error", tokenRec.Source, start, err)
		s.logAudit(auditLogFields{
			Event:           "clienttoken.call",
			Action:          req.Action,
			Provider:        string(ctxMeta.ProviderType),
			ProviderCode:    ctxMeta.ProviderCode,
			AppCode:         ctxMeta.AppCode,
			Mode:            ctxMeta.ModeKey,
			RedisKey:        ctxMeta.CacheKey,
			TTL:             maxInt(int(ttlDur.Seconds()), 0),
			StorageBackend:  storageBackend,
			TokenSource:     firstOrDash(tokenRec.Source),
			APITokenSubject: subject,
			Status:          "error",
			Err:             err,
			Extra:           riskExtra,
		})
		return
	}
	ttlSeconds := maxInt(int(ttlDur.Seconds()), 0)
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
	s.logAudit(auditLogFields{
		Event:           "clienttoken.call",
		Action:          req.Action,
		Provider:        string(ctxMeta.ProviderType),
		ProviderCode:    ctxMeta.ProviderCode,
		AppCode:         ctxMeta.AppCode,
		Mode:            ctxMeta.ModeKey,
		RedisKey:        ctxMeta.CacheKey,
		TTL:             ttlSeconds,
		StorageBackend:  storageBackend,
		TokenSource:     firstOrDash(tokenRec.Source),
		APITokenSubject: subject,
		Status:          "success",
		Extra:           riskExtra,
	})
}

func (s *clientTokenServer) executeAPICall(ctx context.Context, meta *providerContext, req *callRequest) (map[string]any, error) {
	switch meta.ProviderType {
	case providerTypeWechat:
		return s.executeWechatAPICall(ctx, meta, req)
	case providerTypeByteDance:
		return s.executeByteDanceAPICall(ctx, meta, req)
	default:
		return nil, fmt.Errorf("provider %s 暂不支持 API 调试", meta.ProviderCode)
	}
}

func (s *clientTokenServer) executeWechatAPICall(ctx context.Context, meta *providerContext, req *callRequest) (map[string]any, error) {
	if meta.WechatConfig == nil {
		return nil, fmt.Errorf("provider %s 配置缺失", meta.ProviderCode)
	}
	wechatClient, err := s.mediaX.CreateWechatClientTokenClient(meta.WechatConfig)
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
	result := map[string]any{}
	switch method {
	case http.MethodGet:
		_, err = wechatClient.GetBaseClient().HttpGet(ctx, req.Action, queryMap, nil, nil, &result)
	case http.MethodPost:
		_, err = wechatClient.GetBaseClient().HttpPost(ctx, req.Action, queryMap, bodyData, nil, &result)
	default:
		return nil, fmt.Errorf("暂不支持 method=%s", method)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *clientTokenServer) executeByteDanceAPICall(ctx context.Context, meta *providerContext, req *callRequest) (map[string]any, error) {
	if meta.ByteDanceConfig == nil {
		return nil, fmt.Errorf("provider %s 配置缺失", meta.ProviderCode)
	}
	douyinClient, err := s.mediaX.CreateByteDanceDouYinCTClient(meta.ByteDanceConfig)
	if err != nil {
		return nil, fmt.Errorf("create douyin client: %w", err)
	}
	if douyinClient == nil || douyinClient.ByteDanceClient == nil || douyinClient.ByteDanceClient.BaseClient == nil {
		return nil, fmt.Errorf("provider %s 缺少 base client", meta.ProviderCode)
	}
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = http.MethodPost
	}
	queryMap := parseQueryString(req.Query)
	var bodyData interface{}
	if strings.TrimSpace(req.Body) != "" {
		if err := json.Unmarshal([]byte(req.Body), &bodyData); err != nil {
			return nil, fmt.Errorf("parse body json: %w", err)
		}
	} else if method == http.MethodPost {
		bodyData = map[string]any{}
	}
	result := map[string]any{}
	baseClient := douyinClient.ByteDanceClient.BaseClient
	switch method {
	case http.MethodGet:
		_, err = baseClient.HttpGet(ctx, req.Action, queryMap, nil, nil, &result)
	case http.MethodPost:
		_, err = baseClient.HttpPost(ctx, req.Action, queryMap, bodyData, nil, &result)
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
		s.writeError(w, http.StatusBadRequest, "%s", err.Error())
		s.logMetric("message.validate", req.ProviderCode, ctxAppID(ctxMeta), "error", "", start, err)
		return
	}
	if ctxMeta.ProviderType != providerTypeWechat || ctxMeta.WechatConfig == nil {
		s.writeError(w, http.StatusBadRequest, "该 Provider 不支持消息验证")
		s.logMetric("message.validate", ctxMeta.ProviderCode, ctxMeta.AppID(), "error", "", start, fmt.Errorf("unsupported provider"))
		return
	}
	cred := ctxMeta.WechatConfig.Credentials()
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

func (s *clientTokenServer) handleCallbacks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		records := s.callbackStore.list()
		s.writeJSON(w, http.StatusOK, map[string]any{
			"records": records,
		})
	case http.MethodDelete:
		s.callbackStore.clear()
		s.writeJSON(w, http.StatusOK, map[string]any{"status": "cleared"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
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

func (s *clientTokenServer) refreshToken(ctx context.Context, meta *providerContext) (*tokenCacheRecord, error) {
	switch meta.ProviderType {
	case providerTypeWechat:
		return s.refreshWechatToken(ctx, meta)
	case providerTypeByteDance:
		return s.refreshByteDanceToken(ctx, meta)
	default:
		return nil, fmt.Errorf("provider %s 暂不支持刷新", meta.ProviderCode)
	}
}

func (s *clientTokenServer) refreshByteDanceToken(ctx context.Context, meta *providerContext) (*tokenCacheRecord, error) {
	if meta.ByteDanceConfig == nil {
		return nil, fmt.Errorf("provider %s 配置缺失", meta.ProviderCode)
	}
	douyinClient, err := s.mediaX.CreateByteDanceDouYinCTClient(meta.ByteDanceConfig)
	if err != nil {
		return nil, err
	}
	handler := douyinClient.ClientTokenHandler
	if handler == nil || handler.TokenHandler == nil {
		return nil, fmt.Errorf("provider %s 缺少 token handler", meta.ProviderCode)
	}
	if strings.TrimSpace(meta.CacheKey) != "" {
		handler.TokenHandler.SetCacheKey(meta.CacheKey)
	}
	tokenRes := &douyinresponse.ByteDanceAccessTokenRes{}
	if err := handler.TokenHandler.GetToken(ctx, true, tokenRes); err != nil {
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
		AccessToken: strings.TrimSpace(tokenRes.AccessToken),
		StoredAt:    now,
		ExpireAt:    now.Add(time.Duration(ttl) * time.Second),
		Source:      "refresh",
		ErrorCode:   tokenRes.ErrCode,
		ErrorMsg:    strings.TrimSpace(tokenRes.ErrMsg),
	}
	if record.AccessToken == "" {
		return record, fmt.Errorf("douyin token empty")
	}
	if tokenRes.ErrCode != 0 {
		return record, fmt.Errorf("douyin token error: errcode=%d errmsg=%s", tokenRes.ErrCode, record.ErrorMsg)
	}
	if err := s.tokenStore.save(ctx, meta.CacheKey, record, time.Duration(ttl)*time.Second); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *clientTokenServer) refreshWechatToken(ctx context.Context, meta *providerContext) (*tokenCacheRecord, error) {
	if meta.WechatConfig == nil {
		return nil, fmt.Errorf("provider %s 配置缺失", meta.ProviderCode)
	}
	wechatClient, err := s.mediaX.CreateWechatClientTokenClient(meta.WechatConfig)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(meta.CacheKey) != "" && wechatClient.AccessTokenHandler != nil && wechatClient.AccessTokenHandler.ClientTokenHandler != nil {
		wechatClient.AccessTokenHandler.ClientTokenHandler.SetCacheKey(meta.CacheKey)
	}
	tokenRes := &wechatresponse.WeChatAccessTokenRes{}
	if err := wechatClient.AccessTokenHandler.ClientTokenHandler.GetToken(ctx, true, tokenRes); err != nil {
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
		ErrorCode:   tokenRes.ErrCode,
		ErrorMsg:    strings.TrimSpace(tokenRes.ErrMsg),
	}
	if tokenRes.ErrCode != 0 {
		return record, fmt.Errorf("wechat token error: errcode=%d errmsg=%s", tokenRes.ErrCode, record.ErrorMsg)
	}
	if err := s.tokenStore.save(ctx, meta.CacheKey, record, time.Duration(ttl)*time.Second); err != nil {
		return nil, err
	}
	return record, nil
}

var errTokenRefreshFailed = errors.New("token refresh failed")

func wrapRefreshError(err error) error {
	if err == nil {
		return errTokenRefreshFailed
	}
	return fmt.Errorf("%w: %v", errTokenRefreshFailed, err)
}

func (s *clientTokenServer) ensureToken(ctx context.Context, meta *providerContext) (*tokenCacheRecord, string, time.Duration, error) {
	record, storage, ttl, err := s.tokenStore.fetch(ctx, meta.CacheKey)
	if err != nil {
		s.logger.WarnF("clienttoken: fetch cache failed: %v", err)
	}
	needRefresh := false
	if record == nil {
		needRefresh = true
	} else if ttl <= 0 {
		needRefresh = true
	} else if ttl <= time.Duration(meta.RefreshBeforeSeconds)*time.Second {
		needRefresh = true
	}
	if needRefresh {
		refreshed, refreshErr := s.refreshToken(ctx, meta)
		if refreshErr != nil {
			return nil, s.storageBackendLabel(), ttl, wrapRefreshError(refreshErr)
		}
		ttlDur := refreshed.ExpireAt.Sub(time.Now())
		if ttlDur < 0 {
			ttlDur = 0
		}
		return refreshed, s.storageBackendLabel(), ttlDur, nil
	}
	return record, s.cacheStorageBackend(storage), ttl, nil
}

func (s *clientTokenServer) writeError(w http.ResponseWriter, status int, format string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
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

type auditLogFields struct {
	Event           string
	Provider        string
	ProviderCode    string
	AppCode         string
	Mode            string
	Action          string
	RedisKey        string
	TTL             int
	StorageBackend  string
	TokenSource     string
	APITokenSubject string
	Status          string
	Err             error
	Extra           string
}

func (s *clientTokenServer) logAudit(fields auditLogFields) {
	errMsg := "-"
	if fields.Err != nil {
		errMsg = fields.Err.Error()
	}
	extra := firstOrDash(fields.Extra)
	ttl := fields.TTL
	if ttl < 0 {
		ttl = 0
	}
	s.logger.InfoF("clienttoken_event event=%s provider=%s provider_code=%s app_code=%s mode=%s action=%s redis_key=%s ttl_remaining=%d storage_backend=%s token_source=%s api_token_subject=%s status=%s extra=%s error=%s",
		firstOrDash(fields.Event),
		firstOrDash(fields.Provider),
		firstOrDash(fields.ProviderCode),
		firstOrDash(fields.AppCode),
		firstOrDash(fields.Mode),
		firstOrDash(fields.Action),
		firstOrDash(fields.RedisKey),
		ttl,
		firstOrDash(fields.StorageBackend),
		firstOrDash(fields.TokenSource),
		firstOrDash(fields.APITokenSubject),
		firstOrDash(fields.Status),
		extra,
		firstOrDash(errMsg),
	)
}

func apiTokenSubjectFromRequest(r *http.Request) string {
	if r == nil {
		return "-"
	}
	if auth := strings.TrimSpace(r.Header.Get("Authorization")); auth != "" {
		return "authorization:" + handler.MaskAuthorizationHeader(auth)
	}
	if header := strings.TrimSpace(r.Header.Get("X-API-Token")); header != "" {
		return "header:" + maskToken(header)
	}
	if r.URL != nil {
		if query := strings.TrimSpace(r.URL.Query().Get("api_token")); query != "" {
			return "query:" + maskToken(query)
		}
	}
	return "-"
}

func firstOrDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}
