package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	if stored.FlowTTLSeconds != server.flowTTLSeconds {
		t.Fatalf("FlowTTLSeconds = %d, want %d", stored.FlowTTLSeconds, server.flowTTLSeconds)
	}
	wantExpire := fresh.StoredAt.Add(time.Duration(server.flowTTLSeconds) * time.Second)
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
	if resp.FlowTTLSeconds != server.flowTTLSeconds {
		t.Fatalf("flow ttl seconds mismatch: %d vs %d", resp.FlowTTLSeconds, server.flowTTLSeconds)
	}
	if resp.FlowExpireAt.IsZero() {
		t.Fatalf("flow_expire_at should not be zero")
	}
}

func newTestAccessTokenServer(t *testing.T, configYAML string) *accessTokenServer {
	t.Helper()
	t.Setenv("ACCESSTOKEN_REDIS_ADDR", "")
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
