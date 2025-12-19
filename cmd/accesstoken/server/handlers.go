package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	app "github.com/ArtisanCloud/MediaX/cmd/accesstoken/internal/app"
	mask "github.com/ArtisanCloud/MediaX/internal/accesstoken/handler"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

type tokenRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode,omitempty"`
	ConfigPath       string `json:"config_path"`
	AccessToken      string `json:"access_token"`
	AccessTokenTTL   int    `json:"access_token_ttl"`
}

type callRequest struct {
	app.Options
}

type oauthStartRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode"`
	ConfigPath       string `json:"config_path"`
}

type providerSelectionRequest struct {
	ProviderCode     string `json:"provider_code"`
	ProviderApp      string `json:"provider_app"`
	ProviderAuthMode string `json:"provider_auth_mode"`
	ConfigPath       string `json:"config_path"`
}

type flowReplayRequest struct {
	FlowID            string                    `json:"flow_id"`
	ProviderSelection *providerSelectionRequest `json:"provider_fallback"`
}

func (s *accessTokenServer) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req tokenRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid json: %v", err)
		return
	}
	ctx, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, source, sourceDetail, storedRec := s.resolveAccessTokenValue(ctx, req.AccessToken)
	if token == "" {
		s.writeError(w, http.StatusBadRequest, "missing access token for provider %s: provide access_token or configure env/config token", ctx.displayName())
		return
	}
	ttl := req.AccessTokenTTL
	if ttl <= 0 && storedRec != nil && storedRec.ExpiresIn > 0 {
		ttl = storedRec.ExpiresIn
	}
	if ttl <= 0 {
		ttl = app.DefaultAccessTokenTTLSeconds
	}

	clientCfg := ctx.clientConfig()
	httpDebug := false
	if clientCfg != nil && clientCfg.BaseConfig != nil {
		httpDebug = clientCfg.HttpDebug
	}

	oauthKey := extractOauthKey(ctx.Mode)

	resp := map[string]any{
		"access_token":        token,
		"masked_token":        app.MaskToken(token),
		"token_source":        source,
		"access_token_ttl":    ttl,
		"oauth_key":           oauthKey,
		"config_path":         ctx.ConfigPath,
		"provider":            ctx.displayName(),
		"provider_code":       ctx.ProviderCode,
		"provider_app":        ctx.AppCode,
		"provider_auth_mode":  ctx.ModeKey,
		"http_debug":          httpDebug,
		"storage_backend":     s.storageBackend,
		"token_source_detail": sourceDetail,
	}
	if storedRec != nil {
		s.enrichFlowMetadata(storedRec)
		resp["flow_id"] = storedRec.FlowID
		resp["flow_expire_at"] = storedRec.FlowExpireAt
		resp["flow_ttl_seconds"] = storedRec.FlowTTLSeconds
		resp["storage_backend"] = storedRec.StorageBackend
		flowID := storedRec.FlowID
		s.logFlowAction("token.resolve", ctx, flowID, source, sourceDetail)
	} else {
		s.logFlowAction("token.resolve", ctx, "", source, sourceDetail)
	}
	s.writeJSON(w, http.StatusOK, resp)
}
func (s *accessTokenServer) handleCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var request callRequest
	if err := decodeJSON(r.Body, &request); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid json: %v", err)
		return
	}
	opts := &request.Options
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, err := s.resolveProviderConfig(opts.ProviderCode, opts.ProviderApp, opts.AuthMode, opts.ConfigPath)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if ctx.ProviderCode != providerGoogleYouTube {
		s.writeError(w, http.StatusBadRequest, "provider %s 暂未开放 API 调试，仅支持 Google YouTube", ctx.displayName())
		return
	}

	token, source, sourceDetail, storedRec := s.resolveAccessTokenValue(ctx, opts.AccessToken)
	if token == "" {
		s.writeError(w, http.StatusBadRequest, "missing access token: provide access_token 或配置 GOOGLE_YOUTUBE_ACCESS_TOKEN / oauth.access_token")
		return
	}
	opts.AccessToken = token
	opts.TokenSource = sourceDetail
	if storedRec != nil && storedRec.ExpiresIn > 0 {
		opts.AccessTokenTTL = storedRec.ExpiresIn
	} else if opts.AccessTokenTTL <= 0 {
		opts.AccessTokenTTL = app.DefaultAccessTokenTTLSeconds
	}

	cfg := ctx.Youtube
	configPath := ctx.ConfigPath
	cfg.GetOAuthToken = func(key string, refresh bool) object.HashMap {
		return object.HashMap{
			"access_token": token,
			"expires_in":   float64(opts.AccessTokenTTL),
		}
	}

	opts.ProviderCode = ctx.ProviderCode
	opts.ProviderApp = ctx.AppCode
	if opts.AuthMode == "" {
		opts.AuthMode = ctx.ModeKey
	}

	app.LogInvocation(s.logger, opts, cfg.OauthKey)
	ytClient, err := s.mediaX.CreateGoogleYouTubeACClient(cfg)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "create youtube client: %v", err)
		return
	}
	data, err := app.ExecuteAction(r.Context(), s.logger, ytClient, opts)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	resp := map[string]any{
		"result":              data,
		"token_source":        source,
		"token_source_detail": sourceDetail,
		"masked_token":        app.MaskToken(opts.AccessToken),
		"oauth_key":           cfg.OauthKey,
		"config_path":         configPath,
		"provider":            ctx.displayName(),
		"provider_code":       ctx.ProviderCode,
		"provider_app":        ctx.AppCode,
		"provider_auth_mode":  ctx.ModeKey,
		"access_token_ttl":    opts.AccessTokenTTL,
		"storage_backend":     s.storageBackend,
		"action":              opts.Action,
		"part":                opts.Part,
		"request_timestamp":   time.Now().UTC(),
	}
	if storedRec != nil {
		s.enrichFlowMetadata(storedRec)
		resp["flow_id"] = storedRec.FlowID
		resp["flow_expire_at"] = storedRec.FlowExpireAt
		resp["flow_ttl_seconds"] = storedRec.FlowTTLSeconds
	}
	flowID := ""
	if storedRec != nil {
		flowID = storedRec.FlowID
		resp["storage_backend"] = storedRec.StorageBackend
	}
	s.logFlowAction("token.call", ctx, flowID, source, sourceDetail)
	s.writeJSON(w, http.StatusOK, resp)
}
func (s *accessTokenServer) handleOAuthStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req oauthStartRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		s.logFlowAction("oauth.start.error", nil, "", "invalid_json", err.Error())
		s.writeError(w, http.StatusBadRequest, "invalid json: %v", err)
		return
	}
	if strings.TrimSpace(req.ProviderCode) == "" || strings.TrimSpace(req.ProviderApp) == "" || strings.TrimSpace(req.ProviderAuthMode) == "" || strings.TrimSpace(req.ConfigPath) == "" {
		s.logFlowAction("oauth.start.error", nil, "", "invalid_request", "provider_code/provider_app/provider_auth_mode/config_path 不能为空")
		s.writeError(w, http.StatusBadRequest, "provider_code/provider_app/provider_auth_mode/config_path 不能为空")
		return
	}
	ctx, err := s.resolveProviderConfig(req.ProviderCode, req.ProviderApp, req.ProviderAuthMode, req.ConfigPath)
	if err != nil {
		s.logFlowAction("oauth.start.error", ctx, "", "resolve_provider", err.Error())
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	authURL, state, err := s.buildOAuthAuthorizeURL(ctx)
	if err != nil {
		s.logFlowAction("oauth.start.error", ctx, "", "build_authorize_url", err.Error())
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	flowID := fmt.Sprintf("oauth-%s", state)
	oauthKey := extractOauthKey(ctx.Mode)
	expiresIn := s.flowTTLSeconds
	expireAt := time.Now().UTC().Add(time.Duration(expiresIn) * time.Second)
	response := map[string]any{
		"flow_id":         flowID,
		"state":           state,
		"authorize_url":   authURL,
		"auth_url":        authURL,
		"expires_in":      expiresIn,
		"expire_at":       expireAt,
		"listen_addr":     s.listenAddr,
		"storage_backend": s.storageBackend,
		"oauth_key":       oauthKey,
		"config_path":     ctx.ConfigPath,
		"provider": map[string]string{
			"code":      ctx.ProviderCode,
			"app":       ctx.AppCode,
			"auth_mode": ctx.ModeKey,
		},
	}
	s.logFlowAction("oauth.start", ctx, flowID, "pending", state)
	s.writeJSON(w, http.StatusOK, response)
}

func (s *accessTokenServer) handleListFlows(w http.ResponseWriter, r *http.Request) {
	provider := strings.TrimSpace(r.URL.Query().Get("provider_code"))
	appCode := strings.TrimSpace(r.URL.Query().Get("provider_app"))
	mode := strings.TrimSpace(r.URL.Query().Get("provider_auth_mode"))
	flowID := strings.TrimSpace(r.URL.Query().Get("flow_id"))
	limit := 20
	if v := strings.TrimSpace(r.URL.Query().Get("limit")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			if parsed > 200 {
				parsed = 200
			}
			limit = parsed
		}
	}
	cursor := strings.TrimSpace(r.URL.Query().Get("cursor"))
	var (
		records    []*oauthTokenRecord
		nextCursor string
		err        error
	)
	if flowID != "" {
		rec, lookupErr := s.lookupFlowByID(flowID)
		if lookupErr == nil && s.matchesProvider(rec, provider, appCode, mode) {
			records = []*oauthTokenRecord{rec}
		} else {
			records = []*oauthTokenRecord{}
		}
	} else if s.redis != nil {
		records, nextCursor, err = s.scanFlowRecords(r.Context(), provider, appCode, mode, limit, cursor)
		if err != nil {
			s.writeError(w, http.StatusInternalServerError, "scan flow indexes failed: %v", err)
			return
		}
		if len(records) == 0 {
			records = s.listOAuthTokenRecords(provider, appCode, mode)
		}
	} else {
		records = s.listOAuthTokenRecords(provider, appCode, mode)
	}
	if len(records) > 1 {
		sort.Slice(records, func(i, j int) bool {
			if records[i].StoredAt.Equal(records[j].StoredAt) {
				return records[i].FlowID > records[j].FlowID
			}
			return records[i].StoredAt.After(records[j].StoredAt)
		})
	}
	if len(records) > limit {
		records = records[:limit]
	}
	items := make([]map[string]any, 0, len(records))
	for _, rec := range records {
		if rec == nil {
			continue
		}
		items = append(items, s.buildFlowSummary(rec))
	}
	resp := map[string]any{
		"flows": items,
	}
	if nextCursor != "" {
		resp["next_cursor"] = nextCursor
	}
	s.writeJSON(w, http.StatusOK, resp)
}

