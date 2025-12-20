package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type oauthStartResponse struct {
	FlowID         string            `json:"flow_id"`
	State          string            `json:"state"`
	AuthorizeURL   string            `json:"authorize_url"`
	ExpiresIn      int               `json:"expires_in"`
	ExpireAt       time.Time         `json:"expire_at"`
	ListenAddr     string            `json:"listen_addr"`
	StorageBackend string            `json:"storage_backend"`
	OAuthKey       string            `json:"oauth_key"`
	ConfigPath     string            `json:"config_path"`
	Provider       map[string]string `json:"provider"`
}

func TestHandleOAuthStartBuildsAuthorizeURL(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	payload := map[string]any{
		"provider_code":      "bilbili",
		"provider_app":       "content_center",
		"provider_auth_mode": "default",
		"config_path":        server.defaultConfigPath,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/accesstoken/oauth/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.handleOAuthStart(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(res.Body)
		t.Fatalf("status=%d body=%s", res.StatusCode, string(errBody))
	}
	var resp oauthStartResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.FlowID == "" || !strings.HasPrefix(resp.FlowID, "oauth-") {
		t.Fatalf("unexpected flow_id: %s", resp.FlowID)
	}
	if resp.State == "" {
		t.Fatalf("state should be populated")
	}
	if resp.ListenAddr != server.listenAddr {
		t.Fatalf("listen_addr = %s, want %s", resp.ListenAddr, server.listenAddr)
	}
	if resp.ExpiresIn != server.flowTTLSeconds {
		t.Fatalf("expires_in = %d, want %d", resp.ExpiresIn, server.flowTTLSeconds)
	}
	if resp.ExpireAt.IsZero() {
		t.Fatalf("expire_at should be set")
	}
	if resp.StorageBackend != server.storageBackend {
		t.Fatalf("storage_backend = %s, want %s", resp.StorageBackend, server.storageBackend)
	}
	if resp.ConfigPath != server.defaultConfigPath {
		t.Fatalf("config_path = %s, want %s", resp.ConfigPath, server.defaultConfigPath)
	}
	u, err := url.Parse(resp.AuthorizeURL)
	if err != nil {
		t.Fatalf("parse authorize_url: %v", err)
	}
	q := u.Query()
	if got := q.Get("client_id"); got != "demo-client" {
		t.Fatalf("client_id = %s, want demo-client", got)
	}
	if got := q.Get("redirect_uri"); got != "http://localhost:7071/debug/callback" {
		t.Fatalf("redirect_uri = %s, want http://localhost:7071/debug/callback", got)
	}
	if got := q.Get("scope"); got != "arc_base" {
		t.Fatalf("scope = %s, want arc_base", got)
	}
	if got := q.Get("response_type"); got != "code" {
		t.Fatalf("response_type = %s, want code", got)
	}
	if got := q.Get("state"); got == "" {
		t.Fatalf("state query param empty")
	}
	if got := q.Get("state"); got != resp.State {
		t.Fatalf("state mismatch: query=%s body=%s", got, resp.State)
	}
	server.oauthStateMu.Lock()
	_, ok := server.oauthStates[resp.State]
	server.oauthStateMu.Unlock()
	if !ok {
		t.Fatalf("state %s not registered", resp.State)
	}
}

func TestHandleDebugCallbackPersistsMaskedPayload(t *testing.T) {
	tokenEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"flow-token-xyz","refresh_token":"refresh-secret","token_type":"Bearer","scope":"arc_base","expires_in":3600}`)
	}))
	defer tokenEndpoint.Close()

	server := newTestAccessTokenServer(t, testAccessTokenConfigWithTokenURL(tokenEndpoint.URL))
	state := "callback-state"
	server.oauthStateMu.Lock()
	server.oauthStates[state] = &oauthState{
		ProviderCode: "bilbili",
		ProviderApp:  "content_center",
		ModeKey:      "default",
		ConfigPath:   server.defaultConfigPath,
		CreatedAt:    time.Now(),
	}
	server.oauthStateMu.Unlock()

	body := `{"code":"abc123","access_token":"body-secret"}`
	req := httptest.NewRequest(http.MethodGet, "/debug/callback?code=abc123&state="+state+"&provider=bilibili", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.handleDebugCallback(rec, req)

	res := rec.Result()
	bodyBytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("callback status=%d body=%s", res.StatusCode, string(bodyBytes))
	}
	flowID := "oauth-" + state
	stored := server.fetchTokenByFlowID(flowID)
	if stored == nil {
		t.Fatalf("expected stored Flow %s; response=%s", flowID, string(bodyBytes))
	}
	if stored.Callback == nil {
		t.Fatalf("callback payload missing in Flow record")
	}
	if stored.Callback.FlowID != flowID {
		t.Fatalf("callback flow_id = %s, want %s", stored.Callback.FlowID, flowID)
	}
	if stored.Callback.Method != http.MethodGet {
		t.Fatalf("callback method = %s, want GET", stored.Callback.Method)
	}
	if stored.Callback.Body == "" || strings.Contains(stored.Callback.Body, "body-secret") {
		t.Fatalf("callback body not masked: %s", stored.Callback.Body)
	}
	if stored.Callback.Query == "" || strings.Contains(stored.Callback.Query, "abc123") {
		t.Fatalf("callback query not masked: %s", stored.Callback.Query)
	}
	if ct := stored.Callback.Headers["Content-Type"]; ct != "application/json" {
		t.Fatalf("callback header Content-Type = %s, want application/json", ct)
	}
}

func TestHandleOAuthStartSupportsRedbook(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	payload := map[string]any{
		"provider_code":      "redbook_juguang",
		"provider_app":       "juguang",
		"provider_auth_mode": "default",
		"config_path":        server.defaultConfigPath,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/accesstoken/oauth/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.handleOAuthStart(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(res.Body)
		t.Fatalf("status=%d body=%s", res.StatusCode, string(errBody))
	}
	var resp oauthStartResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.OAuthKey != "redbook-default" {
		t.Fatalf("oauth_key=%s want redbook-default", resp.OAuthKey)
	}
	authURL, err := url.Parse(resp.AuthorizeURL)
	if err != nil {
		t.Fatalf("parse authorize_url: %v", err)
	}
	if host := authURL.Host; host != "ad.xiaohongshu.com" {
		t.Fatalf("authorize host=%s want ad.xiaohongshu.com", host)
	}
	scope := authURL.Query().Get("scope")
	if scope != "notes.read notes.write" {
		t.Fatalf("scope=%s want notes.read notes.write", scope)
	}
	if resp.Provider["code"] != "redbook_juguang" {
		t.Fatalf("provider code=%s want redbook_juguang", resp.Provider["code"])
	}
}

func TestHandleOAuthStartRequiresScope(t *testing.T) {
	cfg := strings.Replace(testAccessTokenConfig, "scope: \"notes.read,notes.write\"", "scope: \"\"", 1)
	server := newTestAccessTokenServer(t, cfg)
	payload := map[string]any{
		"provider_code":      "redbook_juguang",
		"provider_app":       "juguang",
		"provider_auth_mode": "default",
		"config_path":        server.defaultConfigPath,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/accesstoken/oauth/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.handleOAuthStart(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d, want 400", res.StatusCode)
	}
	bodyBytes, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(bodyBytes), "OAuth scope 未配置") {
		t.Fatalf("unexpected error body: %s", bodyBytes)
	}
}

func TestHandleDebugCallbackPersistsRedbookFlow(t *testing.T) {
	tokenEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"redbook-flow-token","token_type":"Bearer","expires_in":7200}`)
	}))
	defer tokenEndpoint.Close()

	server := newTestAccessTokenServer(t, testAccessTokenConfigWithRedbookTokenURL(tokenEndpoint.URL))
	state := "redbook-callback"
	server.oauthStateMu.Lock()
	server.oauthStates[state] = &oauthState{
		ProviderCode: "redbook_juguang",
		ProviderApp:  "juguang",
		ModeKey:      "default",
		ConfigPath:   server.defaultConfigPath,
		CreatedAt:    time.Now(),
	}
	server.oauthStateMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/debug/callback?code=abc123&state="+state+"&provider=redbook_juguang", strings.NewReader(`{"debug":"body"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.handleDebugCallback(rec, req)

	res := rec.Result()
	bodyBytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("callback status=%d body=%s", res.StatusCode, string(bodyBytes))
	}
	flowID := "oauth-" + state
	stored := server.fetchTokenByFlowID(flowID)
	if stored == nil {
		t.Fatalf("expected stored Flow %s; response=%s", flowID, string(bodyBytes))
	}
	if stored.ProviderCode != "redbook_juguang" {
		t.Fatalf("provider_code=%s want redbook_juguang", stored.ProviderCode)
	}
	if stored.Scope != "notes.read notes.write" {
		t.Fatalf("stored scope=%s want notes.read notes.write", stored.Scope)
	}
	if stored.AccessToken != "redbook-flow-token" {
		t.Fatalf("access_token=%s want redbook-flow-token", stored.AccessToken)
	}
	if stored.Callback == nil {
		t.Fatalf("expected callback payload recorded")
	}
}
