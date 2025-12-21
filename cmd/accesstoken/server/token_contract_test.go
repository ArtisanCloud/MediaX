package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	app "github.com/ArtisanCloud/MediaX/cmd/accesstoken/internal/app"
)

const testAccessTokenConfig = `
access_token_providers:
  providers:
    - code: "bilibili"
      name: "BiliBili"
      apps:
        - code: "content_center"
          name: "Content Center"
          provider_code: "bilbili"
          api_version: "v1"
          auth_modes:
            - key: "default"
              label: "默认"
              bilbili_config:
                api_url: "https://member.bilibili.com"
                proxy_api_url: ""
                timeout: 5
                http_debug: false
                oauth:
                  oauth_url: "https://member.bilibili.com/oauth2/authorize"
                  access_token_url: "https://member.bilibili.com/oauth2/token"
                  client_id: "demo-client"
                  client_secret: "demo-secret"
                  redirect_url: "http://localhost:7071/debug/callback"
                  scope: "arc_base"
                  access_token: "config-access-token"
    - code: "redbook"
      name: "RedBook"
      apps:
        - code: "juguang"
          name: "JuGuang"
          provider_code: "redbook_juguang"
          api_version: "v1"
          auth_modes:
            - key: "default"
              label: "默认"
              redbook_juguang_config:
                api_url: "https://adapi.xiaohongshu.com"
                proxy_api_url: ""
                timeout: 5
                http_debug: false
                oauth:
                  oauth_url: "https://ad.xiaohongshu.com/oauth2/authorize"
                  access_token_url: "https://ad.xiaohongshu.com/oauth2/token"
                  client_id: "redbook-client"
                  client_secret: "redbook-secret"
                  redirect_url: "http://localhost:7071/debug/callback"
                  scope: "notes.read,notes.write"
                oauth_key: "redbook-default"
    - code: "byte_dance"
      name: "ByteDance"
      apps:
        - code: "default"
          app_key: "douyin"
          name: "DouYin"
          provider_code: "byte_dance_douyin"
          api_version: "v1"
          auth_modes:
            - key: "default"
              label: "默认 OAuth"
              byte_dance_douyin_config:
                api_url: "https://open.douyin.com"
                proxy_api_url: ""
                timeout: 5
                http_debug: false
                oauth:
                  oauth_url: "https://open.douyin.com/platform/oauth/connect"
                  access_token_url: "https://open.douyin.com/oauth/access_token/"
                  client_id: "douyin-client"
                  client_secret: "douyin-secret"
                  redirect_url: "http://localhost:7071/debug/callback"
                  scope: "video.list,im.message.send"
                oauth_key: "douyin-default"
`

const testAccessTokenConfigNoDefault = `
access_token_providers:
  providers:
    - code: "bilibili"
      name: "BiliBili"
      apps:
        - code: "content_center"
          name: "Content Center"
          provider_code: "bilbili"
          api_version: "v1"
          auth_modes:
            - key: "default"
              label: "默认"
              bilbili_config:
                api_url: "https://member.bilibili.com"
                proxy_api_url: ""
                timeout: 5
                http_debug: false
                oauth:
                  oauth_url: "https://member.bilibili.com/oauth2/authorize"
                  access_token_url: "https://member.bilibili.com/oauth2/token"
                  client_id: "demo-client"
                  client_secret: "demo-secret"
                  redirect_url: "http://localhost:7071/debug/callback"
                  scope: "arc_base"
`

func testAccessTokenConfigWithTokenURL(tokenURL string) string {
	return strings.Replace(testAccessTokenConfig, "https://member.bilibili.com/oauth2/token", tokenURL, 1)
}

func testAccessTokenConfigWithRedbookTokenURL(tokenURL string) string {
	return strings.Replace(testAccessTokenConfig, "https://ad.xiaohongshu.com/oauth2/token", tokenURL, 1)
}

func testAccessTokenConfigWithDouyinTokenURL(tokenURL string) string {
	return strings.Replace(testAccessTokenConfig, "https://open.douyin.com/oauth/access_token/", tokenURL, 1)
}

func testAccessTokenConfigWithDouyinAPI(configYAML, apiURL string) string {
	return strings.Replace(configYAML, "https://open.douyin.com", apiURL, 1)
}