func (s *accessTokenServer) handleFlowReplay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req flowReplayRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid json: %v", err)
		return
	}
	flowID := strings.TrimSpace(req.FlowID)
	if flowID == "" {
		s.writeError(w, http.StatusBadRequest, "flow_id required")
		return
	}
	rec, err := s.lookupFlowByID(flowID)
	if err != nil {
		switch {
		case errors.Is(err, errFlowExpired):
			s.writeJSON(w, http.StatusNotFound, map[string]any{
				"error":   "FLOW_EXPIRED",
				"message": "Flow 已过期，请重新授权。",
			})
		case errors.Is(err, errFlowNotFound):
			s.writeJSON(w, http.StatusNotFound, map[string]any{
				"error":   "FLOW_NOT_FOUND",
				"message": "Flow 不存在或已被清理。",
			})
		default:
			s.writeError(w, http.StatusInternalServerError, "load flow failed: %v", err)
		}
		return
	}
	payload := s.buildFlowPayload(rec)
	s.writeJSON(w, http.StatusOK, map[string]any{
		"flow_id":         rec.FlowID,
		"status":          s.flowStatus(rec),
		"storage_backend": rec.StorageBackend,
		"payload":         payload,
	})
}

func (s *accessTokenServer) handleListFlowIndexes(w http.ResponseWriter, r *http.Request) {
	if s.redis == nil {
		s.writeError(w, http.StatusServiceUnavailable, "redis 缓存未启用，无法列出 Flow")
		return
	}
	limit := 50
	if v := strings.TrimSpace(r.URL.Query().Get("limit")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			if parsed > 200 {
				parsed = 200
			}
			limit = parsed
		}
	}
	var cursor uint64
	if cur := strings.TrimSpace(r.URL.Query().Get("cursor")); cur != "" {
		if parsed, err := strconv.ParseUint(cur, 10, 64); err == nil {
			cursor = parsed
		}
	}
	keys, nextCursor, err := s.redis.Scan(r.Context(), cursor, flowIndexPrefix+"*", int64(limit)).Result()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "scan redis flow index failed: %v", err)
		return
	}
	tokens := make([]*oauthTokenRecord, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
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
		if rec := s.fetchTokenByFlowID(flowID); rec != nil {
			tokens = append(tokens, rec)
		}
	}
	if len(tokens) > 1 {
		sort.Slice(tokens, func(i, j int) bool {
			if tokens[i].StoredAt.Equal(tokens[j].StoredAt) {
				return tokens[i].FlowID > tokens[j].FlowID
			}
			return tokens[i].StoredAt.After(tokens[j].StoredAt)
		})
	}
	s.writeJSON(w, http.StatusOK, map[string]any{
		"tokens": tokens,
		"cursor": nextCursor,
	})
}

