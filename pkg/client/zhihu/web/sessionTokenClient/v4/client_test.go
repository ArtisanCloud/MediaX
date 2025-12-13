package v4

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
	zhmiddleware "github.com/ArtisanCloud/MediaX/server/zhihu/sessionToken/middleware"
)

func TestMeFollowingsSuccess(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Cookie"); !strings.Contains(got, "SESSION=abc") {
			t.Fatalf("expected cookie header to propagate session token, got %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []string{"foo"}})
	}))
	t.Cleanup(upstream.Close)
	stub := &stubSessionTokenManager{
		flow: &sessiontoken.Flow{
			FlowID:          "stf_test",
			TenantUUID:      "tenant_demo",
			ProviderAppCode: "zhihu_article",
			Status:          sessiontoken.FlowStatusSucceeded,
		},
	}
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil, WithBaseURL(upstream.URL), WithHTTPClient(upstream.Client()))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/me/followings?limit=5", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=abc")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	var resp apiSuccessResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Meta.FlowID != stub.flow.FlowID {
		t.Fatalf("expected flow id %s, got %s", stub.flow.FlowID, resp.Meta.FlowID)
	}
	if resp.Meta.TenantUUID != stub.flow.TenantUUID {
		t.Fatalf("expected tenant uuid %s, got %s", stub.flow.TenantUUID, resp.Meta.TenantUUID)
	}
	if resp.Meta.ProviderApp != stub.flow.ProviderAppCode {
		t.Fatalf("expected provider app code %s, got %s", stub.flow.ProviderAppCode, resp.Meta.ProviderApp)
	}
	if stub.lastToken != "SESSION=abc" {
		t.Fatalf("expected token lookup to use raw header, got %s", stub.lastToken)
	}
}

func TestChannelsArticlesForbiddenMarksInvalid(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "risk detected"})
	}))
	t.Cleanup(upstream.Close)
	stub := &stubSessionTokenManager{
		flow: &sessiontoken.Flow{
			FlowID:          "stf_401",
			TenantUUID:      "tenant_a",
			ProviderAppCode: "zhihu_article",
			Status:          sessiontoken.FlowStatusSucceeded,
		},
	}
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil, WithBaseURL(upstream.URL), WithHTTPClient(upstream.Client()))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/channels/abc/articles", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=def")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", rec.Code, rec.Body.String())
	}
	var resp apiError
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Code != codeRiskBlock {
		t.Fatalf("expected code %s, got %s", codeRiskBlock, resp.Code)
	}
	if len(stub.markCalls) != 1 {
		t.Fatalf("expected MarkTokenInvalid to be called once, got %d", len(stub.markCalls))
	}
	if stub.markCalls[0].Code != codeRiskBlock {
		t.Fatalf("expected failure code recorded, got %+v", stub.markCalls[0])
	}
}

func TestArticlePostValidationError(t *testing.T) {
	stub := &stubSessionTokenManager{
		flow: &sessiontoken.Flow{
			FlowID:          "stf_post",
			TenantUUID:      "tenant",
			ProviderAppCode: "zhihu_article",
			Status:          sessiontoken.FlowStatusSucceeded,
		},
	}
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/zhihu/v1/articles", strings.NewReader(`{"title":""}`))
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=xyz")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
	if len(stub.markCalls) != 0 {
		t.Fatalf("expected mark invalid not to be called on validation error")
	}
}

func TestSanityCheckFlowMissing(t *testing.T) {
	stub := &stubSessionTokenManager{
		findErr: sessiontoken.ErrFlowNotFound,
	}
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/zhihu/v1/sanity/check", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=xyz")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}
	var resp apiError
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Code != codeCookieExpired {
		t.Fatalf("expected code %s, got %s", codeCookieExpired, resp.Code)
	}
}

func TestCallZhihuRetriesWithAPIConfig(t *testing.T) {
	hits := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	t.Cleanup(upstream.Close)
	stub := &stubSessionTokenManager{}
	cfg := &config.ZhihuSessionTokenConfig{
		API: config.ZhihuSessionTokenAPIConfig{
			RetryBackoff: []int{1},
		},
	}
	client, err := NewClient(cfg, stub, nil,
		WithBaseURL(upstream.URL),
		WithHTTPClient(upstream.Client()),
		WithRetrySchedule([]time.Duration{0}),
	)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	status, _, _, err := client.callZhihu(context.Background(), http.MethodGet, "/foo", url.Values{}, nil, "SESSION=abc")
	if err != nil {
		t.Fatalf("callZhihu returned error: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}
	if hits != 2 {
		t.Fatalf("expected upstream to be hit twice, got %d", hits)
	}
}

type stubSessionTokenManager struct {
	flow      *sessiontoken.Flow
	findErr   error
	lastToken string
	markCalls []markInvalidCall
}

type markInvalidCall struct {
	FlowID  string
	Code    string
	Message string
	API     string
}

func (s *stubSessionTokenManager) CreateFlow(ctx context.Context, opts *sessiontoken.CreateFlowOptions) (*sessiontoken.Flow, error) {
	return nil, errors.New("not implemented")
}

func (s *stubSessionTokenManager) GetFlow(ctx context.Context, flowID string) (*sessiontoken.Flow, error) {
	return nil, errors.New("not implemented")
}

func (s *stubSessionTokenManager) CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*sessiontoken.Flow, error) {
	return nil, errors.New("not implemented")
}

func (s *stubSessionTokenManager) CompleteFlowFailed(ctx context.Context, flowID string, failure *sessiontoken.FlowFailure) (*sessiontoken.Flow, error) {
	return nil, errors.New("not implemented")
}

func (s *stubSessionTokenManager) MarkTokenInvalid(ctx context.Context, flowID string, code string, message string, api string) (*sessiontoken.Flow, error) {
	s.markCalls = append(s.markCalls, markInvalidCall{FlowID: flowID, Code: code, Message: message, API: api})
	return s.flow, nil
}

func (s *stubSessionTokenManager) FindFlowBySessionToken(ctx context.Context, token string) (*sessiontoken.Flow, error) {
	s.lastToken = token
	if s.findErr != nil {
		return nil, s.findErr
	}
	if s.flow == nil {
		return nil, sessiontoken.ErrFlowNotFound
	}
	return s.flow, nil
}