func testAccessTokenConfigWithoutDouyinConfig() string {
	return strings.Replace(testAccessTokenConfig, `
            - key: "default"
              label: "默认 OAuth"
              byte_dance_douyin_config:
                api_url: "https://open.douyin.com"
                proxy_api_url: ""
                timeout: 5
                http_debug: false
                oauth:
                  oauth_url: "https://open.douyin.com/platform/oauth/connect"
                  access_token_url: "https://open.douyin.com/oauth/access_token/"
                  client_id: "douyin-client"
                  client_secret: "douyin-secret"
                  redirect_url: "http://localhost:7071/debug/callback"
                  scope: "video.list,im.message.send"
                oauth_key: "douyin-default"
`, `
            - key: "default"
              label: "默认 OAuth"
              byte_dance_douyin_config: null
`, 1)
}

func testAccessTokenConfigWithoutDouyinScope() string {
	return strings.Replace(testAccessTokenConfig, `scope: "video.list,im.message.send"`, `scope: ""`, 1)
}

type tokenAPIResponse struct {
	AccessToken       string    `json:"access_token"`
	MaskedToken       string    `json:"masked_token"`
	TokenSource       string    `json:"token_source"`
	TokenSourceDetail string    `json:"token_source_detail"`
	AccessTokenTTL    int       `json:"access_token_ttl"`
	FlowID            string    `json:"flow_id"`
	FlowExpireAt      time.Time `json:"flow_expire_at"`
	FlowTTLSeconds    int       `json:"flow_ttl_seconds"`
	StorageBackend    string    `json:"storage_backend"`
}

func TestSaveOAuthTokenRecordHonorsFlowTTL(t *testing.T) {
	t.Setenv("ACCESSTOKEN_FLOW_TTL_SECONDS", "2")
	server := newTestAccessTokenServer(t, testAccessTokenConfig)

	expired := &oauthTokenRecord{
		ProviderCode: "bilbili",
		ProviderApp:  "content_center",
		AuthMode:     "default",
		FlowID:       "oauth-expired",
		AccessToken:  "expired-flow-token",
		StoredAt:     time.Now().Add(-5 * time.Second).UTC(),
	}
	if err := server.saveOAuthTokenRecord(expired); err != nil {
		t.Fatalf("saveOAuthTokenRecord expired: %v", err)
	}
	if got := server.fetchTokenByFlowID("oauth-expired"); got != nil {
		t.Fatalf("expected expired Flow to be filtered, got %+v", got)
	}

	fresh := &oauthTokenRecord{
		ProviderCode: "bilbili",
		ProviderApp:  "content_center",
		AuthMode:     "default",
		FlowID:       "oauth-active",
		AccessToken:  "active-flow-token",
		ExpiresIn:    900,
		StoredAt:     time.Now().UTC(),
	}
	if err := server.saveOAuthTokenRecord(fresh); err != nil {
		t.Fatalf("saveOAuthTokenRecord fresh: %v", err)
	}
	stored := server.fetchTokenByFlowID("oauth-active")
	if stored == nil {
		t.Fatalf("expected Flow to be retrievable before TTL expiry")
	}
	if stored.FlowTTLSeconds != fresh.ExpiresIn {
		t.Fatalf("FlowTTLSeconds = %d, want %d", stored.FlowTTLSeconds, fresh.ExpiresIn)
	}
	wantExpire := fresh.StoredAt.Add(time.Duration(fresh.ExpiresIn) * time.Second)
	if !stored.FlowExpireAt.Equal(wantExpire) {
		t.Fatalf("FlowExpireAt = %v, want %v", stored.FlowExpireAt, wantExpire)
	}
	if stored.StorageBackend != server.storageBackend {
		t.Fatalf("StorageBackend = %s, want %s", stored.StorageBackend, server.storageBackend)
	}
}

func TestHandleTokenPrefersRequestPayload(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	payload := defaultTokenPayload(server)
	payload["access_token"] = "direct-token-12345"

	status, resp := invokeTokenEndpoint(t, server, payload)
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if resp.TokenSource != "request" || resp.TokenSourceDetail != "payload" {
		t.Fatalf("unexpected token source: %+v", resp)
	}
	if resp.MaskedToken != app.MaskToken("direct-token-12345") {
		t.Fatalf("masked token mismatch: %s", resp.MaskedToken)
	}
	if resp.AccessTokenTTL != app.DefaultAccessTokenTTLSeconds {
		t.Fatalf("ttl = %d, want %d", resp.AccessTokenTTL, app.DefaultAccessTokenTTLSeconds)
	}
	if resp.StorageBackend != server.storageBackend {
		t.Fatalf("storage backend mismatch: %s", resp.StorageBackend)
	}
}