func (s *accessTokenServer) handleDebugCallback(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	queryValues := r.URL.Query()
	code := strings.TrimSpace(queryValues.Get("code"))
	state := strings.TrimSpace(queryValues.Get("state"))
	queryValues.Del("code")
	queryValues.Del("state")
	maskedQuery := mask.MaskSensitiveQuery(r.URL.RawQuery)
	record := callbackRecord{
		Timestamp:    time.Now(),
		Method:       r.Method,
		Query:        maskedQuery,
		Headers:      flattenHeaders(r.Header),
		Body:         mask.MaskCallbackBody(string(body), r.Header.Get("Content-Type")),
		ProviderCode: queryValues.Get("provider"),
		FlowID:       queryValues.Get("flow_id"),
		State:        state,
	}
	if record.FlowID == "" && state != "" {
		record.FlowID = fmt.Sprintf("oauth-%s", state)
	}
	s.callbackStore.append(record)
	var message string
	if code != "" && state != "" {
		payload := record
		if tokenRec, err := s.completeOAuthFlow(r.Context(), state, code, &payload); err != nil {
			message = fmt.Sprintf("授权失败: %v", err)
		} else if tokenRec != nil {
			message = fmt.Sprintf("授权成功，flow_id=%s，可在调试页刷新授权记录后复用该 Token。", tokenRec.FlowID)
		}
	}
	if acceptsHTML(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if message == "" {
			message = "回调已记录，可返回调试页。"
		}
		fmt.Fprintf(w, "<html><body><h2>AccessToken 调试台</h2><p>%s</p><p><a href=\"/debug\">返回调试页</a></p></body></html>", message)
		return
	}
	resp := map[string]any{"status": "ok"}
	if message != "" {
		resp["message"] = message
	}
	s.writeJSON(w, http.StatusOK, resp)
}

func (s *accessTokenServer) handleListCallbacks(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]any{
		"records": s.callbackStore.list(),
	})
}

func (s *accessTokenServer) handleClearCallbacks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s.callbackStore.clear()
	s.writeJSON(w, http.StatusOK, map[string]any{"status": "cleared"})
}

func (s *accessTokenServer) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.logger.ErrorF("accesstoken-server: write json failed: %v", err)
	}
}

func (s *accessTokenServer) writeError(w http.ResponseWriter, status int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	s.writeJSON(w, status, map[string]any{"error": msg})
}

func decodeJSON(r io.Reader, dst any) error {
	dec := json.NewDecoder(io.LimitReader(r, 1<<20))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func flattenHeaders(h http.Header) map[string]string {
	if len(h) == 0 {
		return nil
	}
	out := make(map[string]string, len(h))
	for k, vals := range h {
		out[k] = mask.MaskHeaderValue(k, strings.Join(vals, ","))
	}
	return out
}

func acceptsHTML(r *http.Request) bool {
	accept := strings.ToLower(r.Header.Get("Accept"))
	return strings.Contains(accept, "text/html")
}