func TestHandleTokenFallsBackToEnvVariable(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	t.Setenv("BILIBILI_ACCESS_TOKEN", "env-token-abcde")
	payload := defaultTokenPayload(server)

	status, resp := invokeTokenEndpoint(t, server, payload)
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if resp.TokenSource != "env" || resp.TokenSourceDetail != "BILIBILI_ACCESS_TOKEN" {
		t.Fatalf("unexpected env source fields: %+v", resp)
	}
	if resp.AccessToken != "env-token-abcde" {
		t.Fatalf("access token mismatch: %s", resp.AccessToken)
	}
}

func TestHandleTokenFallsBackToConfigValue(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	payload := defaultTokenPayload(server)

	status, resp := invokeTokenEndpoint(t, server, payload)
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if resp.TokenSource != "config" || resp.TokenSourceDetail != "config" {
		t.Fatalf("expected config token source, got %+v", resp)
	}
	if resp.AccessToken != "config-access-token" {
		t.Fatalf("access token mismatch, got %s", resp.AccessToken)
	}
}

func TestHandleTokenFallsBackToStoredFlow(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfigNoDefault)
	rec := &oauthTokenRecord{
		ProviderCode: "bilbili",
		ProviderApp:  "content_center",
		AuthMode:     "default",
		FlowID:       "oauth-demo-flow",
		AccessToken:  "flow-token-xyz",
		ExpiresIn:    7200,
		StoredAt:     time.Now().UTC(),
	}
	if err := server.saveOAuthTokenRecord(rec); err != nil {
		t.Fatalf("saveOAuthTokenRecord: %v", err)
	}
	payload := defaultTokenPayload(server)

	status, resp := invokeTokenEndpoint(t, server, payload)
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if resp.TokenSource != "flow" || resp.TokenSourceDetail != "oauth-demo-flow" {
		t.Fatalf("expected flow source, got %+v", resp)
	}
	if resp.FlowID != "oauth-demo-flow" {
		t.Fatalf("flow_id mismatch: %s", resp.FlowID)
	}
	if resp.AccessTokenTTL != rec.ExpiresIn {
		t.Fatalf("ttl = %d, want %d", resp.AccessTokenTTL, rec.ExpiresIn)
	}
	if resp.StorageBackend != server.storageBackend {
		t.Fatalf("storage backend mismatch: %s", resp.StorageBackend)
	}
	if resp.FlowTTLSeconds != rec.ExpiresIn {
		t.Fatalf("flow ttl seconds mismatch: %d vs %d", resp.FlowTTLSeconds, rec.ExpiresIn)
	}
	if resp.FlowExpireAt.IsZero() {
		t.Fatalf("flow_expire_at should not be zero")
	}
}

func TestInitProvidersHidesDouyinWithoutConfig(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfigWithoutDouyinConfig())
	for _, provider := range server.providers {
		if provider.Code == "byte_dance" {
			t.Fatalf("byte_dance provider should not be exposed when config is missing: %+v", provider)
		}
	}
}

func TestDouyinCallRefreshesTokenAutomatically(t *testing.T) {
	var refreshCalls atomic.Int64
	var actionCalls atomic.Int64
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/access_token/":
			refreshCalls.Add(1)
			_ = r.ParseForm()
			if got := r.FormValue("refresh_token"); got != "refresh-token-xyz" {
				t.Fatalf("unexpected refresh token: %s", got)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"access_token":"refreshed-token","refresh_token":"refresh-token-new","expires_in":7200}`)
		case "/douyin/video/list":
			actionCalls.Add(1)
			if got := r.URL.Query().Get("access_token"); got != "refreshed-token" {
				t.Fatalf("expected refreshed token, got %s", got)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"message":"ok"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	configYAML := testAccessTokenConfigWithDouyinTokenURL(mockServer.URL + "/oauth/access_token/")
	configYAML = testAccessTokenConfigWithDouyinAPI(configYAML, mockServer.URL)
	server := newTestAccessTokenServer(t, configYAML)

	rec := &oauthTokenRecord{
		ProviderCode: providerByteDanceDouYin,
		ProviderApp:  "douyin",
		AuthMode:     "default",
		FlowID:       "flow-refresh",
		AccessToken:  "expired-token",
		RefreshToken: "refresh-token-xyz",
		ExpiresIn:    7200,
		StoredAt:     time.Now().Add(-7190 * time.Second).UTC(),
		Source:       "authorization_code",
	}
	rec.TokenExpireAt = rec.StoredAt.Add(time.Duration(rec.ExpiresIn) * time.Second)
	rec.FlowTTLSeconds = rec.ExpiresIn
	if err := server.saveOAuthTokenRecord(rec); err != nil {
		t.Fatalf("saveOAuthTokenRecord: %v", err)
	}

	payload := defaultDouyinCallPayload(server)
	status, resp := invokeCallEndpoint(t, server, payload)
	if status != http.StatusOK {
		t.Fatalf("status=%d resp=%v", status, resp)
	}
	if refreshCalls.Load() == 0 {
		t.Fatalf("expected refresh token endpoint to be called")
	}
	if actionCalls.Load() != 1 {
		t.Fatalf("expected exactly one douyin action call, got %d", actionCalls.Load())
	}
	if got := resp["flow_id"]; got != rec.FlowID {
		t.Fatalf("flow_id mismatch: %v", got)
	}
	if got := resp["token_source"]; got != "flow" {
		t.Fatalf("token_source expected flow, got %v", got)
	}
	if retry := numberValue(resp["retry_count"]); retry != 0 {
		t.Fatalf("retry_count=%d, want 0", retry)
	}
	result, ok := resp["result"].(map[string]any)
	if !ok || result["message"] != "ok" {
		t.Fatalf("unexpected result payload: %v", resp["result"])
	}
	stored := server.fetchTokenByFlowID(rec.FlowID)
	if stored == nil {
		t.Fatalf("flow missing after refresh")
	}
	if stored.AccessToken != "refreshed-token" {
		t.Fatalf("access token not updated: %+v", stored)
	}
	if stored.Source != "refresh_token" {
		t.Fatalf("expected source refresh_token, got %s", stored.Source)
	}
	if stored.LastRefreshAt.IsZero() {
		t.Fatalf("LastRefreshAt not set")
	}
}

func TestDouyinCallRefreshFailureRequiresReauth(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/access_token/" {
			http.Error(w, "invalid refresh", http.StatusBadRequest)
			return
		}
		http.NotFound(w, r)
	}))
	defer mockServer.Close()

	configYAML := testAccessTokenConfigWithDouyinTokenURL(mockServer.URL + "/oauth/access_token/")
	configYAML = testAccessTokenConfigWithDouyinAPI(configYAML, mockServer.URL)
	server := newTestAccessTokenServer(t, configYAML)

	rec := &oauthTokenRecord{
		ProviderCode: providerByteDanceDouYin,
		ProviderApp:  "douyin",
		AuthMode:     "default",
		FlowID:       "flow-reauth",
		AccessToken:  "expired-token",
		RefreshToken: "refresh-token-xyz",
		ExpiresIn:    7200,
		StoredAt:     time.Now().Add(-7190 * time.Second).UTC(),
	}
	rec.TokenExpireAt = rec.StoredAt.Add(time.Duration(rec.ExpiresIn) * time.Second)
	rec.FlowTTLSeconds = rec.ExpiresIn
	if err := server.saveOAuthTokenRecord(rec); err != nil {
		t.Fatalf("saveOAuthTokenRecord: %v", err)
	}

	payload := defaultDouyinCallPayload(server)
	status, resp := invokeCallEndpoint(t, server, payload)
	if status != http.StatusUnauthorized {
		t.Fatalf("status=%d resp=%v", status, resp)
	}
	if resp["error"] != "need reauth" {
		t.Fatalf("expected need reauth error, got %v", resp["error"])
	}
	if stored := server.fetchTokenByFlowID(rec.FlowID); stored != nil {
		t.Fatalf("flow should be invalidated on refresh failure")
	}
}

func TestDouyinCallEnforcesRateLimit(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/douyin/video/list" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ok":true}`)
			return
		}
		if r.URL.Path == "/oauth/access_token/" {
			fmt.Fprint(w, `{"access_token":"token","expires_in":3600}`)
			return
		}
		http.NotFound(w, r)
	}))
	defer mockServer.Close()

	configYAML := testAccessTokenConfigWithDouyinTokenURL(mockServer.URL + "/oauth/access_token/")
	configYAML = testAccessTokenConfigWithDouyinAPI(configYAML, mockServer.URL)
	server := newTestAccessTokenServer(t, configYAML)

	payload := defaultDouyinCallPayload(server)
	payload["access_token"] = "direct-token"

	status, resp := invokeCallEndpoint(t, server, payload)
	if status != http.StatusOK {
		t.Fatalf("first call status=%d resp=%v", status, resp)
	}
	status, resp = invokeCallEndpoint(t, server, payload)
	if status != http.StatusTooManyRequests {
		t.Fatalf("rate limit status=%d resp=%v", status, resp)
	}
	if errMsg, _ := resp["error"].(string); !strings.Contains(errMsg, "QPS") && !strings.Contains(errMsg, "限制") {
		t.Fatalf("unexpected rate limit error: %v", errMsg)
	}
}

func TestDouyinCallRetriesOnServerError(t *testing.T) {
	var attempts atomic.Int64
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/douyin/video/list" {
			count := attempts.Add(1)
			if count < 3 {
				http.Error(w, "server busy", http.StatusBadGateway)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ok":true}`)
			return
		}
		fmt.Fprint(w, `{"access_token":"token","expires_in":3600}`)
	}))
	defer mockServer.Close()

	configYAML := testAccessTokenConfigWithDouyinTokenURL(mockServer.URL + "/oauth/access_token/")
	configYAML = testAccessTokenConfigWithDouyinAPI(configYAML, mockServer.URL)
	server := newTestAccessTokenServer(t, configYAML)

	payload := defaultDouyinCallPayload(server)
	payload["access_token"] = "direct-token"

	start := time.Now()
	status, resp := invokeCallEndpoint(t, server, payload)
	elapsed := time.Since(start)
	if status != http.StatusOK {
		t.Fatalf("status=%d resp=%v", status, resp)
	}
	if attempts.Load() != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts.Load())
	}
	if retry := numberValue(resp["retry_count"]); retry != 2 {
		t.Fatalf("retry_count=%d, want 2", retry)
	}
	if backoff := numberValue(resp["last_backoff_ms"]); backoff < 400 {
		t.Fatalf("last_backoff_ms=%d, want >=400", backoff)
	}
	if elapsed < 500*time.Millisecond {
		t.Fatalf("expected retries with backoff, got duration %v", elapsed)
	}
}

func TestDouyinCallRequiresScope(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfigWithoutDouyinScope())
	payload := defaultDouyinCallPayload(server)
	payload["access_token"] = "direct-token"

	status, resp := invokeCallEndpoint(t, server, payload)
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d resp=%v", status, resp)
	}
	if errMsg, _ := resp["error"].(string); !strings.Contains(errMsg, "scope") {
		t.Fatalf("expected scope error, got %v", errMsg)
	}
}

func TestDouyinCallRequiresAccessToken(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	payload := defaultDouyinCallPayload(server)

	status, resp := invokeCallEndpoint(t, server, payload)
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d resp=%v", status, resp)
	}
	if resp["error"] != "missing access token: 请在 payload.access_token / 环境变量 / Flow 中提供 AccessToken" {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
}

func newTestAccessTokenServer(t *testing.T, configYAML string) *accessTokenServer {
	t.Helper()
	t.Setenv("ACCESSTOKEN_REDIS_ADDR", "memory")
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configYAML), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	server, closer, err := newAccessTokenServer(configPath, defaultListenAddr)
	if err != nil {
		t.Fatalf("newAccessTokenServer: %v", err)
	}
	t.Cleanup(func() {
		closer()
	})
	return server
}

func defaultTokenPayload(server *accessTokenServer) map[string]any {
	return map[string]any{
		"provider_code":      "bilbili",
		"provider_app":       "content_center",
		"provider_auth_mode": "default",
		"config_path":        server.defaultConfigPath,
	}
}

func invokeTokenEndpoint(t *testing.T, server *accessTokenServer, payload map[string]any) (int, tokenAPIResponse) {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/accesstoken/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.handleToken(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		errBody, _ := io.ReadAll(res.Body)
		t.Fatalf("token endpoint status=%d body=%s", res.StatusCode, string(errBody))
	}
	var resp tokenAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return res.StatusCode, resp
}

func defaultDouyinCallPayload(server *accessTokenServer) map[string]any {
	payload := map[string]any{
		"provider_code":      "byte_dance_douyin",
		"provider_app":       "douyin",
		"provider_auth_mode": "default",
		"config_path":        server.defaultConfigPath,
		"action":             "douyin.video.list",
		"payload":            map[string]any{},
	}
	return payload
}

func invokeCallEndpoint(t *testing.T, server *accessTokenServer, payload map[string]any) (int, map[string]any) {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal call payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/accesstoken/call", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.handleCall(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read call response: %v", err)
	}
	var resp map[string]any
	if len(data) > 0 {
		_ = json.Unmarshal(data, &resp)
	}
	return res.StatusCode, resp
}

func numberValue(v any) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case json.Number:
		n, _ := val.Int64()
		return int(n)
	default:
		return 0
	}
}
